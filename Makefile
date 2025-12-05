BUF ?= go run github.com/bufbuild/buf/cmd/buf@v1.33.0
PROTO_SRC ?= third_party/xai-proto

.PHONY: proto tidy ci

proto:
	$(BUF) generate --template buf.gen.yaml $(PROTO_SRC)

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
