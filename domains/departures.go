package domains

type Departures struct {
	Departures []Departure `json:"departures"`
}

type Departure struct {
	Line             string `json:"line"`
	Destination      string `json:"destination"`
	ArrivalInMinutes int    `json:"arrivalInMinutes"`
}
