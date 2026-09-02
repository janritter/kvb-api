package domains

import (
	"encoding/json"
	"testing"
)

func Test_Departures_MarshalJSON(t *testing.T) {
	departures := Departures{
		Departures: []Departure{
			{Line: "149", Destination: "Heumarkt", ArrivalInMinutes: 5},
			{Line: "950", Destination: "Bf Deutz/Messe", ArrivalInMinutes: 12},
		},
	}

	data, err := json.Marshal(departures)
	if err != nil {
		t.Fatalf("Departures.MarshalJSON() unexpected error: %v", err)
	}

	// Verify JSON contains expected fields
	jsonStr := string(data)
	expectedFields := []string{`"departures"`, `"line"`, `"destination"`, `"arrivalInMinutes"`}
	for _, field := range expectedFields {
		if !contains(jsonStr, field) {
			t.Errorf("JSON output missing field %q: %s", field, jsonStr)
		}
	}
}

func Test_Departures_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"departures": [
			{"line": "149", "destination": "Heumarkt", "arrivalInMinutes": 5},
			{"line": "950", "destination": "Bf Deutz/Messe", "arrivalInMinutes": 12}
		]
	}`

	var departures Departures
	err := json.Unmarshal([]byte(jsonData), &departures)
	if err != nil {
		t.Fatalf("Departures.UnmarshalJSON() unexpected error: %v", err)
	}

	if len(departures.Departures) != 2 {
		t.Errorf("Got %d departures, want 2", len(departures.Departures))
	}

	if departures.Departures[0].Line != "149" {
		t.Errorf("First line = %q, want %q", departures.Departures[0].Line, "149")
	}
	if departures.Departures[1].Destination != "Bf Deutz/Messe" {
		t.Errorf("Second destination = %q, want %q", departures.Departures[1].Destination, "Bf Deutz/Messe")
	}
	if departures.Departures[0].ArrivalInMinutes != 5 {
		t.Errorf("First arrival = %d, want 5", departures.Departures[0].ArrivalInMinutes)
	}
}

func Test_Departures_EmptyJSON(t *testing.T) {
	jsonData := `{"departures": []}`

	var departures Departures
	err := json.Unmarshal([]byte(jsonData), &departures)
	if err != nil {
		t.Fatalf("Departures.UnmarshalJSON() unexpected error: %v", err)
	}

	if len(departures.Departures) != 0 {
		t.Errorf("Got %d departures, want 0", len(departures.Departures))
	}
}

func Test_Departure_JsonTags(t *testing.T) {
	tests := []struct {
		name       string
		structName string
		jsonTag    string
	}{
		{"Departures.Departures", "Departures", "departures"},
		{"Departure.Line", "Departure", "line"},
		{"Departure.Destination", "Departure", "destination"},
		{"Departure.ArrivalInMinutes", "Departure", "arrivalInMinutes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify struct exists and has the expected JSON tag
			// by marshaling a zero value and checking the output
			switch tt.structName {
			case "Departures":
				d := Departures{}
				data, _ := json.Marshal(d)
				if !contains(string(data), `"departures"`) {
					t.Errorf("Departures missing json tag %q", tt.jsonTag)
				}
			case "Departure":
				d := Departure{}
				data, _ := json.Marshal(d)
				if !contains(string(data), `"line"`) {
					t.Errorf("Departure.Line missing json tag %q", "line")
				}
				if !contains(string(data), `"destination"`) {
					t.Errorf("Departure.Destination missing json tag %q", "destination")
				}
				if !contains(string(data), `"arrivalInMinutes"`) {
					t.Errorf("Departure.ArrivalInMinutes missing json tag %q", "arrivalInMinutes")
				}
			}
		})
	}
}

func Test_Departure_NegativeArrival(t *testing.T) {
	// Test that negative arrival values (for past departures) are handled
	dep := Departure{
		Line:             "149",
		Destination:      "Heumarkt",
		ArrivalInMinutes: -1,
	}

	data, err := json.Marshal(dep)
	if err != nil {
		t.Fatalf("Departure.MarshalJSON() unexpected error: %v", err)
	}

	var result Departure
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Departure.UnmarshalJSON() unexpected error: %v", err)
	}
	if result.ArrivalInMinutes != -1 {
		t.Errorf("Negative arrival not preserved: got %d", result.ArrivalInMinutes)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
