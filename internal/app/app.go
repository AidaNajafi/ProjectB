package app

import (
	"authentication/internal/config"
	"authentication/internal/controllers"
	"authentication/internal/repository"
	"authentication/internal/server"
	"authentication/internal/service"
	"fmt"
	"log"
	"os"
)

type App struct {
	app *server.Server
}

func AppStart(flags []string) error {
	jwtV := []byte(os.Getenv("JWT_SECRET"))
	jwtSec := service.NewJwtService(jwtV)
	cfg, err := config.LoadConfig(flags)
	if err != nil {
		return fmt.Errorf("loadConfig failed: %w", err)
	}
	addr := fmt.Sprintf(":%d", cfg.Port)
	repo := repository.NewRepo(cfg.FilePath)
	if err := repo.Init(); err != nil {
		return fmt.Errorf("Header creation failed: %w", err)
	}
	svc := service.NewAuthService(repo)
	srv := controllers.NewController(svc, jwtSec)
	newS := server.NewServer(srv)
	log.Printf("Starting server on port: %s", addr)
	router := newS.SetUpRoutes(svc)
	return router.Run(addr)

}
