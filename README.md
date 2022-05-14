# KVB API

[![Maintainability](https://api.codeclimate.com/v1/badges/c962d75f12cb52361657/maintainability)](https://codeclimate.com/github/janritter/kvb-api/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/c962d75f12cb52361657/test_coverage)](https://codeclimate.com/github/janritter/kvb-api/test_coverage)

> Rest API to get upcoming departures per [KVB](https://www.kvb.koeln/) train station

> Implemented in Go with hexagonal architecture and tracing via OpenTelemetry and Jaeger


KVB API provides one endpoint allowing you to get the upcoming departures for a provided station name

`http://localhost:8080/v1/departures/stations/{station_name}`

KVB API tries to find the best matching station name for your request, so it doesn't need to be the exact name

## Example

**Request**

GET `http://localhost:8080/v1/departures/stations/bensberg`

**Response**

```json
{
  "departures": [
    {
      "line": "1",
      "destination": "Weiden West",
      "arrivalInMinutes": 17
    },
    {
      "line": "1",
      "destination": "Junkersdorf",
      "arrivalInMinutes": 44
    }
  ]
}
```

## Build

The binary will be stored at `dist/kvb-api`

```make
make build
```

## Build Docker Image

The binary will be stored at `dist/kvb-api`

```make
make build-docker
```


## Usage

### Without Tracing

Start the KVB API **without tracing** by running:

```bash
./dist/kvb-api
```

To run with Docker, adapt the `docker-compose.yaml` and run:

```bash
docker-compose up
```

### With OpenTelemetry

Start the KVB API **with OpenTelemetry tracing** by running:

```bash
make run-with-tracing
```

To run with Docker, adapt the `docker-compose.yaml` and run:

```bash
docker-compose up
```

Make sure an OpenTelemetry collector is running on the provided URL.


## Development

### Hexagonal Architecture

- Adapters are stored in `adapters`
- Ports are stored in `ports`
- Business logic data structures are stored in `domains`
- Functions offered by the business logic are stored in `services`
