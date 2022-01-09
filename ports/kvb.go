package ports

import (
	"context"

	"github.com/janritter/kvb-api/domains"
)

type KVBAdapter interface {
	GetDeparturesForStationID(ctx context.Context, stationID int) (domains.Departures, error)
}
