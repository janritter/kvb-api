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

The binary binary will be stored at `bin/kvb-api`

```make
make build
```

## Usage

### Without Tracing

Start the KVB API **without tracing** by running:

```bash
./bin/kvb-api
```

### With Jaeger Tracing

Start the KVB API **with Jaeger tracing** by running:

```bash
ENABLE_TRACING=true JAEGER_ENDPOINT="http://localhost:14268/api/traces" ./bin/kvb-api 
```

Make sure Jaeger is running on the provided URL.
To start the included Jager all in one container run the following command:

```bash
docker-compose up -d
```

## Development

### Hexagonal Architecture

- Adapters are stored in `adapters`
- Ports are stored in `ports`
- Business logic data structures are stored in `domains`
- Functions offered by the business logic are stored in `services`
