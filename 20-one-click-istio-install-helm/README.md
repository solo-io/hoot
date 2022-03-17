# Hoot Episode 20: 1-Click Upgrade to Istio 1.13 using Helm
## Recording ##
https://www.youtube.com/watch?v=Q3G5TEmXq7o

[show notes](SHOWNOTES.md)

[slides](1-click-helm-slides.pdf)

## Abstract
No software is ever finished. Software engineers are continuously working on new features, improving performance and addressing security issues.

Istio maintainers are regularly releasing new minor and path versions to provide the best possible control plane for your mesh of services.

Istio provides multiple upgrade strategies, in this session we will look into in-place upgrades with the help of the official Helm charts and helmfile. 


## Platform setup

```console
minikube start --memory=16384 --cpus=4 --kubernetes-version=v1.20.2
```

## Deploying Istio

Clone this repository, go the the folder of this episode, then:

```console
$ helmfile sync
```

> Optionally, you can pass the `--debug` flag in order to get a more verbose output.

Sample output:

```console
UPDATED RELEASES:
NAME                   CHART           VERSION
istio-base             istio/base       1.12.5
istio-discovery        istio/istiod     1.12.5
istio-ingressgateway   istio/gateway    1.12.5
```

With the help of `istioclt`, you can validate the versions of control and dataplanes.

```console
$ istioctl version

client version: 1.13.2
control plane version: 1.12.5
data plane version: 1.12.5 (1 proxies)
```

## Deploying Bookinfo

First, we are enabling the injection of Istio sidecars in the `default` namespace.

```console
$ kubectl label namespace default istio-injection=enabled
```

Then, we are deploying our sample application, called Bookinfo.

```console
$ kubectl apply -f https://raw.githubusercontent.com/istio/istio/master/samples/bookinfo/platform/kube/bookinfo.yaml
```

Then finally, we are adding an Istio gateway, to make it available from outside of the cluster.

```console
$ kubectl apply -f https://raw.githubusercontent.com/istio/istio/master/samples/bookinfo/networking/bookinfo-gateway.yaml
```

We can validate that everything is working by crafting the URL of the application.

```console
$ export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')
$ export INGRESS_HOST=$(minikube ip)
$ export GATEWAY_URL=$INGRESS_HOST:$INGRESS_PORT
```

Then, we can print this:

```console
$ echo $GATEWAY_URL
192.168.64.2:30701
```

After this, we can curl the Productpage.

```console
$ curl -s "http://${GATEWAY_URL}/productpage"
```

## Upgrading Istio

### Upgrading the control plane

Bump version of charts in `helmfile.yaml` to the desired Istio version, e.g. 1.13.2

```console
$ helmfile sync|apply
```

Sample output:

```console
UPDATED RELEASES:
NAME                    CHART                                        VERSION
istio-base              istio/base                                    1.13.2
istio-discovery         istio/istiod                                  1.13.2
istio-ingressgateway    istio/gateway                                 1.13.2
```

Validation via istioctl:

```console
$ istioctl version

client version: 1.13.2
control plane version: 1.13.2
data plane version: 1.12.3 (7 proxies) 
```

### Rolling restart of workloads

First, we perform a rolling restart of the ingress-gateway.

```console
$ kubectl rollout restart deployment istio-ingressgateway -n istio-system
```

```console
$ kubectl rollout restart deployment -n default
```

Validation via istioctl

```console
$ istioctl version

client version: 1.13.2
control plane version: 1.13.2
data plane version: 1.13.2 (7 proxies)
```

## Monitoring the upgrade

Istio are enhanced with different types of metrics, covering both the control and the data plane.
These are all compatible with Prometheus, the de-facto monitoring tool for cloud native system.

If you add the `istioperformance.json` dashboard to your Grafana, you can explore the improvements of your Istio upgrades via these metrics, e.g. performance optimalizations.


## Cleanup

You can delete all the previously installed Helm releases.

```console
$ helmfile destroy
```

Or, you can also delete you Minikube cluster

```console
$ minikube delete
```
