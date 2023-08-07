# Envoy Ext auth filter
See docs for external auth service:
https://www.envoyproxy.io/docs/envoy/latest/api-v2/service/auth/v2/external_auth.proto.html

Docs for the filter:
https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/ext_authz_filter

# OPA plugin architecture
One can extend OPA, and create a new binary with extra plugins.
https://www.openpolicyagent.org/docs/latest/extensions

This method creates a new binary with the extra functionality added.

# OPA Envoy plugin
A plugin that implements the envoy external authorization service in the opa server:

https://github.com/open-policy-agent/opa-envoy-plugin

# Demo
This demo is based on the code from here:
https://github.com/open-policy-agent/opa-envoy-plugin/blob/master/quick_start.yaml
See this URL for more information and examples. It has been adapted to run without kubernetes.

The relevant files:
- envoy.yaml: envoy configuration that uses external auth
- opa-config.yaml: OPA configuration that configures the envoy_ext_authz_grpc
- policy.rego: The policy to use for auth decisions

```
opa_envoy run --server --config-file=opa-config.yaml --addr=localhost:8181 --diagnostic-addr=0.0.0.0:8282 "--ignore=.*" policy.rego
```

```
docker run -p 8080:8080 --rm openpolicyagent/demo-test-server:v1
```

```
envoy -c envoy.yaml
```

```
SERVICE_URL=localhost:8000
ALICE_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiZ3Vlc3QiLCJzdWIiOiJZV3hwWTJVPSIsIm5iZiI6MTUxNDg1MTEzOSwiZXhwIjoxOTQxMDgxNTM5fQ.rN_hxMsoQzCjg6lav6mfzDlovKM9azaAjuwhjq3n9r8"
BOB_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJzdWIiOiJZbTlpIiwibmJmIjoxNTE0ODUxMTM5LCJleHAiOjE5NDEwODE1Mzl9.ek3jmNLPclafELVLTfyjtQNj0QKIEGrbhKqpwXmQ8EQ"
```
The secret of the JWTs is `secret`

Now you can curl!
for example, alice can do a GET but not a POST:
```
curl -i -H "Authorization: Bearer "$ALICE_TOKEN"" http://$SERVICE_URL/people
curl -i -H "Authorization: Bearer "$ALICE_TOKEN"" -d '{"firstname":"Charlie", "lastname":"OPA"}' -H "Content-Type: application/json" -X POST http://$SERVICE_URL/people
curl -i -H "Authorization: Bearer "$BOB_TOKEN"" http://$SERVICE_URL/people
curl -i -H "Authorization: Bearer "$BOB_TOKEN"" -d '{"firstname":"Charlie", "lastname":"OPA"}' -H "Content-Type: application/json" -X POST http://$SERVICE_URL/people
```

See full info here:
https://github.com/open-policy-agent/opa-envoy-plugin
