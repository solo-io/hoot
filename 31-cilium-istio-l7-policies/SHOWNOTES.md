# Hoot Episode 31, July 12, 2022
Cilium L7 Policies vs Istio's

1 min hook to today's hoot.

Istio offers rich L7 traffic management and security policies. Cilium also offers L7 policies and we have gotten a lot of questions from users if they still need Istio's L7 policies. In this livestream, Yuval will join Lin to explain and dive into the difference of the two focusing on security perspective.

**speaker intro** (2 mins)
speakers: intro

General Discussions (10-15 mins)

- What is Cilium L7 policy?
    - Doc: https://docs.cilium.io/en/stable/concepts/ebpf/intro/, search for L7 policy
    - Examples: https://docs.cilium.io/en/stable/policy/language/#layer-7-examples
    - I can use Kubernetes constructs for example SA in my network policy - example: https://docs.cilium.io/en/stable/policy/kubernetes/#serviceaccounts
    - Lin: a quick demo
    
- What is Istio's L7 policy?
    - Lin: explain authz policy: deny all then explicitly allow access, with a quick demo
    
- How do they compare?
    - problem with identity based on label or service account name
    - Encryption, mutual TLS?
        - Wireguard: https://docs.cilium.io/en/stable/gettingstarted/encryption-wireguard/
            - wireguard limitation: https://docs.cilium.io/en/stable/gettingstarted/encryption-wireguard/#limitations
        - ipsec: https://docs.cilium.io/en/stable/gettingstarted/encryption-ipsec/#encryption-ipsec
        - FIPS compliance: google search is wireguard fips compliant
    - Interoperatable
    - Can Envoy handle Multi-tenancy for L7?
    - Eventual consistency?

- Is there anything else you want to add?

- recap:  the CNI is responsible for L3/L4 traffic, and the service mesh for L7

**wrap up** (2 mins)
- Thank speakers! Ask speakers: How do folks reach out to you?
- Is this interesting? What other topics do you want to see to help you on your application networking? I am super grateful for everyone who liked our past hoot livestream and subscribed to our channel. Happy learning, and see you at the next episode!

Resources:
Yuval SMC EU 2022 slide: https://docs.google.com/presentation/d/1y7nTtpmSSJdeZrvFDibLJDEDMEfqRLKn5-DysGVc-Ws/edit#slide=id.g13698316493_2_447

Louis' tweet: https://twitter.com/louiscryan/status/1522661442138238976?s=20&t=nxDYj8oTdN7UHIZCOKqPRQ

Matt's viewpoint on multi-tenancy envoy: https://twitter.com/mattklein123/status/1522757356857085952?s=20&t=ACVDWbAoSYexcosvNXd3_Q

William Morgan's blog: https://buoyant.io/2022/06/07/ebpf-sidecars-and-the-future-of-the-service-mesh/

DPV2 & Cilium: https://www.doit-intl.com/ebpf-cilium-dataplane-v2-and-all-that-buzz-part-2/
