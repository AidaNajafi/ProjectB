package repository

import (
	"authentication/internal/provider"
	"context"
)

type UserInfo struct {
	ID        int
	Name      string
	Username  string
	Email     string
	UserPhone string
	Password  string
}

type ReservationResponse struct {
	ID         int64
	ProviderID int64
	RoomID     int64
	HotelID    int64
	UserID     int64
	UserPhone  string
	DateFrom   string
	DateTo     string
	Status     string
}

type ReservationRequest struct {
	RoomID    int64
	UserID    int64
	UserPhone string
	DateFrom  string
	DateTo    string
}

type Store interface {
	CreateUser(user UserInfo) error
	GetUserByUsername(username string) (*UserInfo, error)
	CreatePendingReservation(ctx context.Context, p ReservationRequest) (int64, error)
	ConfirmedReservation(ctx context.Context, id int64, pr ReservationResponse) error
	FailedReservation(ctx context.Context, id int64, pr ReservationResponse, reason provider.FailureReason) error
}


