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

type NewFlightProviderFunc func(ctx context.Context, departureDate time.Time, originIATA string, destinationIATA string) ([]FlightSolution, error)

func (f NewFlightProviderFunc) Search(ctx context.Context, departureDate time.Time, originIATA string, destinationIATA string) ([]FlightSolution, error) {
	return f(ctx, departureDate, originIATA, destinationIATA)
}
