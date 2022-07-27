# Hoot Episode 33 - Speed your Istio development environment with vcluster

## Recording ##
  https://youtu.be/b7OkYjvLf4Y 

[show notes](SHOWNOTES.md)

## Hands-on: Steps from the demo

### vcluster Demo

Requirements:
- A k8s cluster
- [Istioctl](https://istio.io/latest/docs/setup/install/istioctl/)

Review vcluster values:
```
cat vcluster-values.yaml
```

Get current context for the main cluster:
```
export MAIN_CLUSTER=$(kubectl config current-context)
```

Create a vcluster
```
vcluster create hoot-istio-test --expose -f vcluster-values.yaml --upgrade -n hoot-istio-test --connect=false --context $MAIN_CLUSTER
```

Make sure you do not have the context name already taken:
```
kubectl config delete-context hoot-istio-test
```

Connect to vcluster:
```
vcluster connect hoot-istio-test -n hoot-istio-test --update-current --kube-config-context-name hoot-istio-test --context $MAIN_CLUSTER
```

Install istio
```
istioctl install
```

Label the namespace to istio to inject the proxy (sidecar)
```
kubectl label namespace default istio-injection=enabled
```

Install httpbin app:
```
kubectl apply -f httpbin.yaml
```

Configure Istio to route traffic to httpbin app:
```
kubectl apply -f istio-httpbin-resources.yaml
```

Access nhttpbin app through istio
```
kubectl port-forward svc/istio-ingressgateway -n istio-system  8080:80
```

```
curl localhost:8080/get
```

Pause vcluster (all the resources will be scale to 0):
```
kubectl config use-context $MAIN_CLUSTER

vcluster pause hoot-istio-test
```

Try to connect to the paused cluster to check that it is down:
```
kubectl config use-context hoot-istio-test

kubectl get ns
```
The connection should fail.


Resume vcluster:
```
kubectl config use-context $MAIN_CLUSTER

vcluster resume hoot-istio-test
```

Wait until the cluster is back and try to connect:
```
kubectl config use-context hoot-istio-test

kubectl get ns
```

Connect to httpbin app:
```
kubectl port-forward svc/istio-ingressgateway -n istio-system  8080:80
```

```
curl localhost:8080/get
```


Delete the vcluster:
```
kubectl config use-context $MAIN_CLUSTER
vcluster delete hoot-istio-test -n hoot-istio-test
```
