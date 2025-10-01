package handler

import (
	petEntity "black-pearl/backend-hackathon/internal/domain/pet/entity"
	sectionItemsEntity "black-pearl/backend-hackathon/internal/domain/sectionItems/entity"
	sectionEntity "black-pearl/backend-hackathon/internal/domain/sections/entity"
	prizeEntity "black-pearl/backend-hackathon/internal/domain/prize/entity"
	taskEntity "black-pearl/backend-hackathon/internal/domain/task/entity"
	theoryEntity "black-pearl/backend-hackathon/internal/domain/theory/entity"
	"black-pearl/backend-hackathon/internal/handler/dto"
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PrizeServiceInterface interface {
	AvailablePrizes(ctx context.Context, userID int) (*[]prizeEntity.Prize, error)
	MyPrizes(ctx context.Context, user_id int) (*[]prizeEntity.Prize, error)
}

type TaskServiceInterface interface {
	GetTask(ctx context.Context, taskID int64) (*taskEntity.Task, error)
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
	GetSectionItemsBySectionID(ctx context.Context, sectionID int64) (*[]sectionItemsEntity.SectionItem, error)
	NewSectionItem(ctx context.Context, sectionID int64, title string, isTest bool, itemId int64) (*sectionItemsEntity.SectionItem, error)
}

type TheoryServiceInterface interface {
	GetTheoryByID(ctx context.Context, theoryID int64) (*theoryEntity.Theory, error)
	NewTheory(ctx context.Context, title, content string) (*theoryEntity.Theory, error)
}

type Handler struct {
	taskSvc         TaskServiceInterface
	petSvc          PetServiceInterface
	sectionSvc      SectionServiceInterface
	sectionItemsSvc SectionItemsServiceInterface
	theorySvc       TheoryServiceInterface
	prizeSvc PrizeServiceInterface
}

func NewHandler(
	taskSvc TaskServiceInterface,
	petSvc PetServiceInterface,
	sectionSvc SectionServiceInterface,
	sectionItemsSvc SectionItemsServiceInterface,
	theorySvc TheoryServiceInterface,
	prizeSvc PrizeServiceInterface,
) *Handler {
	return &Handler{
		taskSvc:         taskSvc,
		petSvc:          petSvc,
		sectionSvc:      sectionSvc,
		sectionItemsSvc: sectionItemsSvc,
		theorySvc:       theorySvc,
		prizeSvc: prizeSvc,
	}
}

func (h *Handler) Register(r *gin.Engine) {
	r.GET("/task/:id", h.GetTask)

	// новые ручки
	r.GET("/sections", h.GetSectionsWithItems)
	r.POST("/sections", h.NewSection)
	r.GET("/sections/:id/items", h.GetSectionItems)
	r.POST("/sections/:id/items", h.NewSectionItem)

	r.GET("/theory/:id", h.GetTheory)
	r.POST("/theory", h.NewTheory)
	
	r.GET("/prizes/{id}/my", h.GetMyPrizes)
	r.POST("/prizes/{id}/available", h.GetAvailablePrizes)
	r.POST("/pet/xp", h.PostXP)
	r.POST("/pet/name", h.PostName)
	r.GET("/pet/{id}", h.GetPet)
}

func (h *Handler) GetPet(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	pet, err := h.petSvc.GetPetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if pet == nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.petSvc.SetName(context.Background(), req.Name, req.UserID)
	if err != nil {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []dto.SectionWithItemsResp
	for _, section := range *sections {
		items, err := h.sectionItemsSvc.GetSectionItemsBySectionID(ctx, section.ID)
		if err != nil || items == nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	section, err := h.sectionSvc.NewSection(c.Request.Context(), req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, section)
}

// ------------------- SECTION ITEMS -------------------

func (h *Handler) GetSectionItems(c *gin.Context) {
	sectionIDStr := c.Param("id")
	sectionID, err := strconv.ParseInt(sectionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid section ID"})
		return
	}
	items, err := h.sectionItemsSvc.GetSectionItemsBySectionID(c.Request.Context(), sectionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) NewSectionItem(c *gin.Context) {
	sectionIDStr := c.Param("id")
	sectionID, err := strconv.ParseInt(sectionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid section ID"})
		return
	}

	var req dto.NewSectionItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.sectionItemsSvc.NewSectionItem(c.Request.Context(), sectionID, req.Title, req.IsTest, req.ItemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// ------------------- THEORY -------------------

func (h *Handler) GetTheory(c *gin.Context) {
	idStr := c.Param("id")
	theoryID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid theory ID"})
		return
	}
	theory, err := h.theorySvc.GetTheoryByID(c.Request.Context(), theoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if theory == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "theory not found"})
		return
	}
	c.JSON(http.StatusOK, theory)
}

func (h *Handler) NewTheory(c *gin.Context) {
	var req dto.NewTheoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	theory, err := h.theorySvc.NewTheory(c.Request.Context(), req.Title, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, theory)
}

// -------------------- TASKS -------------------

func (h *Handler) PostXP(c *gin.Context) {
	var req dto.SendXPReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.petSvc.UpdateXP(context.Background(), req.Exp, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) GetMyPrizes(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prizes, err := h.prizeSvc.MyPrizes(context.Background(), userID)
	if err != nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prizes, err := h.prizeSvc.AvailablePrizes(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var resp dto.GetPrizesResp
	resp.Prizes = *prizes

	c.JSON(http.StatusOK, resp)
}


func (h *Handler) GetTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	log.Println("taskIDStr =", taskIDStr) // правильный вариант
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := h.taskSvc.GetTask(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}
