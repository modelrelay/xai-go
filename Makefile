BUF ?= go run github.com/bufbuild/buf/cmd/buf@v1.71.0
PROTO_SRC ?= third_party/xai-proto
# Public scope only: generate the xAI inference API (xai/api), not xAI's
# internal management_api / shared (billing, analytics) surfaces.
PROTO_PATH ?= $(PROTO_SRC)/proto/xai/api

.PHONY: proto tidy ci

proto:
	$(BUF) generate --template buf.gen.yaml --path $(PROTO_PATH) $(PROTO_SRC)

tidy:
	go fmt ./...
	go mod tidy

ci:
	@files=$$(gofmt -l .); if [ -n "$$files" ]; then \
		echo "gofmt check failed:"; \
		echo "$$files"; \
		exit 1; \
	fi
	go test ./...
	go vet ./...
	govulncheck ./...
	$(BUF) lint $(PROTO_SRC)
