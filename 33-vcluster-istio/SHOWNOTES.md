# Hoot Episode 33, July 26, 2022
Speed your Istio development environment with vcluster

Do you have concerns with cost, CPU and networking for your Kubernetes cluster? What is a cluster within a cluster? In this livestream, Fabian and Rich from Loft Labs along with Antonio from Solo.io will join Lin to discuss what, why and when to use vCluster and live demo how to speed your Istio (or others) dev environment with vCluster to ease your cost, CPU and networking concerns.

**speaker intro** (2 mins)
speakers: intro

**News (2 mins)**
Cilium 1.12: https://isovalent.com/blog/post/cilium-service-mesh/?utm_campaign=1.12-release
- "This will make the Cilium Service Mesh data plane compatible with the service meshes such as Istio which are already migrating to Gateway API."
- call out Istio portion: Istio is the existing service mesh control plane that is supported. It currently requires to be run with the sidecar-based datapath.
- call out mixed mode: With Cilium Service Mesh, you have both options available in your platform and can even run a mix of the two
My blog: https://www.cncf.io/blog/2022/07/22/exploring-cilium-layer-7-capabilities-compared-to-istio/

General Discussions (10-15 mins)
- vcluster team: What exactly is vcluster? Is it free, what license?
- vcluser team: Why vcluster?
- vcluster team: How does vcluster work?
- vcluster team and solo: How do I get started with learning vcluster?

* a quick demo vcluster + Istio in local env *

- Does it work with multiclusters?

* a quick demo vcluster + Istio  multiclusters in local env *

- all: When and when not to use vcluster? 
- vcluster team: Can I use vCluster as my team boundry when the boundry is more than 1 namespaces? Also, does it offer better isolation than namespace?
    - Could this be useful for Istio multi-cluster in prod?

**wrap up** (2 mins)
- Thank speakers! Ask speakers: How do folks reach out to you?
- Is this interesting? What other topics do you want to see to help you on your application networking? I am super grateful for everyone who liked our past hoot livestream and subscribed to our channel. Happy learning, and see you at the next episode!

Resources:
https://loft.sh/blog/development-environments-with-vcluster-a/
https://www.vcluster.com/
https://istio.io/latest/docs/setup/install/multicluster/multi-primary/
