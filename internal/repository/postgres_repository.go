package repository

import (
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

func (ps *PostgresStore) CreateReservation(u ReservationInfo) error {
	_, err := ps.db.Exec("INSERT INTO reservation (user_id, user_phone, room_id, hotel_id, date_from, date_to, status) VALUES ($1, $2, $3, $4, $5, $6, $7)", u.UserID, u.UserPhone, u.RoomID, u.HotelID, u.DateFrom, u.DateTo)
	if err != nil {
		log.Println("failed creating reservation")
		return fmt.Errorf("Failed to insert new reservation: %w", err)
	}
	return nil
}
