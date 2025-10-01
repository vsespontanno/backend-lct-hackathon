package app

import (
	"black-pearl/backend-hackathon/internal/config"
	"black-pearl/backend-hackathon/internal/handler"
	"black-pearl/backend-hackathon/internal/infrastructure/db"
	"black-pearl/backend-hackathon/internal/infrastructure/repository/postgres/pet"
	"black-pearl/backend-hackathon/internal/infrastructure/repository/postgres/prize"
	"black-pearl/backend-hackathon/internal/infrastructure/repository/postgres/quiz"
	"black-pearl/backend-hackathon/internal/infrastructure/repository/postgres/sectionItems"
	"black-pearl/backend-hackathon/internal/infrastructure/repository/postgres/sections"
	"black-pearl/backend-hackathon/internal/infrastructure/repository/postgres/theory"
	"black-pearl/backend-hackathon/internal/infrastructure/repository/postgres/user"
	"black-pearl/backend-hackathon/internal/service"
	"black-pearl/backend-hackathon/logger"

	"github.com/gin-gonic/gin"
)

type App struct {
	engine *gin.Engine
}

func NewApp() *App {
	logger.InitLogger()
	defer logger.Log.Sync()

	cfg, err := config.ReadConfig()
	if err != nil {
		logger.Log.Errorw("failed to read config", "error", err, "stage", "readConfig")
	}
	//	log.Printf("DB_HOST=%s DB_PORT=%s DB_USER=%s DB_NAME=%s", cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.DBName)
	dataBase, err := db.ConnectToPostgres(cfg.DB.Username, cfg.DB.Password, cfg.DB.DBName, cfg.DB.Host, cfg.DB.Port, cfg.DB.SSLMode)
	if err != nil {
		logger.Log.Errorw("failed to connect to database", "error", err, "stage", "connectToPostgres")
	}
	logger.Log.Infow("connected to database", "stage", "connectToPostgres")
	r := gin.Default()
	logger.Log.Infow("initialized gin-router", "stage", "gin.Default")
	quizRepo := quiz.NewQuizRepo(dataBase)
	petRepo := pet.NewPetRepo(dataBase)
	userRepo := user.NewUserRepo(dataBase)
	prizeRepo := prize.NewPrizeRepo(dataBase)
	logger.Log.Infow("initialized repositories", "stage", "repositories")
	prizeSvc := service.NewPrizeService(prizeRepo, logger.Log)
	petSvc := service.NewPetService(petRepo, userRepo, logger.Log)
	logger.Log.Infow("initialized services", "stage", "services")
	sectionRepo := sections.NewSectionsRepo(dataBase)
	sectionItemsRepo := sectionitems.NewSectionItemsRepo(dataBase)
	theoryRepo := theory.NewTheoryRepo(dataBase)
	quizSvc := service.NewQuizService(quizRepo, logger.Log)
	sectionSvc := service.NewSectionService(sectionRepo, logger.Log)
	if sectionSvc == nil {
		logger.Log.Errorw("sectionSvc is nil", "stage", "NewSectionService")
		return nil
	}
	sectionItemsSvc := service.NewSectionItemsService(sectionItemsRepo, logger.Log)
	theorySvc := service.NewTheoryService(theoryRepo, logger.Log)
	quizHandler := handler.NewHandler(quizSvc, petSvc, sectionSvc, sectionItemsSvc, theorySvc, prizeSvc, logger.Log)

	logger.Log.Infow("initialized handlers", "stage", "handlers")
	quizHandler.Register(r)
	return &App{
		engine: r,
	}
}

func (a *App) Run(addr string) error {
	logger.Log.Infow("starting server", "stage", "Run", "addr", addr)
	return a.engine.Run(addr)
}
