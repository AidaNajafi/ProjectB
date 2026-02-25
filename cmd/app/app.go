package app

import (
	"authentication/infra"
	"authentication/internal/config"
	"authentication/internal/controllers"
	"authentication/internal/logger"
	"authentication/internal/provider"
	"authentication/internal/repository"
	"authentication/internal/server"
	"authentication/internal/service"
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sony/gobreaker"
)

type App struct {
	UserApp  *server.UserServer
	HotelApp *server.HotelServer
}

func AppStart() error {

	router := gin.Default()
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("Failed to load config: %w", err)
	}
	ctx := context.Background()
	db, err := infra.StartConnectionPool(ctx, infra.ConnectionConfig{
		User:     cfg.DBuser,
		Password: cfg.DBpass,
		Host:     cfg.DBhost,
		Port:     cfg.DBport,
		Name:     cfg.DBname,
		Mode:     cfg.DBsslmode,
	})
	if err != nil {
		return fmt.Errorf("Failed to initialize the connection pool %w", err)
	}

	defer db.Close()

	repo := repository.NewPostgresStore(db)
	jwtToken := service.NewJwtService([]byte(cfg.SecretKey))

	real := provider.NewProvider(cfg.BaseURL, cfg.ApiKey, cfg.TimeOut)
	cb := provider.NewProviderBreaker(provider.BreakerConfig{
		Name:       cfg.CBname,
		MaxRequest: cfg.CBMaxRequest,
		Interval:   cfg.CBTimeOut,
		Timeout:    cfg.TimeOut,
	}, func(name string, from, to gobreaker.State) {
		log.Printf("%s: %s -> %s", name, from, to)
	})
	

	validator := validator.New()

	UserSvc := service.NewAuthService(repo, jwtToken)
	UserCtrl := controllers.NewUserController(UserSvc, validator)
	newUserServer := server.NewUserServer(UserCtrl)


	provider := provider.NewCBProvider(real, cb)
	HotelSvc := service.NewHotelService(provider, repo)
	HotelCtrl := controllers.NewHotelController(HotelSvc)
	logger := logger.NewJSON(cfg.LogLevel)
	newHotelServer := server.NewHotelServer(HotelCtrl, jwtToken, logger)

	newUserServer.SetUpUserRoutes(router)
	newHotelServer.SetUpHotelRoutes(router)
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Starting server on port: %s", addr)
	return router.Run(addr)

}
