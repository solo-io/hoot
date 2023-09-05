# Hoot Episode 44 - Overview of SPIRE

## Recording ##
https://www.youtube.com/watch?v=MuYmhc4mJHI

## Demo

Prerequisites:
- Kubernetes cluster
- Istio 1.16.1

Make sure you update the YAML files with your own cluster name and trust domain you want to use. 

1. Update the `clusterName` and `trustDomain` fields in [`spire-controller-manager-config.yaml`](demo/spire-controller-manager-config.yaml)

1. Create the namespace and deploy the CSI driver, CRDs and the controller manager configuration and webhook:

```shell
demo/deploy-prereqs.sh
```

1. Open `spire-server.yaml` and update the following values in the `spire-server` ConfigMap resource:
    - `trust_domain`
    - `ca_subject`
    - cluster name in the `k8s_psat` NodeAttestor

1. Deploy the SPIRE server:

```shell
kubectl apply -f demo/spire-server.yaml
```

1. Open `spire-agent.yaml` and update the following values in the `spire-agent` ConfigMap resource:

    - `trust_domain`
    - cluster name in the `k8s_psat` NodeAttestor

1. Deploy the SPIRE agent:

```shell
kubectl apply -f demo/spire-agent.yaml
```

1. You can check the registration entries for SPIRE agents from the `spire-server` container:

```shell
SPIRE_SERVER_POD=$(kubectl get pod -l app=spire-server -n spire -o jsonpath="{.items[0].metadata.name}")

kubectl exec -it spire-server-0 -n spire -c spire-server -- ./bin/spire-server agent list  
```

1. Deploy the `clusterspiffeid.yaml` file:

```shell
kubectl apply -f demo/clusterspiffeid.yaml
```

1. Open `istio-spire-config.yaml` and update the `trustDomain` and the `clusterName` fields, then install Istio:

```shell
istioctl install -f demo/istio-spire-config.yaml
```

1. While Istio is being installed, patch the Ingress gateway deployment, so we get a SPIFFE ID for the ingress gateway:

```shell
kubectl patch deployment istio-ingressgateway -n istio-system -p '{"spec":{"template":{"metadata":{"labels":{"spiffe.io/spire-managed-identity": "true"}}}}}'
```

1. We can check the list of registration entries:

```shell
SPIRE_SERVER_POD=$(kubectl get pod -l app=spire-server -n spire -o jsonpath="{.items[0].metadata.name}")
kubectl exec -it spire-server-0 -n spire -c spire-server -- ./bin/spire-server entry show  
```

1. Label the default namespace for Istio sidecar injection:

```shell
kubectl label ns default istio-injection=enabled
```

1. Deploy the sleep workload (note that you have to add the annotation if you're using a different workload/YAML):

```shell
kubectl apply -f demo/sleep-spire.yaml
```

1. Retrieve the sleep workload certificate:

```shell
SLEEP_POD=$(kubectl get pod -l app=sleep -o jsonpath="{.items[0].metadata.name}")

istioctl proxy-config secret $SLEEP_POD -o json | jq -r '.dynamicActiveSecrets[0].secret.tlsCertificate.certificateChain.inlineBytes' | base64 --decode > sleep.pem
```

1. You can use `openssl` to inspect it:

```shell
openssl x509 -in sleep.pem -text
```

## Additional Resources

- SPIFFE spec - https://github.com/spiffe/spiffe/blob/main/standards/SPIFFE.md 
- Official SPIFFE docs - https://spiffe.io/docs/latest/spiffe-about/overview/ 
- SPIRE architecture & components - https://spiffe.io/docs/latest/spire-about/spire-concepts 
- Scaling SPIRE - https://spiffe.io/docs/latest/planning/scaling_spire 
- Istio cert management - https://istio.io/latest/docs/concepts/security/#pki 
- Istio SPIRE integration - https://istio.io/latest/docs/ops/integrations/spire 