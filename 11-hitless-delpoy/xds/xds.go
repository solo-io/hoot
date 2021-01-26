package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	serverv3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	testv3 "github.com/envoyproxy/go-control-plane/pkg/test/v3"

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
	ClusterName  = "upstream"
	UpstreamHost = "127.0.0.1"

	xdsPort                  = 9977
	grpcMaxConcurrentStreams = 1000000

	drainDuration = time.Second * 5
)

func makeEndpoint(ports ...uint32) *endpoint.ClusterLoadAssignment {

	var endpoints []*endpoint.LbEndpoint
	for _, port := range ports {
		endpoints = append(endpoints, &endpoint.LbEndpoint{
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
		})
	}

	return &endpoint.ClusterLoadAssignment{
		ClusterName: ClusterName,
		Endpoints: []*endpoint.LocalityLbEndpoints{{
			LbEndpoints: endpoints,
		}},
	}
}

var (
	version int
)

func GenerateSnapshot(ports ...uint32) cachev3.Snapshot {
	version++
	nextversion := fmt.Sprintf("snapshot-%d", version)
	fmt.Println("publishing version: ", nextversion)
	return cachev3.NewSnapshot(
		nextversion,                              // version needs to be different for different snapshots
		[]types.Resource{makeEndpoint(ports...)}, // endpoints
		[]types.Resource{},
		[]types.Resource{},
		[]types.Resource{},
		[]types.Resource{}, // runtimes
		[]types.Resource{}, // secrets
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
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	l := logger.Sugar()

	nodeGroup := "edge-gateway"

	// Create a cache
	cache := cachev3.NewSnapshotCache(false, ClusterNodeHasher{}, l)

	updateService := func(ports ...uint32) {
		// simulate a delay in propagating state to envoy
		time.Sleep(time.Second)

		snapshot := GenerateSnapshot(ports...)
		// Add the snapshot to the cache
		if err := cache.SetSnapshot(nodeGroup, snapshot); err != nil {
			l.Errorf("snapshot error %q for %+v", err, snapshot)
			os.Exit(1)
		}
	}

	fmt.Println("starting service 1")
	cancel1 := StartService("svc1", 8080)

	// Create the snapshot that we'll serve to Envoy
	updateService(8080)

	// Run the xDS server
	ctx := context.Background()
	cb := &testv3.Callbacks{Debug: true}
	srv := serverv3.NewServer(ctx, cache, cb)
	go RunServer(ctx, srv, xdsPort)

	reader := bufio.NewReader(os.Stdin)

	///////////////////////////////////
	// interesting part starts here: //
	///////////////////////////////////

	fmt.Print("Press enter to deploy svc 2, and remove service 1")
	reader.ReadString('\n')

	fmt.Print("Deploying svc 2")
	// start second service
	cancel2 := StartServiceWithDrain("svc2", 201, 8081)

	// Create the snapshot that we'll serve to Envoy
	updateService(8080, 8081)

	fmt.Print("Removing svc 1")
	// remove first service
	cancel1()

	updateService(8081)

	fmt.Print("Press enter to deploy svc 3, and drain service 2")
	reader.ReadString('\n')

	// deploy and drain
	fmt.Print("Deploying svc 3")
	StartServiceWithDrain("svc3", 202, 8082)
	// Create the snapshot that we'll serve to Envoy

	updateService(8081, 8082)

	// remove service 2
	cancel2()

	// wait with updating envoy during drain period.
	// envoy should also mark service as unhealthy during this period, and stop routing new requests.

	// remove drained service
	updateService(8082)
	fmt.Print("Press enter to exit")
	reader.ReadString('\n')
}

func StartService(text string, port int) func() {
	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(text))
	})

	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler}
	go server.ListenAndServe()

	return func() {
		go server.Shutdown(context.Background())
	}
}

func StartServiceWithDrain(text string, statusCode int, port int) func() {
	var drain uint32
	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		// handle health checks
		isDraining := atomic.LoadUint32(&drain) != 0

		if r.URL.Path == "/health" {
			if isDraining {
				rw.WriteHeader(http.StatusServiceUnavailable)
			} else {
				rw.WriteHeader(http.StatusOK)
			}
			return
		}

		if isDraining {
			rw.Header().Set("x-envoy-immediate-health-check-fail", "1")
		}
		rw.WriteHeader(statusCode)
		rw.Write([]byte(text))
	})

	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler}
	go server.ListenAndServe()

	return func() {
		go func() {
			atomic.AddUint32(&drain, 1)
			time.Sleep(drainDuration)
			server.Shutdown(context.Background())
		}()
	}
}
