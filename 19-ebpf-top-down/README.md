# Overview

eBPF is a Linux kernel technology that is quickly growing in popularity as it provides developers with the ability to inject custom logic into running kernels in a safe and efficient way.

In the application networking space there are a few common use cases such as:
 * tracing/observability
 * security enforcement
 * network acceleration

Due to the breadth of eBPF as a technology it can be challenging to learn and get started.
The slides/video give a quick overview of eBPF as well a top-down view of a specific example -- tracing network connections.
The example was built and run via our open-source tool [BumbleBee](https://bumblebee.io/).

# Example
A [simple example](probe.c) is included that is based off a `bee init` template.

It can be built and run easily via `bee`, e.g.
```bash
bee build probe.c my_probe:v1
bee run my_probe:v1
```

See the [BumbleBee getting started guide](https://github.com/solo-io/bumblebee/blob/main/docs/getting_started.md) for more info.


# References
* https://ebpf.io/
* https://github.com/iovisor/bcc
* https://github.com/libbpf/libbpf
* https://nakryiko.com/posts/libbpf-bootstrap/
* https://github.com/solo-io/bumblebee
