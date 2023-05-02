# Hoot Episode 49 - Configuring External Services with Istio's ServiceEntry
## Recording
https://www.youtube.com/watch?v=FOFzetag9d4

## Demo

### Prerequisites

- Kubernetes cluster with Istio installed
- Install Prometheus, Kiali, and Grafana


Label the namespace for injection and deploy the sleep pod:

```shell
kubectl label ns default istio-injection=enabled
kubectl apply -f  https://github.com/istio/istio/raw/master/samples/sleep/sleep.yaml
```

Send a couple of requests to httpbin.org:

```shell
kubectl exec -it deploy/sleep -- curl httpbin.org/headers
```


If you open Kiali Graph you'll see the calls going to Passthroughcluster.

```shell
kubectl port-forward svc/kiali -n istio-system 20001:20001
```

### Outbound traffic policy mode

The outbound traffic policy mode defines how Envoy sidecars handle external services, specifically services that are not part of the Istio's internal service registry.

The default setting is called `ALLOW_ANY`:

```shell
kubectl get istiooperator installed-state-istio -n istio-system -o jsonpath='{.spec.meshConfig.outboundTrafficPolicy.mode}'
```

It enables sidecars to send requests to unknown external services.

The other setting is called `REGISTRY_ONLY`. When set, Istio sidecars will block all requests to services that are not part of the internal service registry.

### What is PasshthroughCluster?

When `ALLOW_ANY` is set, we can see that requests to the external services go through the PassthroughCluster.

```shell
istioctl pc clusters deploy/sleep --fqdn PassthroughCluster -o yaml
```

PassthroughCluster is a original destination cluster created in Envoy configuration. `ORIGINAL_DST` is a special type of a cluster where the requests get proxied to the destination IP. The clusters in Envoy typically have a collection of endpoints, however, in the original destination cluster the endpoint is the original request destination. So when we make a request to `httpbin.org` (a service that is not in the mesh), the request is passed through to the original destination.

```shell
istioctl pc all deploy/sleep  -o json  > allcfg.json
```

We can search for "PassthroughCluster" or "PassthroughFilterChain" and you'll notice there's a default filter chain defined for listeners. So whenever other filter chains don't match, the default filter chain is selected.


### What is BlackHoleCluster?

Let's change the outbound traffic policy mode to `REGISTRY_ONLY` and see what happens.

```yaml
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
metadata:
  name: istio
spec:
  profile: demo
  meshConfig:
    outboundTrafficPolicy:
      mode: REGISTRY_ONLY
  values:
    global:
      meshID: mesh1
      multiCluster:
        clusterName: "Kubernetes"
      network: ""
    pilot:
      env:
        PILOT_ENABLE_WORKLOAD_ENTRY_AUTOREGISTRATION: "true"
        PILOT_ENABLE_WORKLOAD_ENTRY_HEALTHCHECKS: "true"
```

In this case, the default filter chain is added as well, however, it points to the BlackholeCluster that doesn't have any endpoints.


### Adding external service to the mesh

Setting the outbound traffic polciy to `REGISTRY_ONLY` is recommended, as it gives you control over what external services are allowed in the mesh.

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: ServiceEntry
metadata:
  name: httpbin
spec:
  hosts:
  - httpbin.org
  location: MESH_EXTERNAL
  ports:
  - number: 80
    name: http
    protocol: HTTP
  resolution: DNS
```

When we deploy the above SE, a cluster called `outbound|80||httpbin.org` gets added to the configuration. The cluster sets the endpoint to the `httpbin.org` -- the value we have in the hosts field, and it sets the Envoy cluster type to `STRICT_DNS`. So what servicentry does, is adds (and configures) a cluster in the Envoy configuration. If we'd look at Kiali now, we'd see that the calls to httpbin.org are going through the httpbin.org cluster and the item in the graph is correctly named as well.

#### Fields in the ServiceEntry resource

The **location** field specifies whether the service is external to the mesh, typically used for external services consumed through APIs, or whether the service is considered a part of the mesh, used for services running on VMs, for example. The difference is in how the mTLS authentication and policy enforcement works. With MESH_EXTERNAL services, the mTLS authentication is disabled, and the policy enforcement is performed on the callers' side instead of the server side (as we don't know if there's a proxy running there or not).

The second interesting field is the **resolution** field. We're setting it to DNS because we want to resolve the IP address of the `httpbin.org` host by querying the DNS. The other options are:

- NONE = this is for any connection that donhave already been resolved ,
- STATIC = this tells the proxy to use the IP addresses that are specified in another field in the resource
- DNS_ROUND_ROBIN = resolves the IP addresses by queryying the DNS asynchronlously. IT will usde the first IP address that's returned, without relying on complete DNS resolution. 

With the ServiceEntry defined, we can now use a VirtualService to configure things such as retries, aborts or delays, or even route to different services. The key is to use the same hostname (httpbin.org) in VirtualService as in the ServiceEntry resource.

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: httpbin
spec:
  hosts:
  - httpbin.org
  http:
  - route:
    - destination:
        host: httpbin.org
    fault:
      abort:
        percentage:
          value: 100
        httpStatus: 400
```

### What if we have an IP address?

Let's say that we have an IP address instead of a host name - how can we use the servicenetry to point to that IP address?

I'll use the ip of the httpbin.org site:

```shell
dig +short httpbin.org
```

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: ServiceEntry
metadata:
  name: httpbin
spec:
  hosts:
  - something.blah
  addresses:
    - 10.0.0.0 
  location: MESH_EXTERNAL
  ports:
  - number: 80
    name: http
    protocol: HTTP
  resolution: STATIC
  endpoints:
    - address: 52.86.68.46
```

We'll use the endpoints/addresses field and the STATIC resolution; however, we still have to specify how we want to address this service. We can't use the hosts field as we don't really have a host name, however, we could use the addresses field to specify a virtual IP address (or addresess) that we want to use to access this service. Let's take this example further. We have an IP address, but it would be so much better if we could use a host name instead. How can we do that?

### I have an IP, but I want to use a hostname

We can't use a host name in the serviceentry resource and hope that it will resolve to the IP address that we want. However, we can create a headless Kubernetes service - that would give us a stable hostname that points to the destination address.

So instead of defining a VIP and an endpoint in servieentry, we can define a host name and an endpoint in the headless Kubernetes service 

```yaml
apiVersion: v1
kind: Service
metadata:
  name: myhttpbinsvc
spec:
  clusterIP: None 
  ports:
    - protocol: TCP
      port: 80
      targetPort: 80
---
apiVersion: v1
kind: Endpoints
metadata:
  # Same as the Service name
  name: myhttpbinsvc
  namespace: default
subsets:
- addresses:
  - hostname: myhttpbinsvc-0
    ip: 52.86.68.46
  ports:
  - name: http
    port: 80
    protocol: TCP
```

As Istio treats this service as part of the service mesh (since it's A Kubernetes service), we don't even need a ServiceEntry for it at this point. The cluster created in Envoy is the original destination type, to the actual IP address is in the kubernetes endpoint isn't anything Istio knows about.

A downside of this approach is that the outbound traffic policy we configured earlier won't be respected for this service as it's considered part of the mesh - i.e. Istio doesn't know that this service points to an external IP. the outbound traffic policy isn't a security thing and you shouldn't be using and treating it like that, but it's good to be aware of how different things work and what the implications are. Kubernetes also has a concept of an ExternalName service. The ExternalName service gives us an alias we can use for the extenral DNS name; this is different from the headless services where we get a stable hostname, but we have to provide IP addreses.

Istio will treat the ExternalName as the external servvice, so the outbound policy will apply.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: externalhttpbin
spec:
  type: ExternalName
  externalName: httpbin.org
```

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: ServiceEntry
metadata:
  name: externalhttpbin
spec:
  hosts:
  - externalhttpbin
  location: MESH_EXTERNAL
  ports:
  - number: 80
    name: http
    protocol: HTTP
  resolution: DNS
```

### TLS origination

Another thing that often comes up is how to originate TLS from the proxy. So internally, you can to use http to call the external service, but you want to transparently upgrade the call to HTTPS. How can we do that?

It starts with the ServiceEntry, but this time, we'll define two ports -- port 80 and port 443:

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: ServiceEntry
metadata:
  name: cnn
spec:
  hosts:
  - edition.cnn.com
  ports:
  - number: 80
    name: http-port
    protocol: HTTP
  - number: 443
    name: https-port
    protocol: HTTPS
  resolution: DNS
```

If we send a request:

```shell
kubectl exec -it deploy/sleep -- curl -sSL -o /dev/null -v http://edition.cnn.com/politics
```

Notice the request goes to HTTP, but then it gets redirected to https - so curl is smart enough to handle the redirection, however, that's a redundant request and it doubles the latency. 

What we can do is to use the DestinationRule to redircet any HTTP requests and send them to port 443 instead. 

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: ServiceEntry
metadata:
  name: cnn
spec:
  hosts:
  - edition.cnn.com
  location: MESH_EXTERNAL
  ports:
  - number: 80
    name: http-port
    protocol: HTTP
    targetPort: 443
  - number: 443
    name: https-port
    protocol: HTTPS
  resolution: DNS
```

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: cnn
spec:
  host: edition.cnn.com
  trafficPolicy:
    portLevelSettings:
    - port:
        number: 80
      tls:
        mode: SIMPLE
```