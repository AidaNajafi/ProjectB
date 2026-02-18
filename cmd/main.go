package main

import (
	"log"
	"net/http"
	"project1/internal/api"
	"project1/internal/repository"
	"project1/internal/service"
)

func main() {
	providersNum := 5
	providers := make([]repository.FlightProvider, providersNum)
	for i := 0; i < providersNum; i++ {
		providers[i] = repository.NewFlightProvider()
	}
	svc := service.New(providers)
	svr := api.NewHandler(svc)
	log.Println("Starting server")
	err := http.ListenAndServe(":8080", svr.FlightHandler())
	if err != nil {
		log.Fatal("server failed")
	}
}
