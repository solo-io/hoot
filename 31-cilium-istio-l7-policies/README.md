# Hoot Episode 30 - HTTP/3 With Envoy Explained

## Recording ##
 https://youtu.be/d4mEnVZKum8

[show notes](SHOWNOTES.md)

## Hands-on: Steps from the demo

### Cilium Demo

Setup Cilium CNI based on latest stable doc.  Check status:

```console
cilium status --wait
```

Install apps:

```console
kubectl apply -f samples/sleep.yaml
kubectl apply -f samples/helloworld-with-affinity.yaml
kubectl apply -f samples/notsleep.yaml
```

Call sleep to helloworld:

```console
kubectl exec -it $(k get po -lapp=sleep -ojsonpath='{.items[0].metadata.name}') -- curl helloworld:5000/hello
```

Call notsleep to helloworld:

```console
kubectl exec -it $(k get po -lapp=notsleep -ojsonpath='{.items[0].metadata.name}') -- curl helloworld:5000/hello
```

Apply cilium l4 policy:

```console
kubectl apply -f cilium-policy-l4.yaml
```

Call sleep to helloworld:

```console
kubectl exec -it $(k get po -lapp=sleep -ojsonpath='{.items[0].metadata.name}') -- curl helloworld:5000/hello
```

Call notsleep to helloworld:

```console
kubectl exec -it $(k get po -lapp=notsleep -ojsonpath='{.items[0].metadata.name}') -- curl helloworld:5000/hello
```

```console
kubectl get ciliumendpoint,ciliumidentity
```

Apply cilium l7 policy:

```console
kubectl apply -f cilium-policy-l7.yaml
kubectl exec -it $(k get po -lapp=sleep -ojsonpath='{.items[0].metadata.name}') -- curl helloworld:5000/hello
kubectl exec -it $(k get po -lapp=sleep -ojsonpath='{.items[0].metadata.name}') -- curl helloworld:5000/hello3
kubectl exec -it $(k get po -lapp=notsleep -ojsonpath='{.items[0].metadata.name}') -- curl helloworld:5000/hello
```

Set the cilium pod:

```console
node=kind-worker
pod=$(kubectl -n kube-system get pods -l k8s-app=cilium -o json | jq -r ".items[] | select(.spec.nodeName==\"${node}\") | .metadata.name" | tail -1)
kubectl -n kube-system exec -q ${pod} -- ps -ef
```

Dump envoy config from the Cilium pod:

```console
kubectl -n kube-system exec -q $pod -- apt update
kubectl -n kube-system exec -q $pod -- apt -y install curl
kubectl -n kube-system exec -q $pod -- curl -k -L https://github.com/fullstorydev/grpcurl/releases/download/v1.8.6/grpcurl_1.8.6_linux_x86_64.tar.gz --output /tmp/grpcurl.tar.gz
kubectl -n kube-system exec -q $pod -- bash -c "cd /tmp && tar zxvf /tmp/grpcurl.tar.gz"
#dump envoy config
k exec -n kube-system -it $pod -- curl -s --unix-socket /var/run/cilium/envoy-admin.sock http://localhost/config_dump
# search for cilium.l7policy ^^
```

Sending an xDS request to the control plane:

```console
while true; do
cat <<EOF | kubectl -n kube-system exec -q -i $pod -- /tmp/grpcurl -d @ -plaintext -unix /var/run/cilium/xds.sock cilium.NetworkPolicyDiscoveryService/StreamNetworkPolicies > /tmp/output
{
  "node": $(kubectl -n kube-system exec -q $pod -- curl -s --unix-socket /var/run/cilium/envoy-admin.sock http://localhost/config_dump | jq -r ".configs[0].bootstrap.node"),
  "resourceNames": [
  ],
  "typeUrl": "type.googleapis.com/cilium.NetworkPolicy"
}
EOF
if [ -s /tmp/output ]; then
  cat /tmp/output
  break
fi
done
```
### Istio Demo

Install the latest stable Istio and enable mTLS for the namespace:

```console
kubectl apply -f - <<EOF
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
  namespace: default
spec:
  mtls:
    mode: STRICT
EOF
```

Deploy the pods to the mesh:

```console
kubectl label namespace default istio-injection=enabled
kubectl apply -f samples/sleep.yaml
kubectl apply -f samples/helloworld.yaml
kubectl apply -f samples/notsleep.yaml
```

Show the certs:

```console
kubectl exec -it $(k get po -lapp=sleep -ojsonpath='{.items[0].metadata.name}') -- openssl s_client -connect helloworld:5000 -showcerts
```

Step through the certs:

```console
echo "-----BEGIN CERTIFICATE-----
<--copy the server certificate here-->
-----END CERTIFICATE-----" | step certificate inspect -
```

Apply RBAC rules:

```console
kubectl apply -f - <<EOF
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: allow-nothing
  namespace: default
spec:
  {}
EOF

kubectl exec -it $(k get po -lapp=sleep -ojsonpath='{.items[0].metadata.name}') -- curl helloworld:5000/hello

kubectl apply -f - <<EOF
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: "helloworld-viewer"
  namespace: default
spec:
  selector:
    matchLabels:
      app: helloworld
  action: ALLOW
  rules:
  - from:
    - source:
        principals: ["cluster.local/ns/default/sa/sleep"]
    to:
    - operation:
        methods: ["GET"]
EOF
```

Check sleep to helloworld with the newly applied RBAC rules:

```console
kubectl exec -it $(k get po -lapp=sleep -ojsonpath='{.items[0].metadata.name}') -- curl helloworld:5000/hello

kubectl exec -it $(k get po -lapp=notsleep -ojsonpath='{.items[0].metadata.name}') -- curl helloworld:5000/hello
```