package provider

import (
	"authentication/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type ProviderMethods interface {
	HotelByDate(date string) (HotelByDate, error)
	HotelByID(id int) (HotelDetail, error)
	HotelRooms(id int, date_from, date_to string) (HotelRoom, error)
	ReserveHotel(req ReservationRequest) (ReservationResponse, error)
}

type Provider struct {
	config *config.Config
	Http   *http.Client
}

func NewProvider(cfg *config.Config) *Provider {
	return &Provider{config: cfg}
}

type HotelDetail struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	Status         string   `json:"status"`
	Description    string   `json:"description"`
	ImageUrl       string   `json:"image_url"`
	Address        string   `json:"address"`
	City           string   `json:"city"`
	Amenities      []string `json:"amenities"`
	Rating         float64  `json:"rating"`
	TotalRooms     int      `json:"totalrooms"`
	AvailableRooms int      `json:"availablerooms"`
}

type HotelRoom struct {
	Count   int `json:"count"`
	HotelID int `json:"hotelid"`
	Rooms   []RoomDetail
}

type RoomDetail struct {
	ID          int    `json:"id"`
	HotelID     int    `json:"hotel"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Currency    string `json:"currency"`
	Capacity    int    `json:"capacity"`
	BedType     string `json:"bedtype"`
	SizeSQM     int    `json:"sizesqm"`
	IsAvailable bool   `json:"isavailable"`
}

type HotelByDate struct {
	Count  int    `json:"count"`
	Date   string `json:"date"`
	Hotels []HotelDetail
}

func (p *Provider) HotelByDate(date string) (HotelByDate, error) {
	url := fmt.Sprintf("%s/hotels?date=%s", p.config.BaseURL, date)

	body, err := StartProviderGetServer(url, p.config.ApiKey)
	if err != nil {
		return HotelByDate{}, fmt.Errorf("Failed to connect to provider")
	}
	var hotel HotelByDate
	if err := json.Unmarshal(body, &hotel); err != nil {
		return HotelByDate{}, fmt.Errorf("Failed to unmarshal request body to json")
	}
	return hotel, nil
}

func (p *Provider) HotelByID(id int) (HotelDetail, error) {

	url := fmt.Sprintf("%s/hotels/%d", p.config.BaseURL, id)
	body, err := StartProviderGetServer(url, p.config.ApiKey)
	if err != nil {
		return HotelDetail{}, fmt.Errorf("Failed to connect to provider %w", err)
	}

	var hotelDetail HotelDetail
	if err := json.Unmarshal(body, &hotelDetail); err != nil {
		return HotelDetail{}, fmt.Errorf("Failed to unmarshal request body to json %w", err)
	}
	return hotelDetail, nil

}

func (p *Provider) HotelRooms(id int, date_from, date_to string) (HotelRoom, error) {
	url := fmt.Sprintf("%s/hotels/%d/rooms?date_from=%s&date_to=%s", p.config.BaseURL, id, date_from, date_to)
	body, err := StartProviderGetServer(url, p.config.ApiKey)
	if err != nil {
		return HotelRoom{}, fmt.Errorf("Failed to connect to provider")
	}
	var hotelRooms HotelRoom
	if err := json.Unmarshal(body, &hotelRooms); err != nil {
		return HotelRoom{}, fmt.Errorf("Failed to unmarshal request body to json")
	}
	return hotelRooms, nil
}

func StartProviderGetServer(url, apiKey string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %w", err)
	}
	req.Header.Set("X-API-Key", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to provider:  %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read request %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Provider server not up : %d", resp.StatusCode)
	}

	return body, nil
}

type ReservationResponse struct {
	ID        int    `json:"id"`
	RoomID    int    `json:"room_id"`
	HotelID   int    `json:"hotel_id"`
	UserID    int    `json:"user_id"`
	UserPhone string `json:"user_phone"`
	DateFrom  string `json:"date_from"`
	DateTo    string `json:"date_to"`
	Status    string `json:"status"`
}

type ReservationRequest struct {
	RoomID    int    `json:"room_id"`
	DateFrom  string `json:"date_from"`
	DateTo    string `json:"date_to"`
	UserID    int    `json:"user_id"`
	UserPhone string `json:"user_phone"`
}

func (p *Provider) ReserveHotel(req ReservationRequest) (ReservationResponse, error) {
	base := strings.TrimRight(p.config.BaseURL, "/")
	url := fmt.Sprintf("%s/reservations", base)
	log.Println(url)
	body, err := StartProviderPostServer(url, p.config.ApiKey, req)
	if err != nil {
		return ReservationResponse{}, fmt.Errorf("Failed to start provider post server: %w", err)
	}
	var response ReservationResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return ReservationResponse{}, fmt.Errorf("Failed to unmarshal request body to response", err)
	}
	return response, nil

}

func StartProviderPostServer(url, apikey string, payload any) ([]byte, error) {
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal the request body(payload) %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("Failed to create request %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apikey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to create response: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response :%w", err)
	}
	log.Printf("provider POST %s -> %d", url, resp.StatusCode)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("Provider returned %d , %s", resp.StatusCode, string(body))
	}
	return body, nil
}
