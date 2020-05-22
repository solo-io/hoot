package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"

	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	epb "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	pb "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v2"
	"google.golang.org/grpc/codes"

	status "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
)

type service struct{}

func (s *service) ShouldRateLimit(context.Context, *pb.RateLimitRequest) (*pb.RateLimitResponse, error) {

	if rand.Int31()%2 == 0 {
		return &pb.RateLimitResponse{
			OverallCode: pb.RateLimitResponse_OVER_LIMIT,
		}, nil
	}

	return &pb.RateLimitResponse{
		OverallCode: pb.RateLimitResponse_OK,
	}, nil
}

func (s *service) Check(context.Context, *epb.CheckRequest) (*epb.CheckResponse, error) {
	if rand.Int31()%2 == 0 {
		return &epb.CheckResponse{
			Status: &status.Status{
				Code: int32(codes.PermissionDenied),
			},
		}, nil
	}
	return &epb.CheckResponse{
		Status: &status.Status{},
		HttpResponse: &epb.CheckResponse_OkResponse{OkResponse: &epb.OkHttpResponse{
			Headers: []*core.HeaderValueOption{{
				Header: &core.HeaderValue{
					Key:   "x-yuval",
					Value: "authorized",
				},
			}},
		}},
	}, nil
}

func main() {
	grpcServer := grpc.NewServer()

	lis, err := net.Listen("tcp", ":10004")
	if err != nil {
		panic(err)
	}
	pb.RegisterRateLimitServiceServer(grpcServer, &service{})
	epb.RegisterAuthorizationServer(grpcServer, &service{})
	fmt.Println("Starting")
	grpcServer.Serve(lis)
}
