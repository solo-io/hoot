# Demo

Start envoy
```
python -m http.server --bind 127.0.0.1 8082 &
podman run -ti --rm --net=host -v ${PWD}:${PWD} -w ${PWD} docker.io/envoyproxy/envoy:v1.18.2 -c envoy.yaml &
```
Start server:
```
(cd server; go run server.go)
```

Now you can curl
```
curl http://localhost:10000
```

# References

https://www.envoyproxy.io/docs/envoy/v1.18.2/api-v3/extensions/filters/http/ext_authz/v3/ext_authz.proto.html

https://www.envoyproxy.io/docs/envoy/v1.18.2/api-v3/extensions/filters/http/ratelimit/v3/rate_limit.proto#extensions-filters-http-ratelimit-v3-ratelimit

https://www.envoyproxy.io/docs/envoy/v1.18.2/api-v3/extensions/access_loggers/grpc/v3/als.proto#envoy-v3-api-msg-extensions-access-loggers-grpc-v3-httpgrpcaccesslogconfig

https://www.envoyproxy.io/docs/envoy/v1.18.2/api-v3/config/metrics/v3/metrics_service.proto

https://github.com/open-policy-agent/opa-envoy-plugin

https://github.com/envoyproxy/envoy/tree/main/api/envoy/service