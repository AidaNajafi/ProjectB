package app

import (
	"authentication/internal/config"
	"authentication/internal/controllers"
	"authentication/internal/repository"
	"authentication/internal/server"
	"authentication/internal/service"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

type App struct {
	app *server.Server
}

func AppStart(flags []string) error {
	err := godotenv.Load("my.env")
	if err != nil {
		log.Println("Failed to find env file")
	}
	cfg, err := config.LoadConfig(flags)
	if err != nil {
		return fmt.Errorf("loadConfig failed: %w", err)
	}
	addr := fmt.Sprintf(":%d", cfg.Port)
	repo := repository.NewCSVRepo(cfg.FilePath)
	if err := repo.Init(); err != nil {
		return fmt.Errorf("Header creation failed: %w", err)
	}
	svc := service.NewAuthService(repo)
	srv := controllers.NewController(svc)
	newS := server.NewServer(srv)
	log.Printf("Starting server on port: %s", addr)
	router := newS.SetUpRoutes(svc)
	return router.Run(addr)

}
