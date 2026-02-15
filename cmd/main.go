package main

import (
	"log"
	"net/http"
	"project1/internal/api"
	"project1/internal/repository"
	"project1/internal/service"
)

func main() {
	NewProviders := repository.NewFlightProvider()
	providers := []repository.FlightProvider{NewProviders}
	svc := service.New(providers)
	svr := api.NewHandler(svc)
	log.Println("Starting server")
	err := http.ListenAndServe(":8080", svr.FlightHandler())
	if err != nil {
		log.Fatal("server failed")
	}
}
