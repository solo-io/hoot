
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
 
## Let's use a filter!


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

