# Hoot Episode 52 - Istio and Open Policy Agent (OPA)

## Recording
https://www.youtube.com/watch?v=EnckV6lyre8


## Demo - Istio and OPA

Prerequisites:
- Kubernetes cluster
- Istio installed (1.18.0)


## Deploy sample app (httpbin)

Label the default namespace with `istio-injection=enabled` and deploy the httpbin workl;oad:

```sh
kubectl label namespace default istio-injection=enabled
kubectl apply -f https://raw.githubusercontent.com/istio/istio/master/samples/httpbin/httpbin.yaml
```

To access the httpbin service, we'll create a gateway and a virtual service:

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: gateway
spec:
  selector:
    istio: ingressgateway
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "example.com"
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: httpbin
spec:
  hosts:
  - "example.com"
  gateways:
  - gateway
  http:
  - route:
    - destination:
        port:
          number: 8000
        host: httpbin
```

You can send a request to the external IP to ensure everything works as it should.

## Install OPA (ConfigMap)

We'll install OPA and a ConfigMap that contains the policy we want to enforce:

```sh
kubectl create ns opa
kubectl apply -f opa-cm.yaml
```

The policy we'll enforce is the following:

```rego
package istio.authz

import future.keywords
import input.attributes.request.http as http_request

default allow := false

allow if {
    http_request.method == "GET"
    http_request.path == "/headers"
}
```

We'll deny all requests, expect if they are GET requests to the `/headers` path.

## Configure the external authorizer in mesh config

Edit the ConfigMap and add the extension provider to the meshConfig:

```shell
kubectl edit configmap istio -n istio-system
```

Add the following:

```yaml
    extensionProviders:
    - name: "opa"
      envoyExtAuthzGrpc:
        service: "opa.opa.svc.cluster.local"
        port: "9191"
        includeRequestBodyInCheck:
            maxRequestBytes: 4096
            packAsBytes: false
```

## Enforce the policy

We can now reference the extension provider by name (`opa`) in AuthorizationPolicy with the CUSTOM action:

```yaml
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: httpbin-opa
spec:
  selector:
    matchLabels:
      app: httpbin
  action: CUSTOM
  provider:
    name: "opa"
  # Trigger OPA for all requests to httpbin
  rules:
  - {}
```

If we send a request to the external IP, we'll get a 403 Forbidden response. If we send a request to `/headers`, we'll get a 200 OK response.

## OPA bundles

OPA bundles are a better way to distribute policies (compared to updating the ConfigMap and restarting the OPA workload). We'll create a bundle that contains the policy we want to enforce in the `authz` folder:

```sh
opa build authz
```

Assuming you've configured your GCP storage, you can use the gsutil and copy the `bundle.tar.gz` file to the bucket.

To use configure OPA to use the bundles, update the bucket name in the `opa-bundle.yaml` file and deploy the `opa-bundle.yaml` (instead of the ConfigMap based deployed we did earlier).

To test the policies, you can use the JWT tokens from the `tokens.md` file.

## Testing the policies

To test the policies, we can use `opa test` command. The `tests` folder contains a policy and tests for the policy. To run the tests, run the following command from the `tests` folder:

```sh
opa test -v .
```