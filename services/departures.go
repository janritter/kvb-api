package services

import (
	"context"
	"log"

	"github.com/janritter/kvb-api/domains"
	"github.com/janritter/kvb-api/ports"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type service struct {
	stationMapperAdapter ports.StationMapperAdapter
	kvbAdapter           ports.KVBAdapter
}

func New(stationMapperAdapter ports.StationMapperAdapter, kvbAdapter ports.KVBAdapter) *service {
	return &service{
		stationMapperAdapter: stationMapperAdapter,
		kvbAdapter:           kvbAdapter,
	}
}

func (srv *service) GetDeparturesForMatchingStation(ctx context.Context, station string) (domains.Departures, error) {
	var span trace.Span
	ctx, span = otel.Tracer("kvb-api").Start(ctx, "GetDeparturesForMatchingStation")
	defer span.End()

	span.SetAttributes(attribute.String("station", station))

	stationID, err := srv.stationMapperAdapter.GetStationIDForName(ctx, station)
	if err != nil {
		log.Printf("Error getting station ID for name: %s", err)
		return domains.Departures{}, err
	}

	departures, err := srv.kvbAdapter.GetDeparturesForStationID(ctx, stationID)
	if err != nil {
		log.Printf("Error getting departures for station ID: %s", err)
		return domains.Departures{}, err
	}

	return departures, nil
}
