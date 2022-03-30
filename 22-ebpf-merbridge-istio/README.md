# Hoot Episode 22 - Accelerate your service mesh with eBPF and Merbridge

## Recording ##
 https://youtu.be/r2wgInmsqsU

[show notes](SHOWNOTES.md)

[slides](merbridge.pdf)


## Hands-on: Steps from the demo

```bash
ps -ef | grep 3650d # grep the sleep process id
nsenter -t xxxx --net bash # enter the specific namespace, to do the iptables manipulation
# after entering to the sleep namespace, try the commands below to manipulate iptables
iptables-save > /tmp/x # save current iptables to a file
iptables -t nat -vnL # look at iptables rules
iptables -t nat -F # clean up iptables rules
iptables-restore < /tmp/x # restore the iptables rules



# The commands below are to see if eBPF progs or functions are present in the kernel.
sudo bpftool map
sudo bpftool prog
sudo bpftool cgroup tree


# The command to see the log, please make sure you add the -d arg to the merbridge daemonset, to get redirection fully logged.
sudo cat /sys/kernel/debug/tracing/trace_pipe

# Apply the merbridge resources
kubectl apply -f https://raw.githubusercontent.com/merbridge/merbridge/main/deploy/all-in-one.yaml

# Apply the sleep and the helloworld service
curl https://raw.githubusercontent.com/istio/istio/master/samples/sleep/sleep.yaml | kubectl apply -f -
curl https://raw.githubusercontent.com/istio/istio/master/samples/helloworld/helloworld.yaml| kubectl apply -f -




# The steps are:
# 1. Start kind cluster: 
kind create cluster --name merbridge

# 2. Start Istio
istioctl install -y

# 3. Start Merbridge
kubectl apply -f https://raw.githubusercontent.com/merbridge/merbridge/main/deploy/all-in-one.yaml

# 4. enable debug mode
kubectl edit ds -n istio-system merbridge
########
containers:
        - args:
        - /app/mbctl
        - -d
        - -m
        - istio
        - --ips-file
        - /host/ips/ips.txt
########

# 5. start to log debug info
sudo cat /sys/kernel/debug/tracing/trace_pipe

#6. try to make requests
kubectl exec sleep-xxxxxx-xxxx -it -- curl helloworld:5000/hello

#7. try to manipulate iptables rules of sleep
ps -ef | grep 3650d # grep the sleep process id
nsenter -t xxxx --net bash # enter the specific namespace, to do the iptables manipulation
# after entering to the sleep namespace, try the commands below to manipulate iptables
iptables-save > /tmp/x # save current iptables to a file
iptables -t nat -vnL # look at iptables rules
iptables -t nat -F # clean up iptables rules
iptables-restore < /tmp/x # restore the iptables rules

#8. to see if Merbridge is still working if the iptables rules are flushed. Or checkout the ISTIO_REDIRECT chain
```