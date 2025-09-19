package handler

import (
	"black-pearl/backend-hackathon/internal/service"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	svc *service.TaskService
}

func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) Register(r *gin.Engine) {
	r.GET("/task/:id", h.GetTask)
	r.POST("/tasks/submit", h.SubmitTask)
	r.GET("/tasks/progress", h.GetProgress)
}

// затычка
func (h *TaskHandler) GetTask(c *gin.Context) {
	return
}

// затычка
func (h *TaskHandler) SubmitTask(c *gin.Context) {
	return
}

// затычка
func (h *TaskHandler) GetProgress(c *gin.Context) {
	return
}
