build:
	cargo build --target wasm32-unknown-unknown --release
build-push:
	docker build -t dhawton/wasm-rust-test:v1 .
	docker push dhawton/wasm-rust-test:v1