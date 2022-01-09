package ports

import "context"

type StationMapperAdapter interface {
	GetStationIDForName(ctx context.Context, name string) (int, error)
}
