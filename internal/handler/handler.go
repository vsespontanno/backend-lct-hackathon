package handler

import (
	petEntity "black-pearl/backend-hackathon/internal/domain/pet/entity"
	prizeEntity "black-pearl/backend-hackathon/internal/domain/prize/entity"
	taskEntity "black-pearl/backend-hackathon/internal/domain/task/entity"
	"black-pearl/backend-hackathon/internal/handler/dto"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PrizeServiceInterface interface {
	AvailablePrizes(ctx context.Context, userID int) (*[]prizeEntity.Prize, error)
	MyPrizes(ctx context.Context, user_id int) (*[]prizeEntity.Prize, error)
}

type TaskServiceInterface interface {
	Task(ctx context.Context, taskID int) (*taskEntity.Task, error)
}

type PetServiceInterface interface {
	UpdateXP(ctx context.Context, xp int, userID int) error
	GetPetByUserID(ctx context.Context, userID int) (*petEntity.Pet, error)
	SetName(ctx context.Context, name string, userID int) error
}

type Handler struct {
	taskSvc  TaskServiceInterface
	petSvc   PetServiceInterface
	prizeSvc PrizeServiceInterface
}

func NewHandler(taskSvc TaskServiceInterface, petSvc PetServiceInterface, prizeSvc PrizeServiceInterface) *Handler {
	return &Handler{taskSvc: taskSvc, petSvc: petSvc, prizeSvc: prizeSvc}
}

func (h *Handler) Register(r *gin.Engine) {
	r.GET("/prizes/{id}/my", h.GetMyPrizes)
	r.POST("/prizes/{id}/available", h.GetAvailablePrizes)
	r.POST("/pet/xp", h.PostXP)
	r.POST("/pet/name", h.PostName)
	r.GET("/pet/{id}", h.GetPet)
	r.GET("/task/{id}", h.GetTask)
	r.POST("/tasks/submit", h.SubmitTask)
	r.GET("/tasks/progress", h.GetProgress)
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

// затычка
func (h *Handler) GetTask(c *gin.Context) {

}

// затычка
func (h *Handler) SubmitTask(c *gin.Context) {
}

// затычка
func (h *Handler) GetProgress(c *gin.Context) {
}
