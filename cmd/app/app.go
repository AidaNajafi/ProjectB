package app

import (
	"authentication/internal/config"
	"authentication/internal/controllers"
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

	cfg, err := config.LoadConfig("../config.yaml")
	if err != nil {
		return fmt.Errorf("loading config failed: %w", err)
	}

	repo := repository.NewRepo(cfg.FilePath)
	if err := repo.Init(); err != nil {
		return fmt.Errorf("Header creation failed: %w", err)
	}
	jwtToken := service.NewJwtService([]byte(cfg.SecretKey))
	svc := service.NewAuthService(repo, jwtToken)

	validator := validator.New()
	ctrl := controllers.NewController(svc, validator)
	newS := server.NewServer(ctrl)

	router := newS.SetUpRoutes()

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Starting server on port: %s", addr)
	return router.Run(addr)

}
