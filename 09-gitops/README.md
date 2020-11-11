Git ops!

Perform operations via Git. The idea is that as our infrastructure becomes 
declarative, we can configure it via git, getting the benefits of git:
- Tail log
- Ability to revert
- Ability to review

A purpose built tool for git-ops that we will use in our demo is [Flux](https://github.com/fluxcd/flux2).

# Demo
In this demo we will use a kind cluster and flux v2 command line tool.

In one terminal run `scripts/create_repo.sh`
In another terminal run `scripts/initcluster.sh`

Once cluster is initialized, run:
```
kubectl apply -f scripts/kustomization.yaml
```

You will now see flux installing the app. Track progress:
```
watch flux get kustomizations
```

See the app installed
```
kubectl get pods
``
