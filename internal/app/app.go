package app

import (
	"black-pearl/backend-hackathon/internal/config"
	"black-pearl/backend-hackathon/internal/handler"
	"black-pearl/backend-hackathon/internal/infrastructure/db"
	"black-pearl/backend-hackathon/internal/infrastructure/repository/postgres/pet"
	"black-pearl/backend-hackathon/internal/infrastructure/repository/postgres/prize"
	"black-pearl/backend-hackathon/internal/infrastructure/repository/postgres/task"
	"black-pearl/backend-hackathon/internal/infrastructure/repository/postgres/user"
	"black-pearl/backend-hackathon/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

type App struct {
	engine *gin.Engine
}

func NewApp() *App {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Printf("failed to connect to config: %v", err)
	}
	//	log.Printf("DB_HOST=%s DB_PORT=%s DB_USER=%s DB_NAME=%s", cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.DBName)
	dataBase, err := db.ConnectToPostgres(cfg.DB.Username, cfg.DB.Password, cfg.DB.DBName, cfg.DB.Host, cfg.DB.Port, cfg.DB.SSLMode)
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
	}
	r := gin.Default()
	taskRepo := task.NewTaskRepo(dataBase)
	petRepo := pet.NewPetRepo(dataBase)
	userRepo := user.NewUserRepo(dataBase)
	prizeRepo := prize.NewPrizeRepo(dataBase)
	prizeSvc := service.NewPrizeService(prizeRepo)
	taskSvc := service.NewTaskService(taskRepo)
	petSvc := service.NewPetService(petRepo, userRepo)
	taskHandler := handler.NewHandler(taskSvc, petSvc, prizeSvc)
	taskHandler.Register(r)
	return &App{
		engine: r,
	}
}

func (a *App) Run(addr string) error {
	return a.engine.Run(addr)
}
