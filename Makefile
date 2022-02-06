prepare:
	go mod download

build: prepare
	go build -o bin/kvb-api

run:
	go run main.go

start-jaeger:
	docker-compose up -d

run-with-tracing:
	ENABLE_TRACING=true OTEL_EXPORTER_OTLP_ENDPOINT="localhost:4317" go run main.go
