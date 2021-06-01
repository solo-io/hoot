Run some server (doesn't really matter what, just so we get an OK response):
```
python -m http.server --bind 127.0.0.1 8082&
```

Run Redis, for the rate limit server:
```
redis-server &
```

Export these environment to configure the rate limit server, and run it:
```
export REDIS_SOCKET_TYPE=tcp
export REDIS_URL=localhost:6379
export LOG_LEVEL=debug
export USE_STATSD=false
export GRPC_PORT=10004
export RUNTIME_ROOT=$PWD
export RUNTIME_SUBDIRECTORY=rlconfig
export RUNTIME_WATCH_ROOT=false

go run github.com/envoyproxy/ratelimit/src/service_cmd
```

Run envoy:
```
envoy -c envoy.yaml&
```

Curl to the upstream
```
curl localhost:10000/
curl -XPOST http://localhost:10000/resources
# or use hey to send a bunch of requests
hey http://localhost:10000
hey -m POST http://localhost:10000/resources
```
Curl to the stats page
```
curl localhost:9901/stats|grep ratelimit 
```

Rate limit debug page: http://localhost:6070/rlconfig

References:
- https://www.envoyproxy.io/docs/envoy/latest/configuration/other_features/rate_limit
- https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/rate_limit_filter
- https://github.com/envoyproxy/ratelimit
- https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/route/v3/route_components.proto#config-route-v3-ratelimit-action
