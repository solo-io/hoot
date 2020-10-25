run ssl server:

```
rm /tmp/envoy_admin.log
go run server.go&
envoy -l debug -c edge.yaml
```

See that access admin interface is logged:
```
curl localhost:9090/help
cat /tmp/envoy_admin.log
```


Sanity check, that everything is working:
```
curl --connect-to example.com:443:127.0.0.1:8443 -k -v https://example.com
```
(we use connect-to flag, so that curl creates the correct SNI record)


See envoy rejects request with an invalid header:
```
curl --connect-to example.com:443:127.0.0.1:8443 -k -v https://example.com -H"invalid_header: foo"
```

See how the XFF header is processed:
```
curl --connect-to example.com:443:127.0.0.1:8443 -k https://example.com -H"x-forwarded-for: 1.2.3.4"
envoy -c simple.yaml --disable-hot-restart &
curl http://localhost:10000 -H"x-forwarded-for: 1.2.3.4"
```

You can impact XFF processing using `use_remote_address` and `xff_num_trusted_hops`