package handler

import (
	"context"
	"grpc-project/grpc/invoicer"
	"grpc-project/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandler struct {
	invoicer.UnimplementedSearchFlightServer
	srv *service.ServiceFlightSolution
}

func NewGrpcHandler(srv *service.ServiceFlightSolution) *GrpcHandler {
	return &GrpcHandler{srv: srv}
}

func (s *GrpcHandler) GetAggregatedFlights(ctx context.Context, req *invoicer.GetAggregatedFlightsRequest) (*invoicer.GetAggregatedFlightsResponse, error) {
	if req.DepartureDate == "" || req.Dest == "" || req.Origin == "" {
		return nil, status.Error(codes.InvalidArgument, "origin, destination and departure date are required")
	}
	depTime, err := ParseStringDepartureDateToRFC(req.DepartureDate)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid departure date :%v", err)
	}
	info := service.FlightSolutionQuery{
		Origin:        req.Origin,
		Dest:          req.Dest,
		DepartureDate: depTime,
	}
	flights, err := s.srv.GetAggregateFlights(ctx, info)
	if err != nil {
		return nil, status.Error(codes.Unavailable, "Failed to get aggregated results from service")
	}
	output := make([]*invoicer.FlightSolutions, 0, len(flights))
	for _, f := range flights {
		output = append(output, &invoicer.FlightSolutions{
			AirlineCode: f.AirlineCode,
			Price:       f.Price,
			FareClass:   f.FareClass,
			Aircraft:    f.Aircraft,
		})
	}
	return &invoicer.GetAggregatedFlightsResponse{Solutions: output}, nil

}
