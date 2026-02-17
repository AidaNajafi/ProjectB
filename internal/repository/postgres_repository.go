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
	_, err := ps.db.Exec("INSERT INTO users (name, username, email, password) VALUES ($1, $2, $3, $4)", u.Name, u.Username, u.Email, u.Password)
	if err != nil {
		log.Println("failed creating user")
		return fmt.Errorf("Failed to insert new user: %w", err)
	}
	return nil
}

func (ps *PostgresStore) GetUserByUsername(username string) (*UserInfo, error) {
	var user UserInfo
	err := ps.db.QueryRow("SELECT id, name, email, password FROM users WHERE username = $1", username).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to retrieve user with this username : %w", err)
	}
	return &user, nil
}
