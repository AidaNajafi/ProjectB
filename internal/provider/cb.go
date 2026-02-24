package provider

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/sony/gobreaker"
)

type BreakerConfig struct {
	Name       string
	MaxRequest uint32
	Interval   time.Duration
	Timeout    time.Duration
}

func NewProviderBreaker(cfg BreakerConfig, onStateChange func(name string, from, to gobreaker.State)) *gobreaker.CircuitBreaker {
	st := gobreaker.Settings{
		Name:        cfg.Name,
		MaxRequests: cfg.MaxRequest,
		Interval:    cfg.Interval,
		Timeout:     cfg.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			if counts.Requests < 20 {
				return false
			}
			failures := float64(counts.TotalFailures)
			requests := float64(counts.Requests)
			return (failures / requests) >= 0.50
		},
		IsSuccessful: func(err error) bool {
			if err == nil {
				return true
			}
			if errors.Is(err, FailAlreadyReserved) || errors.Is(err, FailBadRequest) || errors.Is(err, FailUnAuthorized) || errors.Is(err, FailUnknown) {
				return true
			}
			return false
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			slog.Info("name: ", name, "from: ", string(from), "to: ", string(to))
			if onStateChange != nil {
				onStateChange(name, from, to)
			}
		},
	}
	return gobreaker.NewCircuitBreaker(st)

}

type CBProvider struct {
	next    HotelProvider
	breaker *gobreaker.CircuitBreaker
	timeout time.Duration
}

func NewCBProvider(next HotelProvider, cb *gobreaker.CircuitBreaker) *CBProvider {
	return &CBProvider{
		next:    next,
		breaker: cb,
	}
}

func (c *CBProvider) HotelByDate(ctx context.Context, date string) (HotelByDate, error) {
	out, err := c.breaker.Execute(func() (interface{}, error) {
		return c.next.HotelByDate(ctx, date)
	})
	if err != nil {
		return HotelByDate{}, nil
	}
	return out.(HotelByDate), nil
}

func (c *CBProvider) HotelByID(ctx context.Context, id int64) (HotelDetail, error) {
	out, err := c.breaker.Execute(func() (interface{}, error) {
		return c.next.HotelByID(ctx, id)
	})
	if err != nil {
		return HotelDetail{}, nil
	}
	return out.(HotelDetail), nil
}

func (c *CBProvider) HotelRooms(ctx context.Context, id int64, date_from, date_to string) (HotelRoom, error) {
	out, err := c.breaker.Execute(func() (interface{}, error) {
		return c.next.HotelRooms(ctx, id, date_from, date_to)
	})
	if err != nil {
		return HotelRoom{}, nil
	}
	return out.(HotelRoom), nil
}

func (c *CBProvider) ReserveHotel(ctx context.Context, req ReservationRequest) (ReservationResponse, error) {
	out, err := c.breaker.Execute(func() (interface{}, error) {
		return c.next.ReserveHotel(ctx, req)
	})
	if err != nil {
		return ReservationResponse{}, nil
	}
	return out.(ReservationResponse), nil
}
