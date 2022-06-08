# Hoot Episode 27, May 24, 2022
Gloo Cilium and Istio Seamlessly Together

1 min hook to today's hoot.
Network policy are highly recommended along with L7 security policies based on Istio's security best practice doc. These resources are vastly different, how can we make this easier for our users? Welcome to hoot livestream episode 27, today we will discuss Cilium and Istio Seamlessly Together. I'm your host and look forward to learning this topic with you together.

**News (2 mins)**
Solo.io adds Cilium to Gloo Mesh:
https://www.devopsdigest.com/soloio-adds-cilium-to-gloo-mesh

https://www.solo.io/blog/enabling-cilium-gloo-application-networking-platform/

CNCF survey on SM: https://www.cncf.io/wp-content/uploads/2022/05/CNCF_Service_Mesh_MicroSurvey_Final.pdf

Envoy gateway: https://blog.envoyproxy.io/introducing-envoy-gateway-ad385cc59532

https://isovalent.com/blog/post/2022-05-16-tetragon

https://techcrunch.com/2022/05/18/apollo-graphql-launches-its-supergraph/

hoot update: 
Spire+Istio demo scripts are provided!

**speaker intro** (2 mins)
speakers: intro

General Questions (20 mins)

what is cilium?
- L3, IP based
- CNI, best CNI out there?

Can you describe Cilium's security model?

what are issues with network identity?
- overlapping IPs
- Does label help here?

How does network based identity work with multicluster?
- Is flat network required?
- access to k8s API server across multiclusters
https://docs.cilium.io/en/stable/gettingstarted/clustermesh/clustermesh/#limitations

Thoughts on how cilium and Istio be integrated together? Why would someone wants to use one or the other or both?

Why Cilium and Istio with Gloo mesh?

Gloo mesh has workspaces which provides multi-tenancy to different teams, are we also bringing tenancy to Cilium?

Can I use Gloo network for Cilium without Istio?

**Let us dive into demo** (5-10 mins)

Any demo you want to show?


**wrap up** (2 mins)
- Thank speakers! Ask speakers: How do folks reach out to you?
- Is this interesting? What other topics do you want to see to help you on your application networking? I am super grateful for everyone who liked our past hoot livestream and subscribed to our channel. Happy learning, and see you at the next episode!

