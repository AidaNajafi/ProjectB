package service

import (
	"authentication/internal/provider"
	"fmt"
)


type HotelService struct {
	Provider *provider.Provider
}

func NewHotelService(prv *provider.Provider) *HotelService {
	return &HotelService{Provider: prv}
}

type HotelByDate struct {
	Count  int
	Date   string
	Hotels []HotelDetail
}

func (h *HotelService) GetHotelsByDate(date string) (HotelByDate, error) {
	hotels, err := h.Provider.HotelByDate(date)
	if err != nil {
		return HotelByDate{}, fmt.Errorf("Failed to Get hotels for this date %s! %w", date, err)
	}
	var hotelList []HotelDetail
	for _, providerHotel := range hotels.Hotels {
		mappedHotelDetail := MappHotelDetails(providerHotel)
		hotelList = append(hotelList, mappedHotelDetail)
	}
	return HotelByDate{
		Count:  hotels.Count,
		Date:   hotels.Date,
		Hotels: hotelList}, nil
}

func MappHotelDetails(ProviderHotel provider.HotelDetail) HotelDetail {
	return HotelDetail{
		ID:             ProviderHotel.ID,
		Name:           ProviderHotel.Name,
		Status:         ProviderHotel.Status,
		Description:    ProviderHotel.Description,
		ImageUrl:       ProviderHotel.ImageUrl,
		Address:        ProviderHotel.Address,
		City:           ProviderHotel.City,
		Amenities:      ProviderHotel.Amenities,
		Rating:         ProviderHotel.Rating,
		TotalRooms:     ProviderHotel.TotalRooms,
		AvailableRooms: ProviderHotel.AvailableRooms,
	}
}

type HotelDetail struct {
	ID             int
	Name           string
	Status         string
	Description    string
	ImageUrl       string
	Address        string
	City           string
	Amenities      []string
	Rating         float64
	TotalRooms     int
	AvailableRooms int
}

func (h *HotelService) GetHotelByID(id int) (HotelDetail, error) {

	hotelDetail, err := h.Provider.HotelByID(id)
	if err != nil {
		return HotelDetail{}, fmt.Errorf("Failed to find hotel by this id %d : %w", id, err)
	}
	mappedHotelId := MappHotelDetails(hotelDetail)
	return mappedHotelId, nil
}

type HotelRoom struct {
	Count   int
	HotelID int
	Rooms   []RoomDetail
}

type RoomDetail struct {
	ID          int
	HotelID     int
	Name        string
	Price       int
	Currency    string
	Capacity    int
	BedType     string
	SizeSQM     int
	IsAvailable bool
}

func (h *HotelService) GetHotelRooms(id int, date_from, date_to string) (HotelRoom, error) {
	RoomsList, err := h.Provider.HotelRooms(id, date_from, date_to)
	if err != nil {
		return HotelRoom{}, fmt.Errorf("Failed to Get rooms with these params")
	}
	var roomDetail []RoomDetail
	for _, rooms := range RoomsList.Rooms {
		mappedRoomDetail := MapRoomDetails(rooms)
		roomDetail = append(roomDetail, mappedRoomDetail)
	}
	return HotelRoom{
		Count:   RoomsList.Count,
		HotelID: RoomsList.HotelID,
		Rooms:   roomDetail,
	}, nil

}

func MapRoomDetails(provider provider.RoomDetail) RoomDetail {
	return RoomDetail{
		ID:          provider.ID,
		HotelID:     provider.HotelID,
		Name:        provider.Name,
		Price:       provider.Price,
		Currency:    provider.Currency,
		Capacity:    provider.Capacity,
		BedType:     provider.BedType,
		SizeSQM:     provider.SizeSQM,
		IsAvailable: provider.IsAvailable,
	}
}

type ReservationResponse struct {
	ID        int
	RoomID    int
	HotelID   int
	UserID    int
	UserPhone string
	DateFrom  string
	DateTo    string
	Status    string
}

type ReservationRequest struct {
	RoomID    int
	DateFrom  string
	DateTo    string
	UserID    int
	UserPhone string
}

func (h *HotelService) ReserveHotelService(req ReservationRequest) (ReservationResponse, error) {
	if req.RoomID <= 0 || req.UserID <= 0 || req.DateFrom == "" || req.DateTo == "" {
		return ReservationResponse{}, fmt.Errorf("invalid reservation request")
	}
	request := provider.ReservationRequest{
		RoomID:    req.RoomID,
		DateFrom:  req.DateFrom,
		DateTo:    req.DateTo,
		UserID:    req.UserID,
		UserPhone: req.UserPhone,
	}
	resp, err := h.Provider.ReserveHotel(request)
	if err != nil {
		return ReservationResponse{}, fmt.Errorf("Failed to reserve hotel: %w", err)
	}
	result := ReservationResponse{
		ID:        resp.ID,
		RoomID:    resp.RoomID,
		HotelID:   resp.HotelID,
		UserID:    resp.UserID,
		UserPhone: resp.UserPhone,
		DateFrom:  resp.DateFrom,
		DateTo:    resp.DateTo,
		Status:    resp.Status,
	}
	return result, nil
}
