package provider

type HotelProvider interface {
	HotelByDate(date string) (HotelByDate, error)
	HotelByID(id int64) (HotelDetail, error)
	HotelRooms(id int64, date_from, date_to string) (HotelRoom, error)
	ReserveHotel(req ReservationRequest) (ReservationResponse, error)
}
