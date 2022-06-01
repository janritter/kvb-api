clean:
	rm -rf dist

prepare:
	go mod download

build: clean prepare
	go build -o dist/kvb-api

run:
	go run main.go

run-with-tracing:
	GRPC_GO_LOG_VERBOSITY_LEVEL=99 GRPC_GO_LOG_SEVERITY_LEVEL=info ENABLE_TRACING=true OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4317" go run main.go

build-docker: clean
	mkdir -p dist/
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist -ldflags="-w -s -extldflags '-static'" ./...
	docker buildx build --platform=linux/amd64 -t kvb-api .
