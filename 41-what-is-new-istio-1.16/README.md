# Hoot Episode 41 - What is new with Istio 1.16?

## Recording ##
 https://www.youtube.com/watch?v=hO8PlznWbmI

[show notes](SHOWNOTES.md)

## Hands-on: Steps from the demo

### Steps for Lin's demo 1 - Ambient profile
Prepare for the environment:

```
kind create cluster --config=- <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: ambient
nodes:
- role: control-plane
- role: worker
- role: worker
EOF
```

Download Istio 1.16 binary, review and deploy ambient profile:

```
istioctl version -s
cat manifests/profiles/ambient.yaml
istioctl install --set profile=ambient
```

Review pods deployed:
```
kubectl get pods -A
```

Deploy sleep and helloworld app:
```
kubectl label ns default istio-injection=enabled
kubectl apply -f samples/sleep/sleep.yaml
kubectl apply -f samples/helloworld/helloworld.yaml
```

Check the new tunnel label:
```
kubectl get po -lapp=sleep -o yaml
```

Query all pods that have the new tunnel label configured:

```
kubectl get po -lnetworking.istio.io/tunnel=http -A
```

### Steps for Lin's demo 2 - Discovery selectors

Create the configuration to enable the new discovery selectors features:
```
cat ambient-discovery.yaml
```

Output:
```
# The ambient profile has ambient mesh enabled
# Note: currently this only enables HBONE for sidecars, as the full ambient mode is not yet implemented.
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
spec:
  meshConfig:
    discoverySelectors:
      - matchLabels:
          istio-discovery: enabled
    defaultConfig:
      proxyMetadata:
        ISTIO_META_ENABLE_HBONE: "true"
  values:
    pilot:
      env:
        PILOT_ENABLE_INBOUND_PASSTHROUGH: "false"
        PILOT_ENABLE_HBONE: "true"
        ENABLE_ENHANCED_RESOURCE_SCOPING: "true"
```

Apply the ambient discovery config:

```
istioctl install -f manifests/profiles/ambient-discovery.yaml -y
```

Note the istiod pod is restarted.  Why?

```
kubectl get pods -A
```

Create namespace foo and bar, with foo, default and istio-system have the istio-discovery enabled label:

```
kubectl create ns foo
kubectl create ns bar
kubectl label ns foo istio-discovery=enabled
kubectl label ns default istio-discovery=enabled
kubectl label ns istio-system istio-discovery=enabled
kubectl get ns -L istio-discovery
```

Display the route for sleep:
```
istioctl pc route deploy/sleep
```

Apply the review virtual service to the foo and bar namespaces:

```
kubectl apply -f samples/bookinfo/networking/virtual-service-reviews-90-10.yaml -n foo
kubectl apply -f samples/bookinfo/networking/virtual-service-reviews-90-10.yaml -n bar
```

Display the route for the sleep deployment. You’ll notice the reviews.foo virtual service is in the route list but not the reviews.bar.
```
istioctl pc route deploy/sleep
```

View configmaps:

```
kubectl get cm -A
```

### Clean up for Lin's demo 1 and demo 2:

```
kind delete cluster --name ambient
```
