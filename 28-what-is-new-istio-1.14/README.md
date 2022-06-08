# Hoot Episode 28 - What is new with Istio 1.14?

## Recording ##
 https://www.youtube.com/watch?v=XxMAjo8yrLI

[show notes](SHOWNOTES.md)

## Hands-on: Steps from the demo

### Steps for Faseela's demo 1 - workloadSelector and credentialName in DR

Follow Istio doc:
https://istio.io/latest/docs/tasks/traffic-management/egress/egress-tls-origination/#mutual-tls-origination-for-egress-traffic
This document has all required configuration for workloadSelector and credentialName support in DR.


### Steps for Faseela's demo 2 - auto-sni support

Enable auto_sni experimental feature

```
# kubectl exec -n istio-system -it istiod-cc6987bb6-25kbj -- env | grep ENABLE_AUTO_SNI
ENABLE_AUTO_SNI=true
```

Configure DestinationRule with auto_sni enabled

Follow Istio doc:
https://istio.io/latest/docs/tasks/traffic-management/egress/egress-tls-origination/#mutual-tls-origination-for-egress-traffic
with the only difference of removing sni in DestinationRule configuration.

Refer sni section of ClientTLSSettings at https://istio.io/latest/docs/reference/config/networking/destination-rule/#ClientTLSSettings
for further details. The updated DR will look something like below:

```
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: first-destination-rule
spec:
  workloadSelector:
    matchLabels:
      app: sleep
  exportTo:
  - .
  host: my-nginx.mesh-external.svc.cluster.local
  trafficPolicy:
    loadBalancer:
      simple: ROUND_ROBIN
    portLevelSettings:
    - port:
        number: 31443
      tls:
        credentialName: client-credential
        mode: MUTUAL
        # no sni here
```

### Steps for Lin's demo 1 - minTLS for workloads in the mesh

Prepare for the environment:
```
minikube start --memory=16384 --cpus=4 --kubernetes-version=v1.23.7
brew install helmfile
helm repo list
helm repo update
```

From the demo dir of this episode:
```
helmfile sync
helmfile status
kubectl get pods -A
istioctl version
```

Follow Istio doc:
https://istio.io/latest/docs/tasks/security/tls-configuration/workload-min-tls-version/ & https://istio.io/latest/docs/setup/install/helm/#installation-steps, where to add the meshConfig?

```
helm show values istio/base
helm show values istio/istiod
```
Add the following to helmfile for istiod chart, in the demo/helmfile.yaml file:

```
      - meshConfig:
          meshMTLS:
            minProtocolVersion: TLSV1_3
```

Sync changes:

```
helmfile sync
```

View istio configmap in the istio-system namespace to observe the update.
Deploy sample apps:

```
kubectl create ns foo
kubectl apply -f <(istioctl kube-inject -f ~/Downloads/istio-1.14.0/samples/httpbin/httpbin.yaml) -n foo
kubectl apply -f <(istioctl kube-inject -f ~/Downloads/istio-1.14.0/samples/sleep/sleep.yaml) -n foo
```

Test with tls 1.3:

```
kubectl exec "$(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name})" -c istio-proxy -n foo -- openssl s_client -alpn istio -tls1_3 -connect httpbin.foo:8000
```

Test with tls 1.2:

```
kubectl exec "$(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name})" -c istio-proxy -n foo -- openssl s_client -alpn istio -tls1_2 -connect httpbin.foo:8000
```

Q: Why istio-proxy is used above?  And how does this work?

```
istioctl pc listener deploy/httpbin -n foo  --address 0.0.0.0 -o json --port 15006 | grep 'tlsMinimumProtocolVersion": "TLSv1_3"' -A 20 -B 5
```

Why 3 returns?  virtualInbound-catchall-http, virtualInbound, virtualInbound (raw_buffer as transport protocol)

Q: Does this work for gateway?

Yes: https://github.com/istio/api/blob/master/networking/v1beta1/gateway.proto#L723

### Lin's Demo 2: telemetry feature, workload mode for access logs

Enable access logs for the mesh:

```
kubectl apply -f - <<EOF
apiVersion: telemetry.istio.io/v1alpha1
kind: Telemetry
metadata:
  name: mesh-default
  namespace: istio-system
spec:
  accessLogging:
    - providers:
      - name: envoy
    - disabled: true
EOF
```

Enable access log in foo namespace:

```
kubectl apply -f - <<EOF
apiVersion: telemetry.istio.io/v1alpha1
kind: Telemetry
metadata:
  name: foo
  namespace: foo
spec:
  accessLogging:
    - providers:
      - name: envoy
EOF
```

Send some requests and observe logs for sleep and httpbin's proxy:

```
kubectl exec "$(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name})" -c istio-proxy -n foo -- curl httpbin.foo:8000
```

Note outbound access log not shown, why and how to fix it?

```
kubectl exec "$(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name})" -c sleep -n foo -- curl httpbin.foo:8000
```

Configure only to show client side:

```
kubectl apply -f - <<EOF
apiVersion: telemetry.istio.io/v1alpha1
kind: Telemetry
metadata:
  name: foo
  namespace: foo
spec:
  accessLogging:
    - providers:
      - name: envoy
      match:
        mode: CLIENT
EOF
```

Verify httpbin & sleep proxy logs shown up accordingly, e.g. only access logs on the client side.

### Clean up for Lin's demo 1 and demo 2:

```
kubectl delete telemetry -n istio-system mesh-default
kubectl delete telemetry -n foo foo
kubectl delete -f ~/Downloads/istio-1.14.0/samples/httpbin/httpbin.yaml -n foo
kubectl delete -f ~/Downloads/istio-1.14.0/samples/sleep/sleep.yaml -n foo
kubectl delete ns foo
```
