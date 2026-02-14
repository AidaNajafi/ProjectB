package service

import (
	"context"
	"fmt"
	"project1/internal/repository"
	"time"
)

type PassengerInfo struct {
	Origin        string
	Dest          string
	DepartureDate time.Time
}

type FlightInfo struct {
	AirlineCode string
	Price       int64
	FareClass   string
	Aircraft    string
}

type FSolution struct {
	flp repository.FlightProvider
}

func New(flp repository.FlightProvider) *FSolution {
	return &FSolution{flp: flp}
}

func (f *FSolution) GetFlightList(ctx context.Context, info PassengerInfo) ([]FlightInfo, error) {
	solutions, err := f.flp.Search(ctx, info.DepartureDate, info.Origin, info.Dest)
	if err != nil {
		return nil, fmt.Errorf("Search failed: %w", err)
	}
	list := make([]FlightInfo, 0, len(solutions))
	for _, s := range solutions {
		list = append(list, FlightInfo{
			AirlineCode: s.AirlineCode,
			Price:       s.Price,
			FareClass:   s.FareClass,
			Aircraft:    s.Aircraft,
		})
	}
	return list, nil
}
