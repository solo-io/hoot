To perform hitless deploys, each component in our infrastructure needs to be aware of the 
health status of the component in the next layer.

In the context of envoy, this means that:
- The cloud load balancer needs to know the health state of envoy
- Envoy needs to know the health state of the services it routes to

# LB --> Envoy
Configure the envoy health check filter.
Before existing envoy, fail health checks (`POST /healthcheck/fail` on the admin page)


# Envoy -> LB

Configure active or passive health checks, and retries.

# Caveats

If exposing envoy as a NodePort k8s service, then using standard health from the cloud load balancer
is probably something you want to avoid. The load balancer will send health checks to each k8s node.
The node in turn, will send the request to a random pod, resulting in inconsistent health info given
to the cloud load balancer.

Note that envoy is not aware of k8s readiness/liveness probes. Either have your control plane propagate
this info to envoy via EDS, or configure separate health checks on envoy, regardless of the k8s probes.

Also note, that in distributed systems, each component has an eventual consistent state.
This means that when you want to remove a pod, you want to give the pod some time to drain requests.
During this time, the pod should fail health checks, and the control plane should remove it from envoy
This gives enough time for components sending traffic to the pod (i.e. envoy) to reconcile their state 
and stop sending traffic without disruption.

# Demo

## Envoy to upstream
Run the xds server:
```
(cd xds; go run xds.go)
```

Run envoy:
```
envoy -c envoy.yaml
```

See failing requests stats

```
curl -s http://localhost:8001/stats | grep listener.0.0.0.0_8000.http.ingress_http.downstream_rq_5xx
```

send requests:
```
while true; do hey -n 100 http://localhost:8000/ ; sleep 1;done
```

## LB to envoy

Check health:
```
curl -v http://localhost:8000/health
```

Fail health checks:
```
curl -XPOST http://localhost:8001/healthcheck/fail
```


# More resources:

https://www.envoyproxy.io/docs/envoy/latest/faq/load_balancing/transient_failures.html
https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/upstream/health_checking
https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/health_check_filter