# wasm/rust-test

This is a quick test of Wasm for Istio.

## Building

To build the Wasm module, run:

```bash
./build.sh
```

Requirements:

- A *nix based system
- Docker

If you only need the .wasm file and not a docker image, run:

```bash
make build
cp target/wasm32-unknown-unknown/release/rust_test.wasm plugin.wasm
```

## Using

Apply the WasmPlugin yaml, for example:

```yaml
apiVersion: extensions.istio.io/v1alpha1
kind: WasmPlugin
metadata:
  name: httpbin-rust-test
  namespace: httpbin
spec:
  selector:
    matchLabels:
      app: httpbin
  match:
    - ports:
      - number: 80
  url: oci://docker.io/dhawton/wasm-rust-test:v1.0
```

Requests to httpbin's /get should now output HTTP status code 418 along with a teapot ASCII art if the WasmPlugin was successful.

## License

This project is licensed by [Apache 2.0](LICENSE)
