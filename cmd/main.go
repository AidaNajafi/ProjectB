package cmd

import (
	"grpc-project/internal/proto/invoicer"
	"grpc-project/internal/repo"
	"grpc-project/internal/service"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	providersNum := 5
	providers := make([]repository.FlightProvider, providersNum)
	for i := 0; i < providersNum; i++ {
		providers[i] = repository.NewFlightProvider()
	}
	svc := service.New(providers)
	grpcServer := grpc.NewServer()
	invoicer.RegisterFlightServiceServer(grpcServer, svc)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to listen to the server")
	}
	log.Println("Starting server")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("server failed")
	}

}
