package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"project1/internal/service"
	"time"
)

type SHandler struct {
	srv *service.FSolution
}

func NewHandler(srv *service.FSolution) *SHandler {
	return &SHandler{srv: srv}
}

type InInfo struct {
	Origin        string    `json:"origin"`
	Dest          string    `json:"dest"`
	DepartureDate time.Time `json:"departuredate"`
}

type OutInfo struct {
	AirlineCode string `json:"airlinecode"`
	Price       int64  `json:"price"`
	FareClass   string `json:"fareclass"`
	Aircraft    string `json:"aircraft"`
}

func (s *SHandler) FlightHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/flight-lists", s.GetFlights)
	return mux
}

func (s *SHandler) GetFlights(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}
	ctx := r.Context()
	var body InInfo
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "failed to decode input",
		})
		return
	}
	if body.DepartureDate.IsZero() || body.Origin == "" || body.Dest == "" {
		WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "All fields are required!",
		})
		return
	}
	info := service.PassengerInfo{
		Origin:        body.Origin,
		Dest:          body.Dest,
		DepartureDate: body.DepartureDate,
	}
	FlightList, err := s.srv.GetFlightList(ctx, info)
	if err != nil {
		status := http.StatusBadGateway
		if errors.Is(err, context.DeadlineExceeded) {
			status = http.StatusGatewayTimeout
		}
		WriteJSON(w, status, err.Error())
		return
	}
	output := make([]OutInfo, 0, len(FlightList))
	for _, f := range FlightList {
		output = append(output, OutInfo{
			AirlineCode: f.AirlineCode,
			Price:       f.Price,
			FareClass:   f.FareClass,
			Aircraft:    f.Aircraft,
		})
	}
	WriteJSON(w, http.StatusOK, output)

}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)

}
