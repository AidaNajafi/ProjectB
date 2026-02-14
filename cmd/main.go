package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"project1/internal/api"
	"project1/internal/repository"
	"project1/internal/service"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	provider := repository.NewFlightProviderFunc(func(ctx context.Context, departureDate time.Time, originIATA, destinationIATA string) ([]repository.FlightSolution, error) {
		delay := time.Duration(200+rand.Intn(1001)) * time.Millisecond
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
		return []repository.FlightSolution{
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
	svc := service.New(provider)
	svr := api.NewHandler(svc)
	log.Println("Starting server")
	err := http.ListenAndServe(":8080", svr.FlightHandler())
	if err != nil {
		log.Fatal("server failed")
	}
}
