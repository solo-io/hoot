# Intro

## Waypoint - what does it do?
- Takes care of the build / deploy / expose part of the software lifecycle

## Waypoint, Summary and My thoughts:

- Client / Server architecture.
- In kubernetes, it is installed as a StatefulSet with a data folder.
  I couldn't find guidance on HA Deployments; but I'm sure these will come as the project grows. State is stored in one file, and the server needs to be off to back it up. This is not very cloud native, but OK for a project at this early stage.
- I think it is trying to solve the problem that today we use Makefiles to hack around; provide a consistent way to build and deploy applications.
- Each deployment has a separate URL, this is usually done to separate the release phase from the deployment phase.
- It adds a component for the deployment, the waypoint entrypoint. This component communicates with the waypoint server. The server needs ample bandwidth to stream data from all the entrypoints.
The entrypoint allows the exec and logs commands to work and configuration of the application.
The entrypoint can be disabled, but then these functions won't work. [see more info here](https://www.waypointproject.io/docs/entrypoint/disable).

- While it supports a laptop deployment, laptop style dev-loop doesn't seem like the main use case (unlike for example: draft, Skaffold) as far as I can tell. Perhaps it is the most similar to Okteto. Target Use cases seem to be:
  - Deploy to test / prod cluster
  - Use in CI to create preview environments. 
I see it's strength where deployed resources need to be tracked, deployments are shared amount multiple people, and need to be cleaned up.


# Demo & Discussion

## Installation (one-time)
Install to kubernetes:
```
kubectl create namespace waypoint
waypoint install --platform=kubernetes -accept-tos -k8s-namespace=waypoint
```

## Build and deploy 
As simple as:
```
waypoint up
```

## Debugging / Logs:

Note that these actions below, don't use the kubernetes API. they use the waypoint API and reach the pod using the injected entry point. Which brings me to think about the need to add more security knobs, as kubernetes RBAC will not apply here. it essentially bypasses it.

If I'm understanding this correctly, this also means that to get the most value from waypoint, you **have** to build your application with waypoint. I hope they can come up with an easier way to accomplish this in kubernetes, as this will mean that we'll need to change existing build pipelines. Perhaps they could provide an `inject` to modify existing containers?

```
waypoint exec /bin/bash
# you can use `ss -p` to see the waypoint entrypoint connecting to the waypoint server. waypoint server must be reachable to the entrypoint.
waypoint logs
```

# Final thoughts

Waypoint seems at a very early stage. I think it generally makes sense to abstract kubernetes from 
the developers with a simpler experience. I.e. `waypoint up` is easier to learn than `helm install release-name repo/chart-name --set randomvalue=foo`. `waypoint logs` is easier than `kubectl logs -n ns deploy/app-v2`. It does have a server that requires care and feeding, and doesn't seem to currently have means
of more sophisticated access control.

# More Resources:
- The announcement: https://www.youtube.com/watch?v=nasVKN7Wbtk
- Demo: https://www.youtube.com/watch?v=azoQYaJsxGk
- Docs: https://www.waypointproject.io/docs/
- My favorite part of the docs, the internals: https://www.waypointproject.io/docs/internals/architecture
- Review by the DevOps toolkit: https://www.youtube.com/watch?v=7qrovZjdgz8 (I don't nessecarly agree with his critiques)