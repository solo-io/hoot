# Examples
Use istioctl analyze to detect issues

```
istioctl analyze
```

Pilot's ControlZ interface:

```bash
kubectl port-forward -n istio-system deploy/istiod 9876
```

Or

```bash
istioctl dashboard controlz deployment/istiod.istio-system
```

## Envoy admin page:

```bash
kubectl port-forward -n default deploy/reviews-v1 15000 &
```

## Envoy bootstrap config

```bash
kubectl exec -n default deploy/reviews-v1 -c istio-proxy -- cat /etc/istio/proxy/envoy-rev0.json
```

Or:

```bash
istioctl proxy-config bootstrap -n default deploy/reviews-v1
```

## Envoy in sync:

```bash
istioctl proxy-status
```

Or see if update_rejected gets incremented for a pdo:

```bash
curl localhost:15000/stats | grep update_rejected
```

## Envoy config dump

```bash
curl localhost:15000/config_dump
```

Or

```bash
istioctl proxy-status deploy/reviews-v1
istioctl proxy-config cluster deploy/reviews-v1
istioctl proxy-config route deploy/productpage-v1
istioctl proxy-config listener deploy/details-v1
```

# Debug logs

Control plane, through the ControlZ interface

Data plane:

```bash
curl 'localhost:15000/logging?level=debug'
kubectl logs -n default deploy/productpage-v1 -c istio-proxy
```



# Further reading

- https://istio.io/latest/docs/ops/diagnostic-tools/
- https://istio.io/latest/docs/ops/common-problems/