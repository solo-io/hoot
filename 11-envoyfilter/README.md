# Intro

## `EnvoyFilter` - what is it?
- Provides a means of customizing Istio Service Mesh via Opaque config patches which load directly to Envoy.
- Can be used to enable and configure features in Envoy that are not exposed by Istio.

## EnvoyFilter, Summary

- Istio provides an opinionated way of configuring Envoy. While this is convenient for high level users, it makes customizing the behavior of proxies difficult.
- The `EnvoyFilter` CRD provides a general tool for gaining more direct control over Envoy configuration, when Istio's opinionated API is not sufficient for advanced use cases
- `EnvoyFilter` is useful for use cases which require an alternative Proxy configuration than what is possible via the "vanilla" Istio CRDs.
- `EnvoyFilter` is considered "break-glass"/"use-at-your-own-risk" configuration as it bypasses the core Istio translation logic and is essentially opaque to Istio's internal logic.
- `EnvoyFilter` semantics can be thought of like a "JSON patch" for Envoy configuration. It is composed of a `patch` field containing a raw YAML/JSON patch for Envoy configuration, and `match` field which selects which part of the Istio-translated Envoy configuration into which the patch will be inserted/merged.
- `EnvoyFilter` is considered "break-glass"/"use-at-your-own-risk" configuration as it bypasses the core Istio translation logic and is essentially opaque to Istio's internal logic.


## EnvoyFilter thoughts

- Due to the complexity and fragility of managing EnvoyFilters, they are probably best managed by a k8s controller and a higher level CRD rather than directly managed by mesh operators. This was the basis of Solo.io's reasoning in [Gloo Mesh](https://docs.solo.io/gloo-mesh) which provides use-case specific CRDs which abstract away the EnvoyFilter.
- Providing an abstraction layer between the user and the EnvoyFilter adds a layer of protection against breaks, as the compatibility of EnvoyFilter patches with translated Istio proxy config is not guaranteed to be stable across Istio versions.
- As Istio's core API is fairly opinionated about data plane (proxy) configuration, an extension point such as the EnvoyFilter is absolutely necessary for advanced use cases like an [Enterprise-ready Istio distro](https://www.solo.io/blog/announcing-gloo-mesh-enterprise/).
- EnvoyFilter and WASM provide a rough but sufficient interface for extending Istio to a large variety of use cases.

# Demo & Discussion

## Installation (one-time)
Kubernetes setup: (requires `kind`, `kubectl`, `istioctl`)
```
# set up 2 kind clusters with istio + bookinfo installed
git clone https://github.com/solo-io/gloo-mesh
cd gloo-mesh
source ci/setup-funcs.sh
# create clusters
create_kind_cluster mgmt-cluster 32001
create_kind_cluster remote-cluster 32000
# install istio
install_istio mgmt-cluster 32001
install_istio remote-cluster 32000

# set up bookinfo namespaces
kubectl --context kind-mgmt-cluster create namespace bookinfo
kubectl --context kind-mgmt-cluster label ns bookinfo istio-injection=enabled --overwrite
kubectl --context kind-remote-cluster create namespace bookinfo
kubectl --context kind-remote-cluster label ns bookinfo istio-injection=enabled --overwrite

# install bookinfo without reviews-v3 to management cluster
kubectl --context kind-mgmt-cluster -n bookinfo apply -f ./ci/bookinfo.yaml -l 'app notin (details),version notin (v3)'
kubectl --context kind-mgmt-cluster -n bookinfo apply -f ./ci/bookinfo.yaml -l 'account'

# install only reviews-v3 to remote cluster
kubectl --context kind-remote-cluster -n bookinfo apply -f ./ci/bookinfo.yaml -l 'app notin (details),version in (v3)'
kubectl --context kind-remote-cluster -n bookinfo apply -f ./ci/bookinfo.yaml -l 'service=reviews'
kubectl --context kind-remote-cluster -n bookinfo apply -f ./ci/bookinfo.yaml -l 'account=reviews'
kubectl --context kind-remote-cluster -n bookinfo apply -f ./ci/bookinfo.yaml -l 'app=ratings'
kubectl --context kind-remote-cluster -n bookinfo apply -f ./ci/bookinfo.yaml -l 'account=ratings'

# wait for deployments to finish
kubectl --context kind-mgmt-cluster -n bookinfo rollout status deployment/productpage-v1 --timeout=300s
kubectl --context kind-mgmt-cluster -n bookinfo rollout status deployment/reviews-v1 --timeout=300s
kubectl --context kind-mgmt-cluster -n bookinfo rollout status deployment/reviews-v2 --timeout=300s

kubectl --context kind-remote-cluster -n bookinfo rollout status deployment/reviews-v3 --timeout=300s
```

## Set up Multi-Cluster Networking

These steps will allow you to send requests from a service in `mgmt-cluster` to `reviews-v3` in the `remote-cluster` on the dns address `reviews.bookinfo.svc.remote-cluster.global`:

- Apply ServiceEntry to `mgmt-cluster`:

```yaml
kubectl apply --context kind-mgmt-cluster -f - <<EOF
## apply to client cluster
apiVersion: networking.istio.io/v1beta1
kind: ServiceEntry
metadata:
  name: reviews.bookinfo.svc.remote-cluster.global
  namespace: istio-system
spec:
  addresses:
    - 241.208.99.7
  endpoints:
    - address: 172.18.0.3 #<- this should be the Node IP of the Remote Cluster 
      labels:
        cluster: remote-cluster
      ports:
        http: 32000
  hosts:
    - reviews.bookinfo.svc.remote-cluster.global
  location: MESH_INTERNAL
  ports:
    - name: http
      number: 9080
      protocol: HTTP
  resolution: DNS
  EOF
```

This will allow clients in mgmt-cluster.

- Apply Gateway to `remote-cluster`:

```yaml
kubectl apply --context kind-remote-cluster -f - <<EOF
#### apply to remote cluster

apiVersion: networking.istio.io/v1beta1
kind: Gateway
metadata:
  name: bookinfo-federation-bookinfo
  namespace: istio-system
spec:
  selector:
    istio: ingressgateway
  servers:
    - hosts:
        - '*.global'
      port:
        name: tls
        number: 15443
        protocol: TLS
      tls:
        mode: AUTO_PASSTHROUGH

---

apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: bookinfo-federation.bookinfo
  namespace: istio-system
spec:
  configPatches:
    - applyTo: NETWORK_FILTER
      match:
        context: GATEWAY
        listener:
          filterChain:
            filter:
              name: envoy.filters.network.sni_cluster
          portNumber: 15443
      patch:
        operation: INSERT_AFTER
        value:
          name: envoy.filters.network.tcp_cluster_rewrite
          typed_config:
            '@type': type.googleapis.com/istio.envoy.config.filter.network.tcp_cluster_rewrite.v2alpha1.TcpClusterRewrite
            cluster_pattern: \.remote-cluster.global$
            cluster_replacement: .cluster.local
  workloadSelector:
    labels:
      istio: ingressgateway

EOF
```

## Test via cURL

```bash
kubectl alpha debug --image=curlimages/curl@sha256:aa45e9d93122a3cfdf8d7de272e2798ea63733eeee6d06bd2ee4f2f8c4027d7c -n bookinfo $(kubectl get pod -n bookinfo | grep productpage | awk '{print $1}') -i -- curl -v http://reviews.bookinfo.svc.remote-cluster.global:9080/reviews/123
```

If everything worked the `curl` should yield a `200 OK`.

# Final thoughts

`EnvoyFilter` is extremely powerful, but can be dangerous. The Istio community should converge around a higher level solution (management plane) which abstracts the fragility of the `EnvoyFilter`. Gloo Mesh is one such possibility, but due to obvious bias I'll leave the determination of how best to abstract up to the reader.

# More Resources:
- EnvoyFilter API: https://istio.io/latest/docs/reference/config/networking/envoy-filter/
- Demo: https://www.youtube.com/watch?v=sUkeFAERvE8
- Gloo Mesh: https://docs.solo.io/gloo-mesh/
