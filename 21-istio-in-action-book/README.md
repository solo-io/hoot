# Hoot Episode 21 - Istio In Action Book
## Recording ##
https://www.youtube.com/watch?v=gpWuVnOyWnE

[show notes](SHOWNOTES.md)

[slides](Optimizing-the-control-plane-performance.pdf)

## Get the book

Get the book at 41% discount with the promo code: **SOLOIO41**

Go over to http://mng.bz/06Wl, and don't forget to use the promo code!

## Hands-on: Steps from the demo

Create a local kubernetes cluster:
```
kind create cluster
```

Get istioctl and install Istio and the addons in the cluster:

```
curl -L https://istio.io/downloadIstio | \
    ISTIO_VERSION=1.13.2 TARGET_ARCH=x86_64 sh -

## Move istioctl to your path before proceding

istioctl install -y --set profile=demo
kubectl apply -f istio-1.13.2/samples/addons 
```

Get the Istio in Action book source code, which contains the scripts and files used in this demo.
```
git clone https://github.com/istioinaction/book-source-code

cd book-source-code
```

Install the catalog service, another 12 dummy workloads (to increase the number of workloads managed by istiod), and apply 600 resources so that there is more configuration to process and propagate:
```
kubectl create ns istioinaction
kubectl label ns istioinaction istio-injection=enabled
kubectl config set-context $(kubectl config current-context) --namespace=istioinaction
kubectl -n istioinaction apply -f services/catalog/kubernetes/catalog.yaml
kubectl -n istioinaction apply -f ch11/catalog-virtualservice.yaml
kubectl -n istioinaction apply -f ch11/catalog-gateway.yaml
kubectl -n istioinaction apply -f ch11/sleep-dummy-workloads.yaml

kubectl create ns new-namespace
kubectl label ns new-namespace istio-injection=enabled
kubectl -n new-namespace apply -f https://gist.githubusercontent.com/rinormaloku/8766c15ebcf4f76c251092e041f75efc/raw/7ccb2603c7619b0d89ec2c1d5a559e42830cfaea/echo-service.yaml

kubectl -n istioinaction apply -f ./ch11/resources-600.yaml
```

You need a public IP address to route traffic to the Ingress Gateway, I will port-forward to the localhost:
```
kubectl port-forward -n istio-system svc/istio-ingressgateway 8080:80
```

Open grafana to visualize the control plane performance while you execute performance tests in later steps:
```
istioctl dashboard grafana
```

## Ignoring events

The best optimization is not having to do any processing at all. You can do that by ignoring irrelevant events. For example, with the following istio installation configuration, you configure the cluster to not listen for events in namespaces with some the label `istio-exclude=true`:
```
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
metadata:
  namespace: istio-system
spec:
  profile: demo
  meshConfig:
    discoverySelectors:
      - matchExpressions:
        - key: istio-exclude
          operator: NotIn
          values:
            - "true"
```

While the `new-namespace` is not excluded we will see that the service running in this cluster is propagated to workloads
```
$ istioctl pc clusters deploy/catalog.istioinaction | grep echo

SERVICE FQDN                               PORT    SUBSET     DIRECTION     TYPE    
echo-v1.new-namespace.svc.cluster.local    8080    -          outbound      EDS
```
Meanwhile, if we label the namespace to be ignored, this cluster will be removed, and subsequent changes won't be propagated at all.
```
kubectl label ns new-namespace istio-injection-
kubectl label ns new-namespace istio-exclude=true
```

Let's verify that the cluster is removed:
```
$ istioctl pc clusters deploy/catalog.istioinaction | grep echo
```
The output should be empty.

## Receiving only relevant updates

Using the following Sidecar configuration we set a mesh wide default for workloads to be configured for egress only in the listed namespaces:
```
apiVersion: networking.istio.io/v1beta1
kind: Sidecar
metadata:
  name: default
  namespace: istio-system
spec:
  egress:
  - hosts:
    - "istio-system/*"
    - "prometheus/*"
  outboundTrafficPolicy:
    mode: REGISTRY_ONLY
```

All other configuration won't be sent.

Execute a performance test before applying this optimization (if you didn't port forward the gateway to your localhost replace the value of the `--gateway` flag).

```
./bin/performance-test.sh --reps 10 --delay 1.5 \
    --prom-url prometheus.istio-system:9090 \
    --gateway localhost:8080
```

Additionally, you can check the size of the configuration, sent to every proxy:

```
kubectl -n istioinaction exec -ti deploy/catalog -c istio-proxy \
-- curl -s localhost:15000/config_dump > /tmp/config_dump

du -sh /tmp/config_dump
```

Apply the mesh wide sidecar configuraiton and execute the above steps again. Check out the graphs in the Istio Control Plane dashboard to see the difference.

```
kubectl apply -f ch11/sidecar-mesh-wide.yaml
```

#### Cleanup before going to the next section

```
kubectl delete sidecar -n istio-system default
```

## Batching changes for longer periods
Update the istio installation to batch changes for longer periods.
```
istioctl install --set profile=demo \
--set values.pilot.env.PILOT_DEBOUNCE_AFTER="1500ms" -y
```

Execute the performance test again:
```
./bin/performance-test.sh --reps 10 --delay 1 \
    --prom-url prometheus.istio-system:9090 \
    --gateway localhost:8080
```

You'll see that the changes were propagated in fewer pushes.

## Scaling up or out

The book Istio in Action, explains all the metrics that will help determine where is the bottleneck. After you determine the bottleneck:
- Scale up if the bottleneck is the rate of changes
- Scale out when the bottleneck is outgoing traffic, meaning that each istiod replica manages too many workloads.
