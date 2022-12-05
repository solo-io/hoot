package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	serverv3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	testv3 "github.com/envoyproxy/go-control-plane/pkg/test/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"

	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	runtimeservice "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	secretservice "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
)

// This is mostly copied from
// https://github.com/envoyproxy/go-control-plane/tree/master/internal/example
// // Copyright 2020 Envoyproxy Authors

const (
	// don't use dots in resource name
	ClusterName1  = "cluster_a"
	ClusterName2  = "cluster_b"
	RouteName     = "local_route"
	ListenerName  = "listener_0"
	ListenerPort  = 10000
	UpstreamHost  = "127.0.0.1"
	UpstreamPort1 = 8082
	UpstreamPort2 = 8083

	xdsPort                  = 9977
	grpcMaxConcurrentStreams = 1000000
)

func makeCluster(clusterName string) *cluster.Cluster {
	return &cluster.Cluster{
		Name:                 clusterName,
		ConnectTimeout:       ptypes.DurationProto(5 * time.Second),
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_LOGICAL_DNS},
		LbPolicy:             cluster.Cluster_ROUND_ROBIN,
		LoadAssignment:       makeEndpoint(clusterName),
		DnsLookupFamily:      cluster.Cluster_V4_ONLY,
	}
}

func makeEndpoint(clusterName string) *endpoint.ClusterLoadAssignment {
	port := uint32(UpstreamPort1)
	if clusterName == ClusterName2 {
		port = UpstreamPort2
	}
	return &endpoint.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []*endpoint.LocalityLbEndpoints{{
			LbEndpoints: []*endpoint.LbEndpoint{{
				HostIdentifier: &endpoint.LbEndpoint_Endpoint{
					Endpoint: &endpoint.Endpoint{
						Address: &core.Address{
							Address: &core.Address_SocketAddress{
								SocketAddress: &core.SocketAddress{
									Protocol: core.SocketAddress_TCP,
									Address:  UpstreamHost,
									PortSpecifier: &core.SocketAddress_PortValue{
										PortValue: port,
									},
								},
							},
						},
					},
				},
			}},
		}},
	}
}

func makeRoute(routeName string, weight uint32, clusterName1, clusterName2 string) *route.RouteConfiguration {
	routeConfiguration := &route.RouteConfiguration{
		Name: routeName,
		VirtualHosts: []*route.VirtualHost{{
			Name:    "local_service",
			Domains: []string{"*"},
		}},
	}
	switch weight {
	case 0:
		routeConfiguration.VirtualHosts[0].Routes = []*route.Route{{
			Match: &route.RouteMatch{
				PathSpecifier: &route.RouteMatch_Prefix{
					Prefix: "/",
				},
			},
			Action: &route.Route_Route{
				Route: &route.RouteAction{
					ClusterSpecifier: &route.RouteAction_Cluster{
						Cluster: clusterName1,
					},
					HostRewriteSpecifier: &route.RouteAction_HostRewriteLiteral{
						HostRewriteLiteral: UpstreamHost,
					},
				},
			},
		}}

	case 100:
		routeConfiguration.VirtualHosts[0].Routes = []*route.Route{{
			Match: &route.RouteMatch{
				PathSpecifier: &route.RouteMatch_Prefix{
					Prefix: "/",
				},
			},
			Action: &route.Route_Route{
				Route: &route.RouteAction{
					ClusterSpecifier: &route.RouteAction_Cluster{
						Cluster: clusterName2,
					},
					HostRewriteSpecifier: &route.RouteAction_HostRewriteLiteral{
						HostRewriteLiteral: UpstreamHost,
					},
				},
			},
		}}
		// canary-roll out:
	default:
		routeConfiguration.VirtualHosts[0].Routes = []*route.Route{{
			Match: &route.RouteMatch{
				PathSpecifier: &route.RouteMatch_Prefix{
					Prefix: "/",
				},
			},
			Action: &route.Route_Route{
				Route: &route.RouteAction{
					ClusterSpecifier: &route.RouteAction_WeightedClusters{
						WeightedClusters: &route.WeightedCluster{
							TotalWeight: &wrapperspb.UInt32Value{
								Value: 100,
							},
							Clusters: []*route.WeightedCluster_ClusterWeight{
								{
									Name: clusterName1,
									Weight: &wrapperspb.UInt32Value{
										Value: 100 - weight,
									},
								},
								{
									Name: clusterName2,
									Weight: &wrapperspb.UInt32Value{
										Value: weight,
									},
								},
							},
						},
					},
					HostRewriteSpecifier: &route.RouteAction_HostRewriteLiteral{
						HostRewriteLiteral: UpstreamHost,
					},
				},
			},
		}}

	}
	return routeConfiguration
}

func makeHTTPListener(listenerName string, route string) *listener.Listener {
	// HTTP filter configuration
	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.HttpConnectionManager_AUTO,
		StatPrefix: "http",
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				ConfigSource:    makeConfigSource(),
				RouteConfigName: route,
			},
		},
		HttpFilters: []*hcm.HttpFilter{{
			Name: wellknown.Router,
		}},
	}
	pbst, err := ptypes.MarshalAny(manager)
	if err != nil {
		panic(err)
	}

	return &listener.Listener{
		Name: listenerName,
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: ListenerPort,
					},
				},
			},
		},
		FilterChains: []*listener.FilterChain{{
			Filters: []*listener.Filter{{
				Name: wellknown.HTTPConnectionManager,
				ConfigType: &listener.Filter_TypedConfig{
					TypedConfig: pbst,
				},
			}},
		}},
	}
}

func makeConfigSource() *core.ConfigSource {
	source := &core.ConfigSource{}
	source.ResourceApiVersion = resource.DefaultAPIVersion
	source.ConfigSourceSpecifier = &core.ConfigSource_ApiConfigSource{
		ApiConfigSource: &core.ApiConfigSource{
			TransportApiVersion:       resource.DefaultAPIVersion,
			ApiType:                   core.ApiConfigSource_GRPC,
			SetNodeOnFirstMessageOnly: true,
			GrpcServices: []*core.GrpcService{{
				TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
					EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "xds_cluster"},
				},
			}},
		},
	}
	return source
}

var (
	version int
)

func GenerateSnapshot(weight uint32) (*cachev3.Snapshot, error) {
	version++
	nextversion := fmt.Sprintf("snapshot-%d", version)
	fmt.Println("publishing version: ", nextversion)
	return cachev3.NewSnapshot(
		nextversion, // version needs to be different for different snapshots
		map[resource.Type][]types.Resource{
			resource.EndpointType: {},
			resource.ClusterType: {
				makeCluster(ClusterName1),
				makeCluster(ClusterName2),
			},
			resource.RouteType: {
				makeRoute(RouteName, weight, ClusterName1, ClusterName2),
			},
			resource.ListenerType: {
				makeHTTPListener(ListenerName, RouteName),
			},
			resource.RuntimeType: {},
			resource.SecretType:  {},
		},
	)
}

func registerServer(grpcServer *grpc.Server, server serverv3.Server) {
	// register services
	discoverygrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	endpointservice.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	clusterservice.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	routeservice.RegisterRouteDiscoveryServiceServer(grpcServer, server)
	listenerservice.RegisterListenerDiscoveryServiceServer(grpcServer, server)
	secretservice.RegisterSecretDiscoveryServiceServer(grpcServer, server)
	runtimeservice.RegisterRuntimeDiscoveryServiceServer(grpcServer, server)
}

// RunServer starts an xDS server at the given port.
func RunServer(ctx context.Context, srv3 serverv3.Server, port uint) {
	// gRPC golang library sets a very small upper bound for the number gRPC/h2
	// streams over a single TCP connection. If a proxy multiplexes requests over
	// a single connection to the management server, then it might lead to
	// availability problems.
	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions, grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams))
	grpcServer := grpc.NewServer(grpcOptions...)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	registerServer(grpcServer, srv3)

	log.Printf("management server listening on %d\n", port)
	if err = grpcServer.Serve(lis); err != nil {
		log.Println(err)
	}
}

type ClusterNodeHasher struct{}

// ID uses the node ID field
func (ClusterNodeHasher) ID(node *core.Node) string {
	if node == nil {
		return ""
	}
	return node.Cluster
}

func main() {
	ctx := context.Background()

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	l := logger.Sugar()

	nodeGroup := "edge-gateway"

	// Create a cache
	cache := cachev3.NewSnapshotCache(false, ClusterNodeHasher{}, l)

	// Create the snapshot that we'll serve to Envoy
	snapshot, err := GenerateSnapshot(0)
	if err != nil {
		l.Errorf("could not generate snapshot: %+v", err)
		os.Exit(1)
	}
	if err := snapshot.Consistent(); err != nil {
		l.Errorf("snapshot inconsistency: %+v\n%+v", snapshot, err)
		os.Exit(1)
	}
	l.Debugf("will serve snapshot %+v", snapshot)

	// Add the snapshot to the cache
	if err := cache.SetSnapshot(ctx, nodeGroup, snapshot); err != nil {
		l.Errorf("snapshot error %q for %+v", err, snapshot)
		os.Exit(1)
	}

	// Run the xDS server
	cb := &testv3.Callbacks{Debug: true}
	srv := serverv3.NewServer(ctx, cache, cb)
	go RunServer(ctx, srv, xdsPort)

	reader := bufio.NewReader(os.Stdin)
	for {

		fmt.Print("Enter weight for cluster-b: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		weight, _ := strconv.Atoi(text)
		clampedWeight := clampWeight(weight)
		fmt.Println("setting weight to", clampedWeight)

		snapshot, err = GenerateSnapshot(clampedWeight)
		if err != nil {
			l.Errorf("could not generate snapshot: %+v", err)
			os.Exit(1)
		}
		if err := snapshot.Consistent(); err != nil {
			l.Errorf("snapshot inconsistency: %+v\n%+v", snapshot, err)
			os.Exit(1)
		}
		l.Debugf("will serve snapshot %+v", snapshot)

		// Add the snapshot to the cache
		if err := cache.SetSnapshot(ctx, nodeGroup, snapshot); err != nil {
			l.Errorf("snapshot error %q for %+v", err, snapshot)
			os.Exit(1)
		}

	}
}

func clampWeight(weight int) uint32 {
	// if weight > 100 {
	// 	weight = 100
	// }
	if weight < 0 {
		weight = 0
	}
	return uint32(weight)
}
