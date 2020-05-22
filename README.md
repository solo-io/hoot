# Concept

What is:
- envoy
  - Extendability
  - configuration via API
- listeners
- clusters
- endpoints
- TCP filters
- HCM
- http filters
- routes

## Diagrams
+-------------------+          +----------+           +----------------------------+           +------------------+
|                   |          |          |           |                            |           |                  |
| downstream client +--------->+ listener +---------->+ filters (routing decision) +---------->+ upstream cluster |
|                   |          |          |           |                            |           |                  |
+-------------------+          +----------+           +----------------------------+           +------------------+

How is routing decision done?
+-------------+         +------------+          +--------------+           +---------------+          +----------------+
|             |         |            |          |              |           |               |          |                |
| tcp filters +-------->+ HCM filter +--------->+ http filters +---------->+ router filter +--------->+ host selection |
|             |         |            |          |              |           |               |          |                |
+-------------+         +------------+          +--------------+           +---------------+          +----------------+


# Interactive session

## Simple config

run the following:
```
fuser -k 8082/tcp
fuser -k 10000/tcp
fuser -k 10004/tcp

go run server_tcp.go&
envoy -c simple_tcp.yaml&
echo hi | nc localhost 10000
```

## Http simple config:
```
fuser -k 8082/tcp
fuser -k 10000/tcp
fuser -k 10004/tcp

go run server.go&
envoy -c simple.yaml&
curl http://localhost:10000 -dhi
```

Lets explore the stats:
```
curl localhost:9901/stats|grep somecluster.upstream_rq
```

stats that end with _total are usually counters.
stats that end with _active are usually gauges.

Navigating the config dump:

```
curl localhost:9901/config_dump | grep -C 10 prefix
```

## running with debug logs

```
fuser -k 8082/tcp
fuser -k 10000/tcp
fuser -k 10004/tcp

go run server.go&
envoy -c simple.yaml -l debug&
# or if envoy is already running, turn them on dynamically:
# curl -XPOST "http://localhost:9901/logging?level=debug"
curl http://localhost:10000 -dhi
```

let's talk about what we are seeing in these debug logs:
-- the headers coming in, arriving to the router filter, and being sent upstream.
 
## let's use a filter!

remember: filter order matters.

```
fuser -k 8082/tcp
fuser -k 10000/tcp
fuser -k 10004/tcp

go run server.go&
envoy -c simple_fault.yaml -l debug&

for i in $(seq 10); do
curl http://localhost:10000 -s -o /dev/null -w "%{http_code}" 
echo
done
```

### cors

note: you can provide route level configuration for a filter

```
fuser -k 8082/tcp
fuser -k 10000/tcp
fuser -k 10004/tcp

go run server.go&
envoy -c cors.yaml&
curl -XOPTIONS http://localhost:10000 -H"Origin: solo.io" -v
curl -XOPTIONS http://localhost:10000 -H"Origin: example.com" -v
```

### headers

note: some route level configuration is handled by the router filter

```
fuser -k 8082/tcp
fuser -k 10000/tcp
fuser -k 10004/tcp

go run server.go&
envoy -c response-header.yaml&
curl http://localhost:10000 -v
```

### auth and RL

remember: route order matter.
note: for newer filters, route level configuration is provided using typed_per_filter_config for extendability.

```
fuser -k 8082/tcp
fuser -k 10000/tcp
fuser -k 10004/tcp

go run server.go&
(cd rl_auth; go run main.go)&
envoy -c rl_auth.yaml&
for i in $(seq 10); do
    curl http://localhost:10000 -s -o /dev/null -w "%{http_code}" 
    echo
done
for i in $(seq 10); do
    curl http://localhost:10000/static -s -o /dev/null -w "%{http_code}" 
    echo
done
```

# What is all this xDS stuff?!

Essentially its a method for distributing configuration.
Every piece of dynamic configuration can also be configured statically.
almost true the other way around too.

Envoy has a bunch of xDS:
- LDS
- CDS
- EDS
- RDS

and ADS to merge them all.
the can come from files, grpc or REST api.
can be SotW or Delta updates.

