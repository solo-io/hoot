# Hoot Episode 23, April 5, 2022
HOW to Increase service resiliency for services running on Kubernetes on Spot VMs

https://github.com/murphye/cheap-gke-cluster
https://thenewstack.io/run-a-google-kubernetes-engine-cluster-for-under-25-month/

1 min hook to today's hoot.

Running your app in cloud? Is cost a concern for you? Is high resiliency app a priority for you? This episode, we will introduce methods to increase resiliency for your app, while on spot vms to cut costs!

**speaker intro** (2 mins)
Welcome to hoot livestream, where we bring istio, envoy, k8s, ebpf & graphql technologies to you so you can be well prepared at your job, be the best cloud native developer/operator/architect!
Lin: your host for hoot livestream today.
speakers: intro

**News (2 mins)**

Episode 22 steps are published: https://github.com/solo-io/hoot/tree/master/22-ebpf-merbridge-istio

https://dagger.io/blog/public-launch-announcement

https://techcrunch.com/2022/03/31/as-docker-gains-momentum-it-hauls-in-105m-series-c-on-2b-valuation/

https://twitter.com/evan2645/status/1509607415011954690?s=20&t=u9fkrb4wjtdDHbiQHCeijg

IstioCon schedule out!

https://lp.solo.io/devopsdays-raleigh-networking

**General Questions** (5 mins)

What is spot VM?

Why spot VM for average folks?
- Why spot VM is perfect for development and testing of your services?

I am convinced Spot VMs are interesting, how to
to increase resiliency?
- Retry
- Replica numbers
- Anti-Affinity

What if I have larger clusters? Can I continue to use this?

Is spot VM avail only for Google cloud?

Strategy to cut cost and ensure app healthy all the time?
- Should I consider a mixed of regular VM and spot VM?

**Let us dive into  demo** (5-10 mins)

- gloo edge VS: can you show envoy config there?


Any other tips you would like to share before we wrap up?

**wrap up** (2 mins)
- Thank speakers! Ask speakers: How do folks reach out to you?
- Is this interesting? What other topics do you want to see to help you on your application networking? Remind folks to comment, like and subscribe. See you next next Tues (Debug Envoy Configs and Analyze Access Logs)!