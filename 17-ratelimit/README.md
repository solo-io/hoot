# Intro

Global rate limiting is performed using a server external to envoy that holds the rate limit counter state. The envoy projects provides an implementation that uses Redis. You can create your own server by
implementing the Rate Limit gRPC API.

Rate limits Descriptors are configured in envoy using a list of actions.
Each Action may generate a Descriptor Entry (a key and value). If it does not, the descriptor is not generated (usually).

So for a request that envoy wants to rate limit, it sends a list of descriptors, each descriptor has
an ordered list of descriptor entries, and each descriptor entry has a string key and value.

A visual example for 2 descriptors:
```
[(generic_key, foo), (header_value, bar)]
[(generic_key, abc), (remote_address, 1.2.3.4), (header_value, edf)]
```

The rate limit server then increments a counter for each descriptor. If for any of the descriptors the
counter goes above the limit defined in the server config, the request is rate limited.

# Demo    
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
