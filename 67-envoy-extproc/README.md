This is demonstration code used for [Hoot 67](https://www.linkedin.com/events/7160706454570577920/comments/)

Most of this is adapted from https://github.com/GoogleCloudPlatform/envoy-processor-examples

Notable adjustments from the GCP example:
- Processing server uses a unix socket instead of TCP
- A few additional processing targets were added to the processing server


The three servers can be run by `make run-servers` which will launch the servers in the shell and output logs into the `./logs` directory.

Some of the examples used in the demo:

The following will add headers to the request AND response with the hash of the body
```shell
curl -v -X POST -d 'hello, world' http://127.0.0.1:10000/echohashbuffered
```
The following will block request since it comes from source address beginning `127`
```shell
curl -v http://127.0.0.1:10000/blockLocalHost
```
The following will set dynamic metadata in the filter which is seen by logging it in the access log
```shell
curl -v http://127.0.0.1:10000/dynamicMetadata -H 'x-set-metadata: scoobydoo'
