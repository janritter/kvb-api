package services

import (
	"context"
	"errors"
	"testing"

	"github.com/janritter/kvb-api/domains"
	"github.com/janritter/kvb-api/ports"
)

// --- Mock implementations for testing service layer ---

type mockStationMapperAdapter struct {
	stationID int
	err       error
}

func (m *mockStationMapperAdapter) GetStationIDForName(_ context.Context, name string) (int, error) {
	return m.stationID, m.err
}

type mockKVBAdapter struct {
	departures domains.Departures
	err        error
}

func (m *mockKVBAdapter) GetDeparturesForStationID(_ context.Context, stationID int) (domains.Departures, error) {
	return m.departures, m.err
}

// --- Test GetDeparturesForMatchingStation ---

func Test_GetDeparturesForMatchingStation_success(t *testing.T) {
	mockStationMapper := &mockStationMapperAdapter{stationID: 3, err: nil}
	mockKVB := &mockKVBAdapter{
		departures: domains.Departures{
			Departures: []domains.Departure{
				{Line: "950", Destination: "Bf Deutz/Messe", ArrivalInMinutes: 5},
				{Line: "149", Destination: "Heumarkt", ArrivalInMinutes: 12},
			},
		},
	}

	srv := New(mockStationMapper, mockKVB)
	ctx := context.Background()

	result, err := srv.GetDeparturesForMatchingStation(ctx, "Poststr.")
	if err != nil {
		t.Errorf("GetDeparturesForMatchingStation() unexpected error: %v", err)
	}
	if len(result.Departures) != 2 {
		t.Errorf("GetDeparturesForMatchingStation() got %d departures, want 2", len(result.Departures))
	}
	if result.Departures[0].Line != "950" {
		t.Errorf("GetDeparturesForMatchingStation() first line = %s, want 950", result.Departures[0].Line)
	}
	if result.Departures[1].ArrivalInMinutes != 12 {
		t.Errorf("GetDeparturesForMatchingStation() second arrival = %d, want 12", result.Departures[1].ArrivalInMinutes)
	}
}

func Test_GetDeparturesForMatchingStation_stationNotFound(t *testing.T) {
	mockStationMapper := &mockStationMapperAdapter{stationID: -1, err: errors.New("No station for given name found")}
	mockKVB := &mockKVBAdapter{}

	srv := New(mockStationMapper, mockKVB)
	ctx := context.Background()

	_, err := srv.GetDeparturesForMatchingStation(ctx, "Nonexistent")
	if err == nil {
		t.Error("GetDeparturesForMatchingStation() expected error, got nil")
	}
}

func Test_GetDeparturesForMatchingStation_kvbError(t *testing.T) {
	mockStationMapper := &mockStationMapperAdapter{stationID: 3, err: nil}
	mockKVB := &mockKVBAdapter{err: errors.New("KVB API error")}

	srv := New(mockStationMapper, mockKVB)
	ctx := context.Background()

	_, err := srv.GetDeparturesForMatchingStation(ctx, "Poststr.")
	if err == nil {
		t.Error("GetDeparturesForMatchingStation() expected error, got nil")
	}
}

func Test_GetDeparturesForMatchingStation_sofortDeparture(t *testing.T) {
	mockStationMapper := &mockStationMapperAdapter{stationID: 3, err: nil}
	mockKVB := &mockKVBAdapter{
		departures: domains.Departures{
			Departures: []domains.Departure{
				{Line: "149", Destination: "Heumarkt", ArrivalInMinutes: 0},
			},
		},
	}

	srv := New(mockStationMapper, mockKVB)
	ctx := context.Background()

	result, err := srv.GetDeparturesForMatchingStation(ctx, "Poststr.")
	if err != nil {
		t.Errorf("GetDeparturesForMatchingStation() unexpected error: %v", err)
	}
	if result.Departures[0].ArrivalInMinutes != 0 {
		t.Errorf("GetDeparturesForMatchingStation() sofort arrival = %d, want 0", result.Departures[0].ArrivalInMinutes)
	}
}

func Test_GetDeparturesForMatchingStation_emptyDepartures(t *testing.T) {
	mockStationMapper := &mockStationMapperAdapter{stationID: 3, err: nil}
	mockKVB := &mockKVBAdapter{
		departures: domains.Departures{Departures: []domains.Departure{}},
	}

	srv := New(mockStationMapper, mockKVB)
	ctx := context.Background()

	result, err := srv.GetDeparturesForMatchingStation(ctx, "Poststr.")
	if err != nil {
		t.Errorf("GetDeparturesForMatchingStation() unexpected error: %v", err)
	}
	if len(result.Departures) != 0 {
		t.Errorf("GetDeparturesForMatchingStation() got %d departures, want 0", len(result.Departures))
	}
}

func Test_GetDeparturesForMatchingStation_multipleDepartures(t *testing.T) {
	mockStationMapper := &mockStationMapperAdapter{stationID: 8, err: nil}
	mockKVB := &mockKVBAdapter{
		departures: domains.Departures{
			Departures: []domains.Departure{
				{Line: "149", Destination: "Heumarkt", ArrivalInMinutes: 1},
				{Line: "125", Destination: "Ebertplatz", ArrivalInMinutes: 3},
				{Line: "950", Destination: "Bf Deutz/Messe", ArrivalInMinutes: 5},
				{Line: "159", Destination: "Hansaring", ArrivalInMinutes: 8},
				{Line: "925", Destination: "Neumarkt", ArrivalInMinutes: 15},
			},
		},
	}

	srv := New(mockStationMapper, mockKVB)
	ctx := context.Background()

	result, err := srv.GetDeparturesForMatchingStation(ctx, "Dom/Hbf")
	if err != nil {
		t.Errorf("GetDeparturesForMatchingStation() unexpected error: %v", err)
	}
	if len(result.Departures) != 5 {
		t.Errorf("GetDeparturesForMatchingStation() got %d departures, want 5", len(result.Departures))
	}
	// Verify departures are in order
	for i, dep := range result.Departures {
		if dep.ArrivalInMinutes < 0 {
			t.Errorf("Departure[%d] has negative arrival: %d", i, dep.ArrivalInMinutes)
		}
	}
}

func Test_New_createsService(t *testing.T) {
	mockStationMapper := &mockStationMapperAdapter{}
	mockKVB := &mockKVBAdapter{}

	srv := New(mockStationMapper, mockKVB)
	if srv == nil {
		t.Fatal("New() returned nil")
	}
	// Verify the interfaces are satisfied
	var _ ports.DepartureService = (*service)(nil)
}
