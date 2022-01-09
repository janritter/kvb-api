package ports

import (
	"context"

	"github.com/janritter/kvb-api/domains"
)

type DepartureService interface {
	GetDeparturesForMatchingStation(ctx context.Context, station string) (domains.Departures, error)
}
