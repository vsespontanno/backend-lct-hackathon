package app

import (
	"black-pearl/backend-hackathon/internal/config"
	"black-pearl/backend-hackathon/internal/handler"
	"black-pearl/backend-hackathon/internal/infrastructure/db"
	"black-pearl/backend-hackathon/internal/infrastructure/repository/postgres/pet"
	"black-pearl/backend-hackathon/internal/infrastructure/repository/postgres/sectionItems"
	"black-pearl/backend-hackathon/internal/infrastructure/repository/postgres/sections"
	"black-pearl/backend-hackathon/internal/infrastructure/repository/postgres/task"
	"black-pearl/backend-hackathon/internal/infrastructure/repository/postgres/theory"
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
	sectionRepo := sections.NewSectionsRepo(dataBase)
	sectionItemsRepo := sectionitems.NewSectionItemsRepo(dataBase)
	theoryRepo := theory.NewTheoryRepo(dataBase)
	taskSvc := service.NewTaskService(taskRepo)
	petSvc := service.NewPetService(petRepo, userRepo)
	sectionSvc := service.NewSectionService(sectionRepo)
	if sectionSvc == nil {
		log.Fatal("failed to initialize sectionSvc")
	}
	sectionItemsSvc := service.NewSectionItemsService(sectionItemsRepo)
	theorySvc := service.NewTheoryService(theoryRepo)
	taskHandler := handler.NewHandler(taskSvc, petSvc, sectionSvc, sectionItemsSvc, theorySvc)
	taskHandler.Register(r)
	return &App{
		engine: r,
	}
}

func (a *App) Run(addr string) error {
	return a.engine.Run(addr)
}
