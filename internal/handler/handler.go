package handler

import (
	petEntity "black-pearl/backend-hackathon/internal/domain/pet/entity"
	taskEntity "black-pearl/backend-hackathon/internal/domain/task/entity"
	"black-pearl/backend-hackathon/internal/handler/dto"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskServiceInterface interface {
	Task(ctx context.Context, taskID int64) (*taskEntity.Task, error)
}

type PetServiceInterface interface {
	GetPetByUserID(ctx context.Context, userID int) (*petEntity.Pet, error)
	SetName(ctx context.Context, name string, userID int) error
}

type Handler struct {
	taskSvc TaskServiceInterface
	petSvc  PetServiceInterface
}

func NewHandler(taskSvc TaskServiceInterface, petSvc PetServiceInterface) *Handler {
	return &Handler{taskSvc: taskSvc, petSvc: petSvc}
}

func (h *Handler) Register(r *gin.Engine) {
	r.POST("/pet/{id}", h.PostName)
	r.GET("/pet/{id}", h.GetPet)
	r.GET("/task/{id}", h.GetTask)
	r.POST("/tasks/submit", h.SubmitTask)
	r.GET("/tasks/progress", h.GetProgress)
}

func (h *Handler) GetPet(c *gin.Context) {
	userIDStr := c.Param("userID")
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

	resp := dto.GetPetReq{
		ID:    pet.ID,
		Name:  pet.Name,
		Age:   pet.Age,
		Exp:   pet.Exp,
		Level: pet.Level,
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

// затычка
func (h *Handler) GetTask(c *gin.Context) {

}

// затычка
func (h *Handler) SubmitTask(c *gin.Context) {
}

// затычка
func (h *Handler) GetProgress(c *gin.Context) {
}
