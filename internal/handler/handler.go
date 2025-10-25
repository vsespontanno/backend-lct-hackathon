package handler

import (
	petEntity "black-pearl/backend-hackathon/internal/domain/pet/entity"
	prizeEntity "black-pearl/backend-hackathon/internal/domain/prize/entity"
	quizEntity "black-pearl/backend-hackathon/internal/domain/quiz/entity"
	sectionItemsEntity "black-pearl/backend-hackathon/internal/domain/sectionItems/entity"
	sectionEntity "black-pearl/backend-hackathon/internal/domain/sections/entity"
	taskEntity "black-pearl/backend-hackathon/internal/domain/task/entity"
	theoryEntity "black-pearl/backend-hackathon/internal/domain/theory/entity"
	"black-pearl/backend-hackathon/internal/handler/dto"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PrizeServiceInterface interface {
	AvailablePrizes(ctx context.Context, userID int) (*[]prizeEntity.Prize, error)
	MyPrizes(ctx context.Context, user_id int) (*[]prizeEntity.Prize, error)
}

type QuizServiceInterface interface {
	GetQuiz(ctx context.Context, quizID int) (*quizEntity.Quiz, error)
}

type PetServiceInterface interface {
	UpdateXP(ctx context.Context, xp int, userID int) error
	GetPetByUserID(ctx context.Context, userID int) (*petEntity.Pet, error)
	SetName(ctx context.Context, name string, userID int) error
}

type SectionServiceInterface interface {
	GetSections(ctx context.Context) (*[]sectionEntity.Sections, error)
	NewSection(ctx context.Context, title string) (*sectionEntity.Sections, error)
}

type SectionItemsServiceInterface interface {
	GetSectionItemsBySectionID(ctx context.Context, sectionID int) (*[]sectionItemsEntity.SectionItem, error)
	NewSectionItem(ctx context.Context, sectionID int, title string, isTest bool, itemId int) (*sectionItemsEntity.SectionItem, error)
}

type TheoryServiceInterface interface {
	GetTheoryByID(ctx context.Context, theoryID int) (*theoryEntity.Theory, error)
	NewTheory(ctx context.Context, title, content string) (*theoryEntity.Theory, error)
}

type TaskServiceInterface interface {
	GetTasks(ctx context.Context) (*[]taskEntity.Task, error)
}

type Handler struct {
	quizSvc         QuizServiceInterface
	petSvc          PetServiceInterface
	sectionSvc      SectionServiceInterface
	sectionItemsSvc SectionItemsServiceInterface
	theorySvc       TheoryServiceInterface
	prizeSvc        PrizeServiceInterface
	taskSvc         TaskServiceInterface
	logger          *zap.SugaredLogger
}

func NewHandler(
	quizSvc QuizServiceInterface,
	petSvc PetServiceInterface,
	sectionSvc SectionServiceInterface,
	sectionItemsSvc SectionItemsServiceInterface,
	theorySvc TheoryServiceInterface,
	prizeSvc PrizeServiceInterface,
	taskSvc TaskServiceInterface,
	logger *zap.SugaredLogger,
) *Handler {
	return &Handler{
		quizSvc:         quizSvc,
		petSvc:          petSvc,
		sectionSvc:      sectionSvc,
		sectionItemsSvc: sectionItemsSvc,
		theorySvc:       theorySvc,
		prizeSvc:        prizeSvc,
		taskSvc:         taskSvc,
		logger:          logger,
	}
}

func (h *Handler) Register(r *gin.Engine) {
	r.GET("/quiz/:id", h.GetQuiz)

	// новые ручки
	r.GET("/sections", h.GetSectionsWithItems)
	r.POST("/sections", h.NewSection)
	r.GET("/sections/:id/items", h.GetSectionItems)
	r.POST("/sections/:id/items", h.NewSectionItem)

	r.GET("/theory/:id", h.GetTheory)
	r.POST("/theory", h.NewTheory)

	r.GET("/prizes/:id/my", h.GetMyPrizes)
	r.GET("/prizes/:id/available", h.GetAvailablePrizes)
	r.POST("/pet/xp", h.PostXP)
	r.POST("/pet/name", h.PostName)
	r.GET("/pet/:id", h.GetPet)

	r.GET("/tasks/daily", h.GetDailyTasks)
}

func (h *Handler) GetPet(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		h.logger.Errorw("Invalid user ID", "error", err, "stage", "GetPet.Atoi")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	pet, err := h.petSvc.GetPetByUserID(c.Request.Context(), userID)
	if err != nil {
		h.logger.Errorw("failed to get pet", "error", err, "stage", "GetPet.GetPetByUserID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if pet == nil {
		h.logger.Errorw("pet not found", "stage", "GetPet.GetPetByUserID")
		c.JSON(http.StatusNotFound, gin.H{"error": "pet not found"})
		return
	}

	resp := dto.GetPetResp{
		ID:   pet.ID,
		Name: pet.Name,
		Age:  pet.Age,
		Exp:  pet.Exp,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) PostName(c *gin.Context) {
	var req dto.SetPetNameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("failed to bind request", "error", err, "stage", "PostName.ShouldBindJSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.petSvc.SetName(context.Background(), req.Name, req.UserID)
	if err != nil {
		h.logger.Errorw("failed to set pet name", "error", err, "stage", "PostName.SetName")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// ------------------- SECTIONS -------------------

func (h *Handler) GetSectionsWithItems(c *gin.Context) {
	ctx := c.Request.Context()

	sections, err := h.sectionSvc.GetSections(ctx)
	if err != nil || sections == nil {
		h.logger.Errorw("failed to get sections", "error", err, "stage", "GetSectionsWithItems.GetSections")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []dto.SectionWithItemsResp
	for _, section := range *sections {
		items, err := h.sectionItemsSvc.GetSectionItemsBySectionID(ctx, section.ID)
		if err != nil {
			h.logger.Errorw("failed to get section items", "error", err, "stage", "GetSectionsWithItems.GetSectionItemsBySectionID")
		}
		if items == nil {
			items = &[]sectionItemsEntity.SectionItem{}
		}

		result = append(result, dto.SectionWithItemsResp{
			ID:    section.ID,
			Title: section.Title,
			Items: *items,
		})
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) NewSection(c *gin.Context) {
	var req dto.NewSectionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("failed to bind request", "error", err, "stage", "NewSection.ShouldBindJSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	section, err := h.sectionSvc.NewSection(c.Request.Context(), req.Title)
	if err != nil {
		h.logger.Errorw("failed to create section", "error", err, "stage", "NewSection.NewSection")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, section)
}

// ------------------- SECTION ITEMS -------------------

func (h *Handler) GetSectionItems(c *gin.Context) {
	sectionIDStr := c.Param("id")
	sectionID, err := strconv.Atoi(sectionIDStr)
	if err != nil {
		h.logger.Errorw("invalid section ID", "error", err, "stage", "GetSectionItems.Atoi")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid section ID"})
		return
	}
	items, err := h.sectionItemsSvc.GetSectionItemsBySectionID(c.Request.Context(), sectionID)
	if err != nil {
		h.logger.Errorw("failed to get section items", "error", err, "stage", "GetSectionItems.GetSectionItemsBySectionID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) NewSectionItem(c *gin.Context) {
	sectionIDStr := c.Param("id")
	sectionID, err := strconv.Atoi(sectionIDStr)
	if err != nil {
		h.logger.Errorw("invalid section ID", "error", err, "stage", "NewSectionItem.Atoi")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid section ID"})
		return
	}

	var req dto.NewSectionItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("failed to bind request", "error", err, "stage", "NewSectionItem.ShouldBindJSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.sectionItemsSvc.NewSectionItem(c.Request.Context(), sectionID, req.Title, req.IsTest, req.ItemID)
	if err != nil {
		h.logger.Errorw("failed to create section item", "error", err, "stage", "NewSectionItem.NewSectionItem")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// ------------------- THEORY -------------------

func (h *Handler) GetTheory(c *gin.Context) {
	idStr := c.Param("id")
	theoryID, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Errorw("invalid theory ID", "error", err, "stage", "GetTheory.Atoi")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid theory ID"})
		return
	}
	theory, err := h.theorySvc.GetTheoryByID(c.Request.Context(), theoryID)
	if err != nil {
		h.logger.Errorw("failed to get theory", "error", err, "stage", "GetTheory.GetTheoryByID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if theory == nil {
		h.logger.Errorw("theory not found", "error", err, "stage", "GetTheory.GetTheoryByID")
		c.JSON(http.StatusNotFound, gin.H{"error": "theory not found"})
		return
	}
	c.JSON(http.StatusOK, theory)
}

func (h *Handler) NewTheory(c *gin.Context) {
	var req dto.NewTheoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("failed to bind request", "error", err, "stage", "NewTheory.ShouldBindJSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	theory, err := h.theorySvc.NewTheory(c.Request.Context(), req.Title, req.Content)
	if err != nil {
		h.logger.Errorw("failed to create theory", "error", err, "stage", "NewTheory.NewTheory")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, theory)
}

// -------------------- TASKS -------------------

func (h *Handler) PostXP(c *gin.Context) {
	var req dto.SendXPReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("failed to bind request", "error", err, "stage", "PostXP.ShouldBindJSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.petSvc.UpdateXP(context.Background(), req.Exp, req.UserID)
	if err != nil {
		h.logger.Errorw("failed to update XP", "error", err, "stage", "PostXP.UpdateXP")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) GetMyPrizes(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		h.logger.Errorw("invalid user ID", "error", err, "stage", "GetMyPrizes.Atoi")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prizes, err := h.prizeSvc.MyPrizes(context.Background(), userID)
	if err != nil {
		h.logger.Errorw("failed to get my prizes", "error", err, "stage", "GetMyPrizes.MyPrizes")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var resp dto.GetPrizesResp
	resp.Prizes = *prizes

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetAvailablePrizes(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		h.logger.Errorw("invalid user ID", "error", err, "stage", "GetAvailablePrizes.Atoi")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prizes, err := h.prizeSvc.AvailablePrizes(context.Background(), userID)
	if err != nil {
		h.logger.Errorw("failed to get available prizes", "error", err, "stage", "GetAvailablePrizes.AvailablePrizes")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var resp dto.GetPrizesResp
	resp.Prizes = *prizes

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetQuiz(c *gin.Context) {
	quizIDStr := c.Param("id")
	quizID, err := strconv.Atoi(quizIDStr)
	if err != nil {
		h.logger.Errorw("invalid quiz ID", "error", err, "stage", "GetQuiz.Atoi")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID"})
		return
	}

	quiz, err := h.quizSvc.GetQuiz(c.Request.Context(), quizID)
	if err != nil {
		h.logger.Errorw("failed to get quiz", "error", err, "stage", "GetQuiz.GetQuizByID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if quiz == nil {
		h.logger.Errorw("quiz not found", "error", err, "stage", "GetQuiz.GetQuizByID")
		c.JSON(http.StatusNotFound, gin.H{"error": "quiz not found"})
		return
	}

	c.JSON(http.StatusOK, quiz)
}

func (h *Handler) GetDailyTasks(c *gin.Context) {
	tasks, err := h.taskSvc.GetTasks(c.Request.Context())
	if err != nil {
		h.logger.Errorw("failed to get tasks", "error", err, "stage", "GetTasks.GetTasks")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	today := time.Now().Day()
	var result []taskEntity.Task
	for _, task := range *tasks {
		if task.ID == today {
			result = append(result, task)
		}
	}

	c.JSON(http.StatusOK, result)
}
