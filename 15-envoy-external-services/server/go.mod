module github.com/solo-io/ext-svc-demo

go 1.16

require (
	contrib.go.opencensus.io/exporter/zipkin v0.1.2
	github.com/envoyproxy/go-control-plane v0.9.8
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/pkg/errors v0.9.1 // indirect
	github.com/uber/jaeger-client-go v2.25.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	go.opencensus.io v0.22.4
	go.uber.org/atomic v1.6.0 // indirect
	google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55
	google.golang.org/grpc v1.27.0
)
