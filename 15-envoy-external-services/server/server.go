package main

import (
	"context"
	"fmt"
	"net"

	als_pb "github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v3"
	auth_pb "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	metrics_pb "github.com/envoyproxy/go-control-plane/envoy/service/metrics/v3"
	rl_pb "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v3"

	"google.golang.org/grpc"
)

var (
	localEndpointURI = "127.0.0.1:10004"
)

type service struct{}

func (s *service) ShouldRateLimit(ctx context.Context, r *rl_pb.RateLimitRequest) (*rl_pb.RateLimitResponse, error) {
	fmt.Println("received should rate limit")
	// return nil, fmt.Errorf("error")
	return &rl_pb.RateLimitResponse{
		OverallCode: rl_pb.RateLimitResponse_OK,
	}, nil
}
func (s *service) Check(ctx context.Context, r *auth_pb.CheckRequest) (*auth_pb.CheckResponse, error) {
	fmt.Println("received check request")
	// return nil, fmt.Errorf("error")
	return &auth_pb.CheckResponse{}, nil
}

func (s *service) StreamAccessLogs(r als_pb.AccessLogService_StreamAccessLogsServer) error {
	for {
		msg, err := r.Recv()
		if err != nil {
			return err
		}
		for _, le := range msg.GetHttpLogs().LogEntry {
			fmt.Println("received access log", le.GetRequest().GetPath(), le.GetResponse().GetResponseCode().GetValue())
		}
	}
}

func (s *service) StreamMetrics(r metrics_pb.MetricsService_StreamMetricsServer) error {
	for {
		msg, err := r.Recv()
		if err != nil {
			return err
		}
		metrics := msg.EnvoyMetrics
		for _, em := range metrics {
			name := em.GetName()
			if name == "cluster.somecluster.upstream_rq_total" {
				for _, m := range em.Metric {
					fmt.Println("received metric", name, m.GetCounter().GetValue())
				}
			}
		}
	}
}

func main() {

	grpcServer := grpc.NewServer()

	lis, err := net.Listen("tcp", localEndpointURI)
	if err != nil {
		panic(err)
	}
	s := &service{}
	rl_pb.RegisterRateLimitServiceServer(grpcServer, s)
	auth_pb.RegisterAuthorizationServer(grpcServer, s)
	als_pb.RegisterAccessLogServiceServer(grpcServer, s)
	metrics_pb.RegisterMetricsServiceServer(grpcServer, s)
	fmt.Println("Starting")
	grpcServer.Serve(lis)
}
