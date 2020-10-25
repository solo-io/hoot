
## Simple config

run the following in the background, to generate some traffic:
```
go run server.go&
(cd rl; go run rl.go)&
while true; do
 curl localhost:10000/
 curl localhost:10000/foo
 sleep 1
done
```

# Admin page
```
envoy -c stats.yaml
```
http://localhost:9901/


# Envoy debug logs

```
envoy -c simple.yaml -l debug
```

OR

```
envoy -c simple.yaml
```
and:
```
curl -XPOST "localhost:9901/logging?level=debug"
```

Envoy will now output logs in debug level.

# Access Logs

```
envoy -c accesslogs.yaml --file-flush-interval-msec 1
```

# Promethues
```
prometheus --config.file=prometheus.yml --web.listen-address="127.0.0.1:9090" --storage.tsdb.path=$(mktemp -d)
```
```
envoy -c stats.yaml
```

UI is in:
http://localhost:9090/

example query:
```
rate(envoy_listener_http_downstream_rq_xx{envoy_response_code_class="2"}[10s])
```

# Jaeger
Distributed tracing the the distributed analogy to a stack trace.
For example, you can see a regular go program stack trace here:
http://localhost:6060/debug/pprof/goroutine?debug=2

To generate some distributed traces:


Run jaeger:
```
docker run --rm --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 14250:14250 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.18
```

Run envoy:
```
envoy -c jaeger.yaml
```

Wait a few seconds, for the traffic generating command from the first section to 
generate some traffic.

See the traces in the jaeger UI: http://localhost:16686/

# Notes:
## tracing
The jaeger shared object for envoy with libstdc++:
https://github.com/jaegertracing/jaeger-client-cpp/releases/download/v0.4.2/libjaegertracing_plugin.linux_amd64.so

Or, for libc++ envoy, see here:
https://github.com/envoyproxy/envoy/issues/11382#issuecomment-638012072
(https://github.com/tetratelabs/getenvoy-package/files/3518103/getenvoy-centos-jaegertracing-plugin.tar.gz)


## access logs
more info here:
https://www.envoyproxy.io/docs/envoy/v1.15.0/api-v3/config/accesslog/v3/accesslog.proto#envoy-v3-api-msg-config-accesslog-v3-statuscodefilter
format string:
https://www.envoyproxy.io/docs/envoy/v1.15.0/configuration/observability/access_log/usage#config-access-log-format-strings