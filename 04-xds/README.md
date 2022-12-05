# Setup

## Envoy Terminal
In the background, start:
```
go run server.go&
envoy -c xds.yaml
```

## XDS Terminal:
Change into the xDS package:
```
cd xds
```
Start XDS server:
```
go run xds.go
```

This server demonstrates dynamic configuration using traffic shifting based on user input. You can also test how envoy behaves when given invalid configuration, by providing input that is out of range.

## Request Terminal:
generate traffic, and see stats:
```
while true; do hey -n 100 -c 5 -t 1  http://localhost:10000/ ; sleep 1;done
```

Red server returns response code 200,
Blue server returns response code 201.

# Self-demo

On the XDS server terminal, type in traffic weight (0-100 are valid numbers) and see the traffic shifting live in the second terminal (where the `hey` loop runs).

Also try providing an invalid configuration (i.e. weight of 150), and see that envoy
does not change it's behavior.

Stats related to the control plane are:
- `control_plane.connected_state`
- stats that end with `.version_text`

Specifically, when you provide an invalid weight, you notice that the `http.http.rds.local_route.version_text`
does not match up with the server's published snapshot version.

Note that when envoy nacks:
- it retains the last known good configuration
- the version text in the stats will not match the one from the server
