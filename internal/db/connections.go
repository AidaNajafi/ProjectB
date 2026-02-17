package db

import (
	"authentication/internal/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Init(cfg *config.Config) (*sql.DB, error) {

	connectionStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBuser, cfg.DBpass, cfg.DBhost,
		cfg.DBport, cfg.DBname, cfg.DBsslmode)
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
