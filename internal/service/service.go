package service

import (
	"context"
	repository "grpc-project/internal/repo"
	"grpc-project/proto"
	"math/rand"
	"sync"
)

// type FlightSolutionQuery struct {
// 	Origin        string
// 	Dest          string
// 	DepartureDate time.Time
// }

// type FlightSolutions struct {
// 	AirlineCode string
// 	Price       int64
// 	FareClass   string
// 	Aircraft    string
// }

type ServiceFlightSolution struct {
	providers []repository.FlightProvider
}

func New(p []repository.FlightProvider) *ServiceFlightSolution {
	return &ServiceFlightSolution{providers: p}
}

func (s *ServiceFlightSolution) GetAggregateFlights(ctx context.Context, filter *proto.FlightSolutionQuery) ([]*proto.FlightSolutions, error) {
	results := make(chan []*proto.FlightSolutions)
	var wg sync.WaitGroup
	taskCount := 10000
	for i := 0; i < taskCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			provider := s.providers[rand.Intn(len(s.providers))]
			flights, err := provider.Search(ctx, filter.DepartureDate, filter.Origin, filter.Dest)
			if err != nil {
				results <- nil
				return
			}
			mappedFlights := make([]*proto.FlightSolutions, len(flights))
			for i, flight := range flights {
				mappedFlights[i] = *proto.FlightSolutions{
					AirlineCode: flight.AirlineCode,
					Price:       flight.Price,
					FareClass:   flight.FareClass,
					Aircraft:    flight.Aircraft,
				}
			}
			results <- mappedFlights
		}(i)

	}
	go func() {
		wg.Wait()
		close(results)
	}()

	aggregatedFlightSolutions := make([]*proto.FlightSolutions, 0)
	for flights := range results {
		if flights != nil {
			aggregatedFlightSolutions = append(aggregatedFlightSolutions, flights...)
		}
	}

	return aggregatedFlightSolutions, nil

}
