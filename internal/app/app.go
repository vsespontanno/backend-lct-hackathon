package app

import (
	"black-pearl/backend-hackathon/internal/config"
	"black-pearl/backend-hackathon/internal/handler"
	"black-pearl/backend-hackathon/internal/infrastructure/db"
	"black-pearl/backend-hackathon/internal/infrastructure/repository/postgres/pet"
	"black-pearl/backend-hackathon/internal/infrastructure/repository/postgres/prize"
	"black-pearl/backend-hackathon/internal/infrastructure/repository/postgres/quiz"
	sectionitems "black-pearl/backend-hackathon/internal/infrastructure/repository/postgres/sectionItems"
	"black-pearl/backend-hackathon/internal/infrastructure/repository/postgres/sections"
	"black-pearl/backend-hackathon/internal/infrastructure/repository/postgres/task"
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
	logger.Log.Infow("initialized quiz repository", "stage", "quizRepo")
	petRepo := pet.NewPetRepo(dataBase)
	logger.Log.Infow("initialized pet repository", "stage", "petRepo")
	userRepo := user.NewUserRepo(dataBase)
	logger.Log.Infow("initialized user repository", "stage", "userRepo")
	theoryRepo := theory.NewTheoryRepo(dataBase)
	logger.Log.Infow("initialized theory repository", "stage", "theoryRepo")
	taskRepo := task.NewTaskRepo(dataBase, logger.Log)
	logger.Log.Infow("initialized task repository", "stage", "taskRepo")
	sectionItemsRepo := sectionitems.NewSectionItemsRepo(dataBase)
	logger.Log.Infow("initialized section items repository", "stage", "sectionItemsRepo")
	prizeRepo := prize.NewPrizeRepo(dataBase)
	logger.Log.Infow("initialized repositories", "stage", "repositories")
	sectionRepo := sections.NewSectionsRepo(dataBase)
	logger.Log.Infow("initialized sections repository", "stage", "sectionRepo")
	prizeSvc := service.NewPrizeService(prizeRepo, logger.Log)
	logger.Log.Infow("initialized prize service", "stage", "prizeSvc")
	petSvc := service.NewPetService(petRepo, userRepo, logger.Log)
	logger.Log.Infow("initialized pet service", "stage", "petSvc")
	quizSvc := service.NewQuizService(quizRepo, logger.Log)
	logger.Log.Infow("initialized quiz service", "stage", "quizSvc")
	sectionSvc := service.NewSectionService(sectionRepo, logger.Log)
	logger.Log.Infow("initialized section service", "stage", "sectionSvc")
	if sectionSvc == nil {
		logger.Log.Errorw("sectionSvc is nil", "stage", "NewSectionService")
		return nil
	}
	sectionItemsSvc := service.NewSectionItemsService(sectionItemsRepo, logger.Log)
	logger.Log.Infow("initialized section items service", "stage", "sectionItemsSvc")
	theorySvc := service.NewTheoryService(theoryRepo, logger.Log)
	logger.Log.Infow("initialized theory service", "stage", "theorySvc")
	taskSvc := service.NewTaskService(taskRepo, logger.Log)
	logger.Log.Infow("initialized task service", "stage", "taskSvc")
	quizHandler := handler.NewHandler(quizSvc, petSvc, sectionSvc, sectionItemsSvc, theorySvc, prizeSvc, taskSvc, logger.Log)

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
