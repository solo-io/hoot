.PHONY: run-servers
run-servers: run-processor run-httptarget run-envoy

.PHONY: restart-servers
restart-servers: restart-processor restart-httptarget restart-envoy

.PHONY: stop-servers
stop-servers: stop-processor stop-httptarget stop-envoy

.PHONY: run-processor
run-processor:
	go run ./pkg/processor -s /tmp/processor.sock -d 1>logs/processor.log 2>&1 &

.PHONY: stop-processor
stop-processor:
	kill -s SIGINT $(shell pidof processor)

.PHONY: restart-processor
restart-processor: stop-processor
restart-processor: run-processor

.PHONY: run-httptarget
run-httptarget:
	go run ./pkg/httptarget -p 9001 -d 1>logs/httptarget.log 2>&1 &

.PHONY: stop-httptarget
stop-httptarget:
	kill -s SIGINT $(shell pidof httptarget)

.PHONY: restart-httptarget
restart-httptarget: stop-httptarget
restart-httptarget: run-httptarget

.PHONY: download-envoy
download-envoy:
	([ ! -e ./envoy ] && \
		docker run -d --entrypoint /bin/sh --name get-envoy-binary envoyproxy/envoy:dev && \
		docker cp get-envoy-binary:/usr/local/bin/envoy ./envoy && \
		docker stop get-envoy-binary && \
		docker rm get-envoy-binary) || echo "envoy exists or copy from docker failed"

.PHONY: run-envoy
run-envoy: download-envoy
	mkdir -p ./logs
	./envoy -c ./envoy.yaml --component-log-level 'ext_proc:debug' 1>logs/envoy.log 2>&1 &

.PHONY: stop-envoy
stop-envoy:
	kill -s SIGINT $(shell pidof envoy)

.PHONY: restart-envoy
restart-envoy: stop-envoy
restart-envoy: run-envoy
