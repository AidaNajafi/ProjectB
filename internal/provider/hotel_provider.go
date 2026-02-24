package provider

import "context"

type HotelProvider interface {
	HotelByDate(ctx context.Context, date string) (HotelByDate, error)
	HotelByID(ctx context.Context, id int64) (HotelDetail, error)
	HotelRooms(ctx context.Context, id int64, date_from, date_to string) (HotelRoom, error)
	ReserveHotel(ctx context.Context, req ReservationRequest) (ReservationResponse, error)
}
