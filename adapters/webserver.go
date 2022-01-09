package adapters

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/janritter/kvb-api/ports"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

type WebserverAdapter struct {
	departureService ports.DepartureService
}

func NewWebserverAdapter(departureService ports.DepartureService) *WebserverAdapter {
	return &WebserverAdapter{
		departureService: departureService,
	}
}

func (adapter *WebserverAdapter) RunWebserver() {
	r := mux.NewRouter()
	r.Use(otelmux.Middleware("kvb-api-webserver"))

	r.HandleFunc("/v1/departures/stations/{key}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		searchStation := vars["key"]

		departures, _ := adapter.departureService.GetDeparturesForMatchingStation(r.Context(), searchStation)
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
