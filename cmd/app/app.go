package app

import (
	"authentication/internal/config"
	"authentication/internal/controllers"
	"authentication/internal/db"
	"authentication/internal/repository"
	"authentication/internal/server"
	"authentication/internal/service"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
)

type App struct {
	app *server.Server
}

func AppStart() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("Failed to load config: %w", err)
	}

	db, err := db.Init(cfg)
	if err != nil {
		return fmt.Errorf("Failed to initialize the connection pool %w", err)
	}
	
	defer db.Close()
	repo := repository.NewPostgresStore(db)
	jwtToken := service.NewJwtService([]byte(cfg.SecretKey))
	svc := service.NewAuthService(repo, jwtToken)

	validator := validator.New()
	ctrl := controllers.NewController(svc, validator)
	newS := server.NewServer(ctrl)

	router := newS.SetUpRoutes()

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Starting server on port: %s", addr)
	return router.Run(addr)

}
