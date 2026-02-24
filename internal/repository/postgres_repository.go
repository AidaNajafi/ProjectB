package repository

import (
	"authentication/internal/provider"
	"context"
	"database/sql"
	"fmt"
	"log"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (ps *PostgresStore) CreateUser(u UserInfo) error {
	_, err := ps.db.Exec("INSERT INTO users (name, username, email, userphone,  password) VALUES ($1, $2, $3, $4, $5)", u.Name, u.Username, u.Email, u.UserPhone, u.Password)
	if err != nil {
		log.Println("failed creating user")
		return fmt.Errorf("Failed to insert new user: %w", err)
	}
	return nil
}

func (ps *PostgresStore) GetUserByUsername(username string) (*UserInfo, error) {
	var user UserInfo
	err := ps.db.QueryRow("SELECT id, name, email, userphone, password FROM users WHERE username = $1", username).Scan(&user.ID, &user.Name, &user.Email, &user.UserPhone, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to retrieve user with this username : %w", err)
	}
	return &user, nil
}

func (ps *PostgresStore) CreatePendingReservation(ctx context.Context, p ReservationRequest) (int64, error) {
	const q = `
	INSERT INTO reservation
		(user_id, user_phone, room_id, date_from, date_to, reservation_status)
	VALUES
		($1, $2, $3, $4, $5, $6)
	RETURNING id
	`
	var id int64
	err := ps.db.QueryRowContext(ctx, q,
		p.UserID,
		p.UserPhone,
		p.RoomID,
		p.DateFrom,
		p.DateTo,
		"pending",
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Failed to create pending reservation: %w", err)
	}
	return id, nil
}

func (ps *PostgresStore) ConfirmedReservation(ctx context.Context, id int64, pr ReservationResponse) error {
	const q = `
	UPDATE reservation
	SET reservation_status='confirmed',
		provider_id=$2,
		hotel_id=$3
	WHERE id=$1 AND reservation_status='pending'
	`
	_, err := ps.db.ExecContext(ctx, q, id, pr.ProviderID, pr.HotelID)
	if err != nil {
		return fmt.Errorf("Failed to confirm reservation: %w", err)
	}
	return nil
}

func (ps *PostgresStore) FailedReservation(ctx context.Context, id int64, pr ReservationResponse, reason provider.FailureReason) error {
	const q = `
	UPDATE reservation
	SET reservation_status='failed',
		provider_id=$2,
		hotel_id=$3
	WHERE id=$1 AND reservation_status='pending'
	`
	_, err := ps.db.ExecContext(ctx, q, id, pr.ProviderID, pr.HotelID)
	if err != nil {
		return fmt.Errorf("Failed to mark failure for reservation: %w", err)
	}
	log.Println("Reservation Failed for id %d: %s", id, reason)
	return nil
}
