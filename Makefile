prepare:
	go mod download

build: prepare
	go build -o bin/kvb-api

run:
	go run main.go

start-jaeger:
	docker-compose up -d

run-with-tracing: start-jaeger
	ENABLE_TRACING=true JAEGER_ENDPOINT="http://localhost:14268/api/traces" go run main.go
