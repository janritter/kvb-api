package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/janritter/kvb-api/adapters"
	"github.com/janritter/kvb-api/services"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	service = "kvb-api"
)

func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
		)),
	)
	return tp, nil
}

func main() {
	tp, err := tracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}
	otel.SetTracerProvider(tp)

	kvbAdapter := adapters.NewKVBAdapter()
	stationMapperAdapter := adapters.NewStationMapperAdapter()
	departureService := services.New(stationMapperAdapter, kvbAdapter)

	r := mux.NewRouter()
	r.Use(otelmux.Middleware("kvb-api-webserver"))

	r.HandleFunc("/v1/departures/stations/{key}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		searchStation := vars["key"]

		departures, _ := departureService.GetDeparturesForMatchingStation(r.Context(), searchStation)
		payload, err := json.Marshal(departures)
		if err != nil {
			log.Printf("Error marshalling departures: %s", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)
	}))

	srv := &http.Server{
		Handler: r,
		Addr:    ":8080",

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Running webserver on port 8080")

	log.Fatal(srv.ListenAndServe())
}
