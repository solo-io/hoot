FROM rust:1.65 as builder

WORKDIR /usr/src/app
COPY . .
RUN rustup target add wasm32-unknown-unknown && \
    make build

FROM scratch
COPY --from=builder /usr/src/app/target/wasm32-unknown-unknown/release/rust_test.wasm ./plugin.wasm