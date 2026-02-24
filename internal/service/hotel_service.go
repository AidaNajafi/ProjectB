package service

import (
	"authentication/internal/provider"
	"authentication/internal/repository"
	"context"
	"fmt"
	"log"
)

type HotelService struct {
	Provider *provider.CBProvider
	Store    repository.Store
}

func NewHotelService(prv *provider.CBProvider, store repository.Store) *HotelService {
	return &HotelService{Provider: prv, Store: store}
}

type HotelByDate struct {
	Count  int64
	Date   string
	Hotels []HotelDetail
}

func (h *HotelService) GetHotelsByDate(ctx context.Context, date string) (HotelByDate, error) {
	hotels, err := h.Provider.HotelByDate(ctx, date)
	if err != nil {
		return HotelByDate{}, fmt.Errorf("Failed to Get hotels for this date %s! %w", date, err)
	}
	hotelList:= make([]HotelDetail,0, len(hotels.Hotels))
	for _, providerHotel := range hotels.Hotels {
		mappedHotelDetail := mapHotelDetails(providerHotel)
		hotelList = append(hotelList, mappedHotelDetail)
	}
	return HotelByDate{
		Count:  hotels.Count,
		Date:   hotels.Date,
		Hotels: hotelList}, nil
}

func mapHotelDetails(ProviderHotel provider.HotelDetail) HotelDetail {
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
	ID             int64
	Name           string
	Status         string
	Description    string
	ImageUrl       string
	Address        string
	City           string
	Amenities      []string
	Rating         float64
	TotalRooms     int64
	AvailableRooms int64
}

func (h *HotelService) GetHotelByID(ctx context.Context, id int64) (HotelDetail, error) {

	hotelDetail, err := h.Provider.HotelByID(ctx, id)
	if err != nil {
		return HotelDetail{}, fmt.Errorf("Failed to find hotel by this id %d : %w", id, err)
	}
	mappedHotelId := mapHotelDetails(hotelDetail)
	return mappedHotelId, nil
}

type HotelRoom struct {
	Count   int64
	HotelID int64
	Rooms   []RoomDetail
}

type RoomDetail struct {
	ID          int64
	HotelID     int64
	Name        string
	Price       int64
	Currency    string
	Capacity    int64
	BedType     string
	SizeSQM     int64
	IsAvailable bool
}

func (h *HotelService) GetHotelRooms(ctx context.Context, id int64, date_from, date_to string) (HotelRoom, error) {
	RoomsList, err := h.Provider.HotelRooms(ctx, id, date_from, date_to)
	if err != nil {
		return HotelRoom{}, fmt.Errorf("Failed to Get rooms with these params")
	}
	roomDetail:= make([]RoomDetail, 0, len(RoomsList.Rooms))
	for _, rooms := range RoomsList.Rooms {
		mappedRoomDetail := mapRoomDetails(rooms)
		roomDetail = append(roomDetail, mappedRoomDetail)
	}
	return HotelRoom{
		Count:   RoomsList.Count,
		HotelID: RoomsList.HotelID,
		Rooms:   roomDetail,
	}, nil

}

func mapRoomDetails(provider provider.RoomDetail) RoomDetail {
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
	DateFrom  string
	DateTo    string
	UserID    int64
	UserPhone string
}

func (h *HotelService) ReserveHotelService(ctx context.Context, req ReservationRequest) (ReservationResponse, error) {
	if req.RoomID <= 0 || req.UserID <= 0 || req.DateFrom == "" || req.DateTo == "" {
		return ReservationResponse{}, fmt.Errorf("invalid reservation request")
	}
	repoReq := mapServiceRequestToRepository(req)
	id, err := h.Store.CreatePendingReservation(ctx, repoReq)
	if err != nil {
		return ReservationResponse{}, fmt.Errorf("failed creating pending reservation %w", err)
	}

	ProviderRequest := mapServiceRequestToProvider(req)

	ProviderResponse, err := h.Provider.ReserveHotel(ctx, ProviderRequest)
	RepoResponse := mapProviderResponseToRepository(ProviderResponse)
	if err != nil {
		reason := provider.ClassifyProviderError(err)
		log.Println(reason)
		storeErr := h.Store.FailedReservation(ctx, id, RepoResponse, reason)
		log.Printf("DEBUG provider err=%v | storeErr=%T %#v | storeErr==nil? %v",
			err, storeErr, storeErr, storeErr == nil,
		)
		if storeErr != nil {
			return ReservationResponse{}, fmt.Errorf("Failed to mark failure for reservation: %w : %s", storeErr, reason)
		}
		return ReservationResponse{}, fmt.Errorf("Failed to reserve hotel: %w : %s", err, reason)
	}

	err = h.Store.ConfirmedReservation(ctx, id, RepoResponse)
	if err != nil {
		reason := provider.ClassifyProviderError(err)
		return ReservationResponse{}, fmt.Errorf("Failed to confirm reservation: %w : %s", err, reason)
	}
	response := mapProviderResponseToService(ProviderResponse)
	return response, nil
}

func mapServiceRequestToProvider(pr ReservationRequest) provider.ReservationRequest {
	return provider.ReservationRequest{
		RoomID:    pr.RoomID,
		DateFrom:  pr.DateFrom,
		DateTo:    pr.DateTo,
		UserID:    pr.UserID,
		UserPhone: pr.UserPhone,
	}
}

func mapProviderResponseToService(pr provider.ReservationResponse) ReservationResponse {
	return ReservationResponse{
		ProviderID: pr.ProviderID,
		RoomID:     pr.RoomID,
		HotelID:    pr.HotelID,
		UserID:     pr.UserID,
		UserPhone:  pr.UserPhone,
		DateFrom:   pr.DateFrom,
		DateTo:     pr.DateTo,
		Status:     pr.Status,
	}
}

func mapServiceRequestToRepository(re ReservationRequest) repository.ReservationRequest {
	return repository.ReservationRequest{
		RoomID:    re.RoomID,
		DateFrom:  re.DateFrom,
		DateTo:    re.DateTo,
		UserID:    re.UserID,
		UserPhone: re.UserPhone,
	}
}

func mapProviderResponseToRepository(re provider.ReservationResponse) repository.ReservationResponse {
	return repository.ReservationResponse{
		ProviderID: re.ProviderID,
		RoomID:     re.RoomID,
		HotelID:    re.HotelID,
		UserID:     re.UserID,
		UserPhone:  re.UserPhone,
		DateFrom:   re.DateFrom,
		DateTo:     re.DateTo,
		Status:     re.Status,
	}
}
