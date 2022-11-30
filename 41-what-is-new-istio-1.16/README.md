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

Call helloworld from sleep:

```
kubectl exec -it $(k get po -lapp=sleep -ojsonpath='{.items[0].metadata.name}') -- curl helloworld:5000/hello
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

### Wasm Plugin demo

Assumes an already provisioned cluster, for example GKE provisioned with [this script](https://raw.githubusercontent.com/dhawton/work-tests/master/clouds/gke/create.sh).

Install Istio 1.16.0 with the default profile:

```bash
istioctl install -y
```

Ensure the LoadBalancer is created and IP assigned:

```bash
kubectl get svc istio-ingressgateway -n istio-system
```

Store this as it will be needed later:
  
```bash
export INGRESS_IP=$(kubectl get svc istio-ingressgateway -n istio-system -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
```

Create a namespace for the demo and label for sidecar injection:

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: Namespace
metadata:
  name: wasm
  labels:
    istio-injection: enabled
EOF
```

Deploy the httpbin application and gateway to the `wasm` namespace (assumes istio is installed in ~/istio):

```bash
kubectl apply -n wasm -f https://raw.githubusercontent.com/istio/istio/release-1.16/samples/httpbin/httpbin.yaml
kubectl apply -n wasm -f https://raw.githubusercontent.com/istio/istio/release-1.16/samples/httpbin/httpbin-gateway.yaml
```

Set the Envoy container to info logging for the Wasm logger, as by default the log level is set to `warn`:

```bash
kubectl exec -it -n wasm deploy/httpbin -c istio-proxy -- curl -X POST localhost:15000/logging?wasm=info
```

Apply the Wasm plugin:

```bash
kubectl apply -n wasm -f - <<EOF
apiVersion: extensions.istio.io/v1alpha1
kind: WasmPlugin
metadata:
  name: httpbin-rust-test
spec:
  selector:
    matchLabels:
      app: httpbin
  match:
  - ports:
    - number: 80
  url: oci://docker.io/dhawton/wasm-rust-test:v1
  imagePullPolicy: Always
EOF
```

This will take a bit to download the image and apply. Monitor the sidecar logs to wait for the Wasm plugin to be loaded:

```text
2022-11-28T22:47:20.264572Z	info	wasm	fetching image dhawton/wasm-rust-test from registry index.docker.io with tag v1
2022-11-28T22:47:21.528261Z	info	envoy wasm	wasm log: on_vm_start
 ** snip **
2022-11-28T22:47:21.591880Z	info	envoy wasm	wasm log: on_vm_start
```

There will be multiple on_vm_start's, that is normal as each thread starts up its V8 VM and loads the Wasm plugin.

Now, send a request to the httpbin service:

```bash
curl -s http://$INGRESS_IP/get
```

Expected output:

```text
❯ curl -s http://$INGRESS_IP/get
I'm a teapot

                       (
            _           ) )
         _,(_)._        ((
    ___,(_______).        )
  ,'__.   /       \    /\_
 /,' /  |""|       \  /  /
| | |   |__|       |,'  /
 \`.|                  /
  `. :           :    /
    `.            :.,'
      `-.________,-'
```

Apply a new Wasm plugin that changes the port to demonstrate new matching:

```bash
kubectl apply -n wasm -f - <<EOF
apiVersion: extensions.istio.io/v1alpha1
kind: WasmPlugin
metadata:
  name: httpbin-rust-test
spec:
  selector:
    matchLabels:
      app: httpbin
  match:
  - ports:
    - number: 81
  url: oci://docker.io/dhawton/wasm-rust-test:v1
  imagePullPolicy: Always
EOF
```

Give a few seconds for the new configuration to be generated and loaded, then re-run the curl command:

```bash
curl -s http://$INGRESS_IP/get
```

A teapot or two may appear initially, but the new configuration should be applied and the normal httpbin response seen:

```json
{
  "args": {}, 
  "headers": {
    "Accept": "*/*", 
    "Host": "35.227.152.117", 
    "User-Agent": "curl/7.81.0", 
    "X-B3-Parentspanid": "c95fed33660096e4", 
    "X-B3-Sampled": "0", 
    "X-B3-Spanid": "b8baaacf7f5ba2f4", 
    "X-B3-Traceid": "afa8fe480fba6538c95fed33660096e4", 
    "X-Envoy-Attempt-Count": "1", 
    "X-Envoy-Internal": "true", 
    "X-Forwarded-Client-Cert": "By=spiffe://cluster.local/ns/wasm/sa/httpbin;Hash=a83168142ce34878164ff018029a1b00c3764d47c08377b9bfc7a4c305f6b382;Subject=\"\";URI=spiffe://cluster.local/ns/istio-system/sa/istio-ingressgateway-service-account"
  }, 
  "origin": "10.138.15.235", 
  "url": "http://35.227.152.117/get"
}
```

### Cleanup Wasm Demo Cluster

Tear down the cluster. If it was built in GCP:

```bash
gcloud container clusters delete (cluster name) --zone=(zone)
```