# 🎉Hoot Episode 50 - Kubernetes Networking and Cilium

## Recording
https://www.youtube.com/watch?v=TxPEpArfjwY&t=5s

## Demo - Kind with iptables

1. Create a Kubernetes cluster with `kind` using iptables mode (default):

  ```shell
  kind create cluster --config=kind-iptables.yaml
  ```

2. Deploy `httpbin` and scale it up to 3 replicas:

  ```shell
  kubectl apply -f https://raw.githubusercontent.com/istio/istio/master/samples/httpbin/httpbin.yaml
  kubectl scale deployment httpbin -n default --replicas=3
  ```

We deployed `httpbin` workload (3 replicas) and a corresponding Kubernetes service. Let's look at the iptables rules that are set up on one of the cluster nodes:

>You can use `docker exec` to go to one of the nodes and run the `iptables` command.

```shell
iptables -L -t nat > iptables.txt
```

From the `PREROUTING` chain, traffic is routed to `KUBE-SERVICES` chain. The `KUBE-SERVICES` chain contains rules for each service in the cluster. For example:

```console
Chain KUBE-SERVICES (2 references)
target                     prot opt source               destination
KUBE-SVC-FREKB6WNWYJLKTHC  tcp  --  anywhere             10.96.120.233        /* default/httpbin:http cluster IP */ tcp dpt:8000
KUBE-SVC-NPX46M4PTMTKRN6Y  tcp  --  anywhere             10.96.0.1            /* default/kubernetes:https cluster IP */ tcp dpt:https
KUBE-SVC-ERIFXISQEP7F7OF4  tcp  --  anywhere             10.96.0.10           /* kube-system/kube-dns:dns-tcp cluster IP */ tcp dpt:domain
KUBE-SVC-JD5MR3NA4I4DYORP  tcp  --  anywhere             10.96.0.10           /* kube-system/kube-dns:metrics cluster IP */ tcp dpt:9153
KUBE-SVC-TCOU7JCQXEZGVUNU  udp  --  anywhere             10.96.0.10           /* kube-system/kube-dns:dns cluster IP */ udp dpt:domain
KUBE-NODEPORTS             all  --  anywhere             anywhere             /* kubernetes service nodeports; NOTE: this must be the last rule in this chain */ ADDRTYPE match dst-type LOCAL
```

The rule applies to any incoming source IP address - i.e. we don't care where the traffic coming from. The destination IP `10.96.120.233` is the IP address of the Kubernetes service (in this case, `httpbin` in the default namespace). The last portion, `tcp dpt:8000`, specifies that this rule applies to TCP traffic destined for port 8000.

In short - the rule is saying that any traffic going to the `httpbin.default` service should be processed by the `KUBE-SVC-FREKB6WNWYJLKTHC` chain. Let's look at that chain:

```console
Chain KUBE-SVC-FREKB6WNWYJLKTHC (1 references)
target                     prot opt source               destination
KUBE-MARK-MASQ             tcp  --  !10.244.0.0/16       10.96.120.233        /* default/httpbin:http cluster IP */ tcp dpt:8000
KUBE-SEP-PFKDKMHHACIMKSSX  all  --  anywhere             anywhere             /* default/httpbin:http -> 10.244.1.2:80 */ statistic mode random probability 0.25000000000
KUBE-SEP-Q7UD3MC3WPZFKDWM  all  --  anywhere             anywhere             /* default/httpbin:http -> 10.244.1.3:80 */ statistic mode random probability 0.33333333349
KUBE-SEP-66R3JWVSZK6BYGSL  all  --  anywhere             anywhere             /* default/httpbin:http -> 10.244.2.2:80 */ statistic mode random probability 0.50000000000
KUBE-SEP-XIWLPFJKMVHRQV3W  all  --  anywhere             anywhere             /* default/httpbin:http -> 10.244.2.3:80 */
```

The first rule in the chain marks all packets that are not from the `10.244.0.0/16` source subnet (i.e. they aren't originating from within the cluster), going to the `httpbin` service, for masquerading (`KUBE-MARK-MASQ` chain). The `KUBE-MARK-MASQ` chain is responsible for marking packets that need to be masqueraded - the chain marks all packets with `0x4000`:

```console
Chain KUBE-MARK-MASQ (16 references)
target     prot opt source               destination
MARK       all  --  anywhere             anywhere             MARK or 0x4000
```

>What is masquerading? Masquerading is a form of SNAT (Source Network Address Translation) used when the source address for the outbound packets should be changed to the address of the outgoing networking interface. For example, when packets are exiting the node and instead of using the internal IP, an external IP should be used.

The other targets in the `KUBE-SVC` chain correspond to the pods that are backing the service. For example, here's one of the chains for the `httpbin` pods:

```console
Chain KUBE-SEP-PFKDKMHHACIMKSSX (1 references)
target          prot opt source               destination
KUBE-MARK-MASQ  all  --  10.244.1.2           anywhere             /* default/httpbin:http */
DNAT            tcp  --  anywhere             anywhere             /* default/httpbin:http */ tcp to:10.244.1.2:80
```

We've seen the first rule already - the `KUBE-MARK-MASQ` rule applies for packets that are going out of the pod (source `10.244.1.2`) and marks the outgoing packets with `0x4000`.

The second rule does the Destination Network Address Translation (DNAT). This rule translates the original destination IP address (IP of the `httpbin` service) to the IP address of the pod (`tcp to:10.244.1.2:80`). This is the rule that does the redirection to the actual pod IP.


## Demo - Kind using IPVS

1. Setup the kind cluster using IPVS:

```shell
kind create cluster --config=kind-ipvs.yaml
```

2. Deploy `httpbin` and scale it up to 3 replicas:

```shell
kubectl apply -f https://raw.githubusercontent.com/istio/istio/master/samples/httpbin/httpbin.yaml
kubectl scale deployment httpbin -n default --replicas=3
```

Let's look at the iptables rules this time (run the command from one of the nodes):

```shell
iptables -L -t nat > iptables-ipvs.txt
```

When using the ipvs here's how the `KUBE-SERVICES` chain looks like:

```console
Chain KUBE-SERVICES (2 references)
target          prot opt source               destination         
KUBE-MARK-MASQ  all  --  !10.244.0.0/16       anywhere             /* Kubernetes service cluster ip + port for masquerade purpose */ match-set KUBE-CLUSTER-IP dst,dst
KUBE-NODE-PORT  all  --  anywhere             anywhere             ADDRTYPE match dst-type LOCAL
ACCEPT          all  --  anywhere             anywhere             match-set KUBE-CLUSTER-IP dst,dst
```

Note that previously, the chain had a target for each service in the cluster. The difference this time is in the use of the `KUBE-CLUSTER-IP` IP set.

IP set is an extension to the iptables and it allows us to match multiple IP addresses at once (i.e. match against a group). The IP set `KUBE-CLUSTER-IP` contains the Kubernetes service IP addresses. We can use the `ipset` tool to list the contents of the set:

```shell
apt-get update -y
apt install -y ipset
```

Let's look at the contents of the `KUBE-CLUSTER-IP` set:

```shell
ipset list KUBE-CLUSTER-IP
```

```console
Name: KUBE-CLUSTER-IP
Type: hash:ip,port
Revision: 6
Header: family inet hashsize 1024 maxelem 65536 bucketsize 12 initval 0xa3af5d2d
Size in memory: 440
References: 2
Number of entries: 5
Members:
10.96.0.10,tcp:9153
10.96.0.10,udp:53
10.96.0.1,tcp:443
10.96.0.10,tcp:53
10.96.54.116,tcp:8000
```

Note the last IP (and port) corresponds to the `httpbin` Kubernetes service. While IPs where stored inside the iptables rules in the previous mode, when using ipvs, the IPs are stored in the ipset and the iptables rules now only reference the ipset. This allows us to only modify the ipsets, instead of traversing the iptables chains and modifying the rules.

Let's install the `ipvsadm` tool, so we can look at the configuration of the IPVS proxy:

```shell
apt-get update -y
apt-get install -y ipvsadm
```

If we run `ipvsadm -L -n` we get the following output:

```console
IP Virtual Server version 1.2.1 (size=4096)
Prot LocalAddress:Port Scheduler Flags
  -> RemoteAddress:Port           Forward Weight ActiveConn InActConn
TCP  10.96.0.1:443 rr
  -> 172.18.0.4:6443              Masq    1      0          0
TCP  10.96.0.10:53 rr
  -> 10.244.0.2:53                Masq    1      0          0
  -> 10.244.0.4:53                Masq    1      0          0
TCP  10.96.0.10:9153 rr
  -> 10.244.0.2:9153              Masq    1      0          0
  -> 10.244.0.4:9153              Masq    1      0          0
TCP  10.96.54.116:8000 rr
  -> 10.244.1.2:80                Masq    1      0          0
  -> 10.244.1.3:80                Masq    1      0          0
  -> 10.244.2.2:80                Masq    1      0          0
  -> 10.244.2.3:80                Masq    1      0          0
UDP  10.96.0.10:53 rr
  -> 10.244.0.2:53                Masq    1      0          0
  -> 10.244.0.4:53                Masq    1      0          0
```

Note the IP `10.96.54.116:8000` that corresponds to the httpbin service and the lines that follow are the IP addresses of the backing pods.

One of the advantages that ipvs has over specifying the individual iptables rules is that it can do proper load balancing - the `rr` in the output stands for round-robin, where as iptables can't do load balancing and it uses probabilities to distribute the traffic.

Also, when we scale up/down the deployments in the pure iptables mode, the rules get added to the chain and have to be processes sequentially. When using IPVS the iptables rules and chains stay the same and the IPVS proxy is updated with the new IP addresses of the pods. This is much more efficient and allows us to scale up/down the deployments without having to modify the iptables rules.