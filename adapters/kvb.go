package adapters

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/janritter/kvb-api/domains"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/text/encoding/charmap"
)

type KVBAdapter struct{}

func NewKVBAdapter() *KVBAdapter {
	return &KVBAdapter{}
}

func (adapter *KVBAdapter) GetDeparturesForStationID(ctx context.Context, stationID int) (domains.Departures, error) {
	var span trace.Span
	ctx, span = otel.Tracer("kvb-api").Start(ctx, "GetDeparturesForStationID")
	defer span.End()

	span.SetAttributes(attribute.Int("stationID", stationID))

	url := fmt.Sprintf("https://www.kvb.koeln/generated/?aktion=show&code=%d", stationID)

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return domains.Departures{}, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return domains.Departures{}, err
	}

	departures := []domains.Departure{}

	_, span = otel.Tracer("kvb-api").Start(ctx, "GetDeparturesForStationID.Parse")
	defer span.End()
	doc.Find("body > div > table:nth-child(2) > tbody > tr").Each(func(i int, s *goquery.Selection) {
		if i != 0 {
			route := s.Find("td:nth-child(1)").Text()
			destination := s.Find("td:nth-child(2)").Text()
			arrivalTimeString := s.Find("td:nth-child(3)").Text()

			// Build correct time from Sofort and 2 Min
			arrivalTime := -1
			if strings.TrimSpace(arrivalTimeString) == "Sofort" {
				arrivalTime = 0
			} else {
				arrivalTimeString = strings.Replace(arrivalTimeString, "Min", "", -1)
				arrivalTime, err = strconv.Atoi(strings.TrimSpace(arrivalTimeString))
				if err != nil {
					log.Println(err)
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
				}
			}

			// Response is ISO-8859-1, transfer to utf-8
			destination, err = charmap.ISO8859_1.NewDecoder().String(destination)
			if err != nil {
				log.Println(err)
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}

			singleDeparture := domains.Departure{
				Line:             strings.TrimSpace(route),
				Destination:      destination,
				ArrivalInMinutes: arrivalTime,
			}
			departures = append(departures, singleDeparture)
		}
	})

	return domains.Departures{
		Departures: departures,
	}, nil
}
