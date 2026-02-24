package infra

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

type ConnectionConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	Mode     string
}

func StartConnectionPool(ctx context.Context, cfg ConnectionConfig) (*sql.DB, error) {

	connectionStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host,
		cfg.Port, cfg.Name, cfg.Mode)
	log.Println(connectionStr)
	var err error
	db, err = sql.Open("postgres", connectionStr)
	if err != nil {
		log.Fatal("Failed to open connection to database!", err)
	}
	db.SetMaxOpenConns(25)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("connection was not established successfully!: %w", err)
	}
	log.Println("Database connection pool is up and running.")
	return db, nil
}
