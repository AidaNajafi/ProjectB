package service

import (
	"context"
	"fmt"
	"project1/internal/repository"
	"time"
)

type FlightSolutionQuery struct {
	Origin        string
	Dest          string
	DepartureDate time.Time
}

type FlightSolutions struct {
	AirlineCode string
	Price       int64
	FareClass   string
	Aircraft    string
}

type ServiceFlightSolution struct {
	providers []repository.FlightProvider
}

func New(p []repository.FlightProvider) *ServiceFlightSolution {
	return &ServiceFlightSolution{providers: p}
}

func (s *ServiceFlightSolution) GetAggregateFlights(ctx context.Context, filter FlightSolutionQuery) ([]FlightSolutions, error) {
	AggregatedFlightSolutions := make([]FlightSolutions, 0)
	for _, provider := range s.providers {
		flights, err := provider.Search(ctx, filter.DepartureDate, filter.Origin, filter.Dest)
		if err != nil {
			return nil, fmt.Errorf("Search failed!: %w", err)
		}
		for _, f := range flights {
			AggregatedFlightSolutions = append(AggregatedFlightSolutions, FlightSolutions{
				AirlineCode: f.AirlineCode,
				Price:       f.Price,
				FareClass:   f.FareClass,
				Aircraft:    f.Aircraft,
			})
		}
	}
	return AggregatedFlightSolutions, nil

}
