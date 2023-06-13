# Hoot Episode 51 - Kubernetes Networking and Cilium - Part 2

## Recording
https://www.youtube.com/watch?v=he2sLeJsMqU


## Demo - Cilium & NetworkPolicy

We'll start by creating a kind cluster without any CNI and kube-proxy.

1. Create a kind cluster, without CNI.

    ```shell
    kind create cluster --config=files/kind-nocni.yaml
    ```

2. Install Cilium - use the strict kube-proxy replacement mode:

    ```shell
    helm install cilium cilium/cilium --version 1.13.2 --namespace kube-system --values files/cilium-values.yaml
    ```


Use the command below to check the status of Cilium:

```shell
cilium status
```


Once Cilium is installed, we'll also deploy Prometheus and Grafana:

```shell
kubectl apply -f https://raw.githubusercontent.com/cilium/cilium/v1.13/examples/kubernetes/addons/prometheus/monitoring-example.yaml
```

## Deploying sample apps

We'll use `sleep` and `httpbin` workloads for the demo:

```shell
kubectl apply -f sleep.yaml
kubectl apply -f https://raw.githubusercontent.com/istio/istio/master/samples/httpbin/httpbin.yaml
```

We can send a request to make sure there's connectivity between the two:

```shell
kubectl exec -it deploy/sleep -- curl httpbin:8000/headers
```

## Ingress policies

Next, let's deploy a deny all ingress policy -- this will deny all traffic between the pods in the default namespace:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-ingress
spec:
  podSelector: {}
  policyTypes:
  - Ingress
```

If we try to send a request this time, it's not going to work - the request will eventually time out.

We can use `hubble`to observe the traces from the command line:

```shell
# port forward to the Hubble relay running in the cluster
cilium hubble port-forward &


# then use hubble observe
hubble observe --to-namespace default -f
```

```console
Jun  9 20:38:30.645: default/sleep-6d68f49858-7ppbf:43145 (ID:43148) <- kube-system/coredns-565d847f94-wq69m:53 (ID:37585) to-overlay FORWARDED (UDP)
Jun  9 20:38:30.645: default/sleep-6d68f49858-7ppbf:43145 (ID:43148) <- kube-system/coredns-565d847f94-wq69m:53 (ID:37585) to-endpoint FORWARDED (UDP)
Jun  9 20:38:30.645: kube-system/coredns-565d847f94-wq69m:53 (ID:37585) <> default/sleep-6d68f49858-7ppbf (ID:43148) pre-xlate-rev TRACED (UDP)
Jun  9 20:38:30.645: kube-system/kube-dns:53 (world) <> default/sleep-6d68f49858-7ppbf (ID:43148) post-xlate-rev TRANSLATED (UDP)
Jun  9 20:38:30.645: kube-system/coredns-565d847f94-wq69m:53 (ID:37585) <> default/sleep-6d68f49858-7ppbf (ID:43148) pre-xlate-rev TRACED (UDP)
Jun  9 20:38:30.645: kube-system/kube-dns:53 (world) <> default/sleep-6d68f49858-7ppbf (ID:43148) post-xlate-rev TRANSLATED (UDP)
Jun  9 20:38:30.645: default/sleep-6d68f49858-7ppbf (ID:43148) <> default/httpbin-797587ddc5-kqb4f:80 (ID:1176) post-xlate-fwd TRANSLATED (TCP)
Jun  9 20:38:30.645: default/sleep-6d68f49858-7ppbf:53352 (ID:43148) <> default/httpbin-797587ddc5-kqb4f:80 (ID:1176) policy-verdict:none INGRESS DENIED (TCP Flags: SYN)
Jun  9 20:38:30.645: default/sleep-6d68f49858-7ppbf:53352 (ID:43148) <> default/httpbin-797587ddc5-kqb4f:80 (ID:1176) Policy denied DROPPED (TCP Flags: SYN)
Jun  9 20:38:31.678: default/sleep-6d68f49858-7ppbf:53352 (ID:43148) <> default/httpbin-797587ddc5-kqb4f:80 (ID:1176) policy-verdict:none INGRESS DENIED (TCP Flags: SYN)
Jun  9 20:38:31.678: default/sleep-6d68f49858-7ppbf:53352 (ID:43148) <> default/httpbin-797587ddc5-kqb4f:80 (ID:1176) Policy denied DROPPED (TCP Flags: SYN)
Jun  9 20:38:33.723: default/sleep-6d68f49858-7ppbf:53352 (ID:43148) <> default/httpbin-797587ddc5-kqb4f:80 (ID:1176) policy-verdict:none INGRESS DENIED (TCP Flags: SYN)
Jun  9 20:38:33.723: default/sleep-6d68f49858-7ppbf:53352 (ID:43148) <> default/httpbin-797587ddc5-kqb4f:80 (ID:1176) Policy denied DROPPED (TCP Flags: SYN)
Jun  9 20:38:37.757: default/sleep-6d68f49858-7ppbf:53352 (ID:43148) <> default/httpbin-797587ddc5-kqb4f:80 (ID:1176) policy-verdict:none INGRESS DENIED (TCP Flags: SYN)
Jun  9 20:38:37.757: default/sleep-6d68f49858-7ppbf:53352 (ID:43148) <> default/httpbin-797587ddc5-kqb4f:80 (ID:1176) Policy denied DROPPED (TCP Flags: SYN)
```

Hubble CLI supports a number of different output formats, including JSON and ability to filter the flows. The command below will show the flows between the sleep and httpbin endpoints:

```shell
hubble observe -o json --from-label "app=sleep" --to-label "app=httpbin"
```

There's also a UI part to Hubble - you can open it with the command below:

```shell
cilium hubble ui
```


With everyting denied, let's explicitly allow traffic from sleep to httpbin endpoints:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-ingress-to-httpbin
  namespace: default
spec:
  podSelector:
    matchLabels:
      app: httpbin
  policyTypes:
  - Ingress
  ingress:
    - from:
      - podSelector:
          matchLabels:
            app: sleep
```

If we repeat the request from the sleep workload, we'll see that the traffic is now allowed.

Let's say we deploy another sleep workload, but this time we'll deploy it to the "sleep" namespace:

```shell
kubectl create ns sleep
kubectl apply -f sleep.yaml -n sleep
```

Send a request to the httpbin workload from the sleep namespace will fail because of our deny all policy.

We can deploy a policy that explicitly allows making calls from the pods in the sleep namespace:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-ingress-from-sleep
  namespace: default
spec:
  podSelector:
    matchLabels:
      app: httpbin
  policyTypes:
  - Ingress
  ingress:
    - from:
      - namespaceSelector:
          matchLabels:
            kubernetes.io/metadata.name: sleep
```

If we send a request from the sleep namespace, we'll be able to access the httpbin workload in the default namespace:

```shell
kubeclt exec -it deploy/sleep -n sleep -- curl httpbin.default:8000/headers
```

Looking closer at this policy, you can see it's not very restrictive as it allows any pods in the sleep namespace to access the httpbin workload.

We can combine the `namespaceSelector` with the `podSelector` to make it more restrictive and only allow the sleep pods from the sleep namespace to access the httpbin pods:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-ingress-from-sleep
spec:
  podSelector:
    matchLabels:
      app: httpbin
  policyTypes:
  - Ingress
  ingress:
    - from:
      - namespaceSelector:
          matchLabels:
            kubernetes.io/metadata.name: sleep
        podSelector:
          matchLabels:
            app: sleep
```

This won't change anything for the requests made from the sleep pods, but it will restrict any other pods from making the requests.

To clean up, let's delete all policies:

```shell
kubectl delete netpol --all
```

Instead of blanket denial of all ingress (and egress) traffic within a namespace, a better starting point might be having a policy that automatically allows all traffic within a namespace, but denies all traffic from other namespaces:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-ns-ingress
spec:
  podSelector: {}
  ingress:
    - from:
      - podSelector: {}
```

We can achieve this using empty selectors, basically translating the above to "apply this policy to all pods in the namespace" and "allow all traffic from all pods in the namespace".

With the above policy in place, we can make requests within the default namespace, however, requests from the sleep namespace will be dropped.

## Egress policies

Since we only modified ingress policies, sending a request from pods works fine:

```shell
kubectl exec -it deploy/sleep -n default -- curl httpbin.default:8000/headers
```

Let's create a policy that denies all egress traffic from the sleep pod:

```yaml
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: deny-sleep-egress
spec:
  podSelector:
    matchLabels:
      app: sleep
  policyTypes:
  - Egress
```

If we retry the above request, it will fail. To allow the request, we can create a policy that allows all egress traffic from the sleep pod:

```yaml
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: deny-sleep-to-httpbin-egress
spec:
  podSelector:
    matchLabels:
      app: sleep
  policyTypes:
  - Egress
  egress:
    - to:
      - podSelector:
          matchLabels:
            app: httpbin
```

If you send a request you'll notice that it fails - why is that? The reason is DNS -- we're denying all egress traffic from the sleep pod, including requests to resolve the httpbin DNS name. We can confirm this by sending a request using the service IP instead of the DNS name and it will work.

To fix this, we can create a policy that allows all egress traffic for all pods in the namespace to the core-dns service running in the kube-system namespace:

```yaml
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: allow-dns-egress
spec:
  podSelector: {}
  policyTypes:
  - Egress
  egress:
    - to:
      - namespaceSelector:
          matchLabels:
            kubernetes.io/metadata.name: kube-system
        podSelector:
          matchLabels:
            k8s-app: kube-dns
```

This time, sending a request to the httpbin service will work. Cleanup the policies:

```shell
kubectl delete netpol --all
```

## L7 policies

We installed Prometheus and Grafana at the beginning, let's start sending some requests so we can see some data in Grafana.

```shell
kubectl exec -it deploy/sleep -n default -- /bin/sh

while true; do curl -X POST httpbin.default:8000/post; sleep 0.3; done & 
while true; do curl httpbin.default:8000/headers; sleep 0.3; done &
while true; do curl httpbin.default:8000/ip; sleep 0.3; done & 
```

Now let's say we only want to allow traffic to the `/headers` path and only for the GET method on the httpbin service. We know we can't use NetworkPolicy for that, but we can use CiliumNetworkPolicy.

Let's try sending a POST request and a request to /ip path just to show that it works:

```shell
kubectl exec -it deploy/sleep -n default -- curl -X POST httpbin.default:8000/post
kubectl exec -it deploy/sleep -n default -- curl httpbin.default:8000/ip
```

Both of these work, so let's deploy a policy that only allows GET requests to `/headers`:

```yaml
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: allow-get-headers
spec:
  endpointSelector:
    matchLabels:
      app: httpbin
  ingress:
  - fromEndpoints:
    - matchLabels:
        app: sleep
    toPorts:
    - ports:
      # Note port 80 here! We aren't talking about the service port, but the "endpoint" port, which is 80.
      - port: "80"
        protocol: TCP
      rules:
        http:
        - method: GET
          path: /headers
```

If we try sending a POST request, it will fail with HTTP 403, but GET requests to `/headers` will work.


The L7 policies are enforced by the Envoy proxy that spins up inside the Cilium agent pod.

Finally, we can also look at the metrics in Grafana. To access Grafana, we need to port-forward the service:

```shell
kubectl port-forward svc/grafana 3000:3000 -n cilium-monitoring
```