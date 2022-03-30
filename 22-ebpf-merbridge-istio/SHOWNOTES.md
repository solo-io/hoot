# Hoot Episode 22, March 29, 2022

1 min hook to today's hoot.
Service mesh & eBPF - these are super hot topics, are they competing, are they complementary to each other?  Can you add ebpf support to Istio without changing Istio at all?

**speaker intro** (2 mins)
Welcome to hoot livestream, where we bring istio, envoy, k8s, ebpf & graphql technologies to you so you can be well prepared at your job, be the best cloud native developer/operator/architect!
Lin: your host for hoot livestream today.
speakers: intro

**News (2 mins)**

Episode 21 steps are published: https://github.com/solo-io/hoot/tree/master/21-istio-in-action-book

IstioCon and SMC EU speakers are notified, congratulations!

https://istio.io/latest/blog/2022/istioctl-proxy/

https://venturebeat.com/2022/03/21/report-89-of-orgs-have-been-attacked-by-kubernetes-ransomware/

https://www.businesswire.com/news/home/20220323005341/en/Spectro-Cloud-Announces-T-Mobile-Ventures-Investment-in-its-Series-B-Funding-Round-to-Drive-Innovation-in-Kubernetes-Management-at-5GEdge-Locations


**General Questions** (5 mins)
DaoCloud gained a steering contribution seat, congrats! 

- What is DaoCloud?

- Why Istio for DaoCloud?

- What triggered you to start merbridge project?

** merbridge ** (10 mins)

https://istio.io/latest/blog/2022/merbridge/

- How merbridge works without modifying Istio?
-- what is the role of init container?
-- how does it work with istio CNI? maybe no need for istio CNI?

- Can you explain how users can accelerate SM adoption with eBPF and merbridge?

- How does merbridge work for services in the service mesh vs not in the service mesh?

- How do I know merbridge is working?

- From your performance test, i think it is about 10% latency improvement. are those for same node or different nodes?
-- if you could modify Istio, would performance be better?  for example, init container may not be needed?

- open to share challenges when building merbridge?

- We've got folks in the community exploring merbridge.  Is merbridge ready for production?

**Let us dive into merbridge demo** (5-10 mins)

**wrap up** (2 mins)
- Thank speakers! Ask speakers: How do folks reach out to you?
-- https://join.slack.com/t/merbridge/shared_invite/zt-11uc3z0w7-DMyv42eQ6s5YUxO5mZ5hwQ
-- https://github.com/merbridge/merbridge
- Is this interesting? What other topics do you want to see to help you on your application networking? Remind folks to comment, like and subscribe. See you next Tues (increase application resiliency with spot vms)!