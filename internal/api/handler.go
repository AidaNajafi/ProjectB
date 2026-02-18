package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"project1/internal/service"
	"time"
)

type Handler struct {
	srv *service.ServiceFlightSolution
}

func NewHandler(srv *service.ServiceFlightSolution) *Handler {
	return &Handler{srv: srv}
}

type HandlerFlightSolutionQuery struct {
	Origin        string `json:"origin"`
	Dest          string `json:"dest"`
	DepartureDate string `json:"departuredate"`
}

type HandlerFlightSolution struct {
	AirlineCode string `json:"airlinecode"`
	Price       int64  `json:"price"`
	FareClass   string `json:"fareclass"`
	Aircraft    string `json:"aircraft"`
}

func (s *Handler) FlightHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/flight-lists", s.GetFlights)
	return mux
}

func ParseDate(dateStr string) (time.Time, error) {
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

func (s *Handler) GetFlights(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodPost {
		WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}
	log.Println("Received request for flights")
	ctx := r.Context()
	var body HandlerFlightSolutionQuery
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("Failed to decode input: %v\n", err)
		WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "failed to decode input",
		})
		return
	}

	if body.DepartureDate == "" || body.Origin == "" || body.Dest == "" {
		WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "All fields are required!",
		})
		return
	}
	date, err := ParseDate(body.DepartureDate)
	if err != nil {
		WriteJSON(w, 500, map[string]string{
			"error": err.Error(),
			"msg":   "Failed to parse date",
		})
		return
	}
	info := service.FlightSolutionQuery{
		Origin:        body.Origin,
		Dest:          body.Dest,
		DepartureDate: date,
	}
	FlightList, err := s.srv.GetAggregateFlights(ctx, info)
	if err != nil {
		status := http.StatusBadGateway
		if errors.Is(err, context.DeadlineExceeded) {
			status = http.StatusGatewayTimeout
		}
		WriteJSON(w, status, err.Error())
		return
	}
	output := make([]HandlerFlightSolution, 0, len(FlightList))
	for _, f := range FlightList {
		output = append(output, HandlerFlightSolution{
			AirlineCode: f.AirlineCode,
			Price:       f.Price,
			FareClass:   f.FareClass,
			Aircraft:    f.Aircraft,
		})
	}
	WriteJSON(w, http.StatusOK, output)
	latency := time.Since(start)
	log.Println("latency: ", latency)

}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)

}
