# Open Policy Agent
- General purpose policy language and a runtime to evaluate it.
- It can run as a go library, as a sidecar with REST api, and (WIP) compiled to WASM (note, this is not the same as the envoy filter wasm and can **not** be used in envoy directly).
- In the context of APIs, The goal of OPA is to get a decision based off of input. i.e. Should we allow a given request?

Some concepts (adapted for API usage):
- Input (json) - the request we want to run our policy on, and make a decision.
- Data (json) - Data that can be used on the policy that is not request specific (for example roles in the organization).
- Policy / Modules (Rego) - the logic of the policy. Can use information from the input and the data and returns an answer (with api requests this is usually true or false; but it doesn't have to be)

I want to evaluate this request (input) using this policy (rego) under this organization data (Data).

In the context of API this can mean: Only allow a POST request IF the user has a proper role in the organization.
The input will have the current user. The data can provide user to role map

# rego
Rego - query language inspired by datalog (which is a subset of prolog)

- Language to evaluate policy
- Order usually doesn't matter

# Examples
--------
Examples with this [data](basic-data.json). Run:
```
opa run basic-data.json
```
And try these:

- All solo.io members: `data.orgs["solo.io"].members`
- Is there a solo.io member named yuval: `data.orgs["solo.io"].members[_] == "yuval"`
- All the orgs yuval is a member of: `some org; data.orgs[org].members[_] == "yuval"`



# Policy examples

given data.json and query.rego, let's we evaluate input.json in 3 different ways:


## server
Use opa in server mode:

```
opa run data.json query.rego -s
curl localhost:8181/v1/data/app/api/allow -d@input.json
```

Note: in prod you will probably use bundles and health checks to determine readiness - this is just a demo!


## command line

```
cat input.json|jq .input|opa eval -I -d data.json -d query.rego data.app.api.allow
```

## unit tests

```
opa test .
```
