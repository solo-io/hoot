# Hoot Episode 28, June 7, 2022
What is new in Istio 1.14

1 min hook to today's hoot.
Istio just turned 5 and the community released 1.14.  Wanna to know what is new with the release?  Today, I am so excited to discuss the Istio 1.14 release with Faseela!

**speaker intro** (2 mins)
speakers: intro

**News (2 mins)**

Istio 1.14 is out, thank RM & community!

Broadcom aquired VMware:
https://investors.broadcom.com/news-releases/news-release-details/broadcom-acquire-vmware-approximately-61-billion-cash-and-stock

IstioCon recap: https://www.youtube.com/watch?v=PyXxLXJRMoU

hoot update: https://github.com/solo-io/hoot/#upcoming-episodes

Cilium workshop: https://app.livestorm.co/solo-io/introduction-to-ebpf-and-cilium-amer-060922

General Questions (20 mins)

what is your contribution experience to Istio?

Discuss istio 1.14, Pull out release blog, release note and upgrade note

Discuss upgrade:

- First minor release without any upgrade caveat.
- Kubernetes warning of removal of deprecated APIs. https://kubernetes.io/blog/2021/07/14/upcoming-changes-in-kubernetes-1-22/#api-changes.  Istio prior to 1.10 won't work with k8s 1.22 or newer.

Discuss highlights of releases
- Spire: refer to episode 26
- Faseela: auto SNI
- Lin: min TLS version: https://preliminary.istio.io/latest/docs/tasks/security/tls-configuration/workload-min-tls-version/
-- Q: does it work for gateway?
- Lin: Telemetry API improvement

Other features of releases that are interesting?
- Faseela?
-- workload selector for DestinationRule
-- credential name support for sidecar egress TLS origination
- Lin?
-- PILOT_SEND_UNHEALTHY_ENDPOINTS
-- ProxyConfig - envoy runtime values: https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/operations/runtime
-- always disable protocol sniffing in production: **Fixed** an issue causing traffic from a gateway to a service with an [undeclared protocol]
-- Istio's default load balancing algorithm from `ROUND_ROBIN` to `LEAST_REQUEST`
-- **Added** support for WasmPlugin pulling image from private repository with `imagePullSecret`.
-- **Added** support of installing gateway helm chart as `daemonset`.
  ([Issue #37610](https://github.com/istio/istio/issues/37610))
-- anything interesting from istioctl?

**Let us dive into demo** (10-15 mins)
- Faseela: auto SNI etc
- Lin: upgrade to 1.14, min TLS, and telemetry API improvement


**wrap up** (2 mins)
- Thank speakers! Ask speakers: How do folks reach out to you?
- Is this interesting? What other topics do you want to see to help you on your application networking? I am super grateful for everyone who liked our past hoot livestream and subscribed to our channel. Happy learning, and see you at the next episode!

