package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Provider struct {
	BaseURL string
	ApiKey  string
	Client  *http.Client
	TimeOut time.Duration
}

func NewProvider(baseUrl, apiKey string, timeout time.Duration) *Provider {
	return &Provider{BaseURL: baseUrl, ApiKey: apiKey, TimeOut: timeout}
}

type RoomDetail struct {
	ID          int64  `json:"id"`
	HotelID     int64  `json:"hotel_id"`
	Name        string `json:"name"`
	Price       int64  `json:"price"`
	Currency    string `json:"currency"`
	Capacity    int64  `json:"capacity"`
	BedType     string `json:"bed_type"`
	SizeSQM     int64  `json:"size_sqm"`
	IsAvailable bool   `json:"is_available"`
}

type HotelByDate struct {
	Count  int64  `json:"count"`
	Date   string `json:"date"`
	Hotels []HotelDetail
}

func (p *Provider) HotelByDate(ctx context.Context, date string) (HotelByDate, error) {
	url := fmt.Sprintf("%s/hotels?date=%s", p.BaseURL, date)

	body, err := p.DoGetRequest(ctx, url, p.ApiKey)
	if err != nil {
		return HotelByDate{}, fmt.Errorf("Failed to connect to provider")
	}
	var hotel HotelByDate
	if err := json.Unmarshal(body, &hotel); err != nil {
		return HotelByDate{}, fmt.Errorf("Failed to unmarshal request body to json")
	}
	return hotel, nil
}

type HotelDetail struct {
	ID             int64    `json:"id"`
	Name           string   `json:"name"`
	Status         string   `json:"status"`
	Description    string   `json:"description"`
	ImageUrl       string   `json:"image_url"`
	Address        string   `json:"address"`
	City           string   `json:"city"`
	Amenities      []string `json:"amenities"`
	Rating         float64  `json:"rating"`
	TotalRooms     int64    `json:"totalrooms"`
	AvailableRooms int64    `json:"availablerooms"`
}

func (p *Provider) HotelByID(ctx context.Context, id int64) (HotelDetail, error) {

	url := fmt.Sprintf("%s/hotels/%d", p.BaseURL, id)
	body, err := p.DoGetRequest(ctx, url, p.ApiKey)
	if err != nil {
		return HotelDetail{}, fmt.Errorf("Failed to connect to provider %w", err)
	}

	var hotelDetail HotelDetail
	if err := json.Unmarshal(body, &hotelDetail); err != nil {
		return HotelDetail{}, fmt.Errorf("Failed to unmarshal request body to json %w", err)
	}
	return hotelDetail, nil

}

type HotelRoom struct {
	Count   int64 `json:"count"`
	HotelID int64 `json:"hotelid"`
	Rooms   []RoomDetail
}

func (p *Provider) HotelRooms(ctx context.Context, id int64, date_from, date_to string) (HotelRoom, error) {
	url := fmt.Sprintf("%s/hotels/%d/rooms?date_from=%s&date_to=%s", p.BaseURL, id, date_from, date_to)
	body, err := p.DoGetRequest(ctx, url, p.ApiKey)
	if err != nil {
		return HotelRoom{}, fmt.Errorf("Failed to connect to provider")
	}
	var hotelRooms HotelRoom
	if err := json.Unmarshal(body, &hotelRooms); err != nil {
		return HotelRoom{}, fmt.Errorf("Failed to unmarshal request body to json")
	}
	return hotelRooms, nil
}

func (p *Provider) DoGetRequest(ctx context.Context, url, apiKey string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %w", err)
	}
	req.Header.Set("X-API-Key", apiKey)

	resp, err := p.Client.Do(req)
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
	ProviderID int64  `json:"id"`
	RoomID     int64  `json:"room_id"`
	HotelID    int64  `json:"hotel_id"`
	UserID     int64  `json:"user_id"`
	UserPhone  string `json:"user_phone"`
	DateFrom   string `json:"date_from"`
	DateTo     string `json:"date_to"`
	Status     string `json:"status"`
}

type ReservationRequest struct {
	RoomID    int64  `json:"room_id"`
	DateFrom  string `json:"date_from"`
	DateTo    string `json:"date_to"`
	UserID    int64  `json:"user_id"`
	UserPhone string `json:"user_phone"`
}

type ProviderError struct {
	StatusCode int
	Body       []byte
}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("provider error: status: %d, message: %s", e.StatusCode, e.Body)
}

func (p *Provider) ReserveHotel(ctx context.Context, req ReservationRequest) (ReservationResponse, error) {
	base := strings.TrimRight(p.BaseURL, "/")
	url := fmt.Sprintf("%s/reservations", base)

	var response ReservationResponse
	err := p.DoPostRequest(ctx, url, p.ApiKey, req, &response)
	if err != nil {
		return ReservationResponse{}, fmt.Errorf("Failed to start provider post server: %w", err)
	}
	return response, nil
}

func (p *Provider) DoPostRequest(ctx context.Context, url, apikey string, payload any, out any) error {
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to marshal the request body(payload) %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("Failed to create request %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apikey)

	resp, err := p.Client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to create response: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to read response :%w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &ProviderError{StatusCode: resp.StatusCode, Body: body}
	}
	log.Println(resp.StatusCode)

	if out == nil {
		return nil
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("Failed to unmarshal request body to response: %w", err)
	}

	return nil
}
