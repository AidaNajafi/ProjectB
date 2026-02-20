package main

import (
	"grpc-project/grpc/invoicer"
	"grpc-project/internal/handler"
	repository "grpc-project/internal/repo"
	"grpc-project/internal/service"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	providersNum := 5
	providers := make([]repository.FlightProvider, providersNum)
	for i := 0; i < providersNum; i++ {
		providers[i] = repository.NewFlightProvider()
	}
	svc := service.New(providers)

	RestSvr := handler.NewRestHandler(svc)
	log.Println("Starting rest server")
	go func() {
		err := http.ListenAndServe(":8080", RestSvr.FlightHandler())
		if err != nil {
			log.Fatal("rest server failed")
		}
	}()

	grpcServer := grpc.NewServer()
	log.Println("has the grpc started its server?")
	grpcSvr := handler.NewGrpcHandler(svc)
	reflection.Register(grpcServer)
	invoicer.RegisterSearchFlightServer(grpcServer, grpcSvr)
	log.Println("is it here?")

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to listen to the server")
	}
	log.Println("Starting grpc server")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("grpc server failed")
	}

}
