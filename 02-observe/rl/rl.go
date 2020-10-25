package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	pb "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v2"

	"time"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/zpages"
	"google.golang.org/grpc"
)

var (
	localEndpointURI = "127.0.0.1:10004"
)

type service struct{}

func (s *service) ShouldRateLimit(ctx context.Context, r *pb.RateLimitRequest) (*pb.RateLimitResponse, error) {
	// this is done automatically:
	// ctx, span := trace.StartSpan(ctx, "ShouldRateLimit")
	// defer span.End()

	time.Sleep(time.Second)
	return &pb.RateLimitResponse{
		OverallCode: pb.RateLimitResponse_OK,
	}, nil
}

func setupZpage() {
	mux := http.NewServeMux()
	zpages.Handle(mux, "/debug")

	// Change the address as needed
	addr := ":8888"
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Failed to serve zPages " + err.Error())
	}
}

func main() {
	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		log.Fatalf("Failed to register ocgrpc server views: %v", err)
	}
	cfg := jaegercfg.Configuration{
		ServiceName: "ratelimit",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		panic("Could not initialize jaeger tracer: " + err.Error())
	}
	defer closer.Close()

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
		otgrpc.OpenTracingServerInterceptor(tracer)),
		grpc.StreamInterceptor(
			otgrpc.OpenTracingStreamServerInterceptor(tracer)))

	lis, err := net.Listen("tcp", localEndpointURI)
	if err != nil {
		panic(err)
	}
	pb.RegisterRateLimitServiceServer(grpcServer, &service{})
	fmt.Println("Starting")
	go setupZpage()
	grpcServer.Serve(lis)
}
