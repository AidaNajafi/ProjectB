package repository

import (
	"context"
	"time"
)

type FlightSolution struct {
	AirlineCode string
	Price       int64
	FareClass   string
	Aircraft    string
}

type FlightProvider interface {
	Search(ctx context.Context, departureDate time.Time, originIATA string, destinationIATA string) ([]FlightSolution, error)
}

type FlightProviderFunc func(ctx context.Context, departureDate time.Time, originIATA string, destinationIATA string) ([]FlightSolution, error)

func (f FlightProviderFunc) Search(ctx context.Context, departureDate time.Time, originIATA string, destinationIATA string) ([]FlightSolution, error) {
	return f(ctx, departureDate, originIATA, destinationIATA)
}


func NewFlightProvider() FlightProvider{
	return FlightProviderFunc (func(ctx context.Context, departureDate time.Time, originIATA string, destinationIATA string) ([]FlightSolution, error){
		time.Sleep(5 *time.Millisecond)
		return []FlightSolution{
			{AirlineCode: "IR", Price: 210_00, FareClass: "Y", Aircraft: "A320"},
			{AirlineCode: "W5", Price: 195_00, FareClass: "M", Aircraft: "MD83"},
			{AirlineCode: "EP", Price: 175_00, FareClass: "Q", Aircraft: "B737"},
			{AirlineCode: "QB", Price: 185_00, FareClass: "H", Aircraft: "A319"},
			{AirlineCode: "NV", Price: 165_00, FareClass: "K", Aircraft: "A320"},
			{AirlineCode: "B9", Price: 205_00, FareClass: "L", Aircraft: "A320"},
			{AirlineCode: "Y9", Price: 225_00, FareClass: "V", Aircraft: "A321"},
			{AirlineCode: "IV", Price: 190_00, FareClass: "T", Aircraft: "B737"},
			{AirlineCode: "JI", Price: 155_00, FareClass: "S", Aircraft: "MD82"},
			{AirlineCode: "ZV", Price: 215_00, FareClass: "N", Aircraft: "A320"},
			{AirlineCode: "HH", Price: 160_00, FareClass: "E", Aircraft: "B737"},
			{AirlineCode: "I3", Price: 150_00, FareClass: "B", Aircraft: "ATR72"},
			{AirlineCode: "VR", Price: 158_00, FareClass: "R", Aircraft: "MD83"},
			{AirlineCode: "SR", Price: 170_00, FareClass: "G", Aircraft: "A320"},
			{AirlineCode: "TB", Price: 240_00, FareClass: "U", Aircraft: "A320"},
		}, nil
	})
}