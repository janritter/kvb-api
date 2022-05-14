package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/janritter/kvb-api/adapters"
	"github.com/janritter/kvb-api/services"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	service = "kvb-api"
)

func tracerProvider() (*tracesdk.TracerProvider, error) {
	ctx := context.Background()

	traceClient := otlptracegrpc.NewClient()
	traceExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		return nil, err
	}

	bsp := tracesdk.NewBatchSpanProcessor(traceExp)
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
		)),
		tracesdk.WithSpanProcessor(bsp),
	)

	return tp, nil
}

func main() {
	if os.Getenv("ENABLE_TRACING") == "true" {
		log.Println("Configuring trace provider")
		tp, err := tracerProvider()
		if err != nil {
			log.Fatal(err)
		}
		otel.SetTracerProvider(tp)
	}

	otel.SetTextMapPropagator(propagation.TraceContext{})

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
