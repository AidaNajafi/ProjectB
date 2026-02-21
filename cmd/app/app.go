package app

import (
	"authentication/internal/config"
	"authentication/internal/controllers"
	"authentication/internal/db"
	"authentication/internal/provider"
	"authentication/internal/repository"
	"authentication/internal/server"
	"authentication/internal/service"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type App struct {
	UserApp *server.UserServer
	HotelApp *server.HotelServer
}

func AppStart() error {
	router := gin.Default()
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
	UserSvc := service.NewAuthService(repo, jwtToken)

	provider:= provider.NewProvider(cfg)
	HotelSvc:= service.NewHotelService(provider)

	validator := validator.New()

	UserCtrl := controllers.NewUserController(UserSvc, validator)
	newUserServer := server.NewUserServer(UserCtrl)
	
	HotelCtrl:= controllers.NewHotelController(HotelSvc)
	newHotelServer:= server.NewHotelServer(HotelCtrl, jwtToken)

	newUserServer.SetUpUserRoutes(router)
	newHotelServer.SetUpHotelRoutes(router)
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Starting server on port: %s", addr)
	return router.Run(addr)

}
