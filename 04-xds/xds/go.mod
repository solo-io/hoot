module github.com/solo-io/xds

go 1.18

// We attempt to mirror the latest dependencies from https://github.com/solo-io/gloo

require (
	github.com/envoyproxy/go-control-plane v0.10.3
	github.com/golang/protobuf v1.5.2
	go.uber.org/zap v1.19.1
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/census-instrumentation/opencensus-proto v0.3.0 // indirect
	github.com/cncf/xds/go v0.0.0-20220314180256-7f1daf1720fc // indirect
	github.com/envoyproxy/protoc-gen-validate v0.6.7 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d // indirect
	golang.org/x/sys v0.0.0-20210816183151-1e6c022a8912 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220329172620-7be39ac1afc7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
