package app

import (
	"authentication/internal/config"
	"authentication/internal/repository"
	"authentication/internal/server"
	"authentication/internal/service"
	"fmt"
	"log"
)

func App(flags []string) error {
	cfg, err := config.LoadConfig(flags)
	if err != nil {
		return fmt.Errorf("loadConfig failed: %w", err)
	}
	repo := repository.NewCSVRepo(cfg.FilePath)
	svc := service.NewAuthService(repo)
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Starting server on port: %s", addr)
	router := server.SetUpRoutes(svc)
	return router.Run(addr)

}
