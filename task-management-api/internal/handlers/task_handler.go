package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/task-management-api/internal/models"
	"github.com/task-management-api/internal/repository"
)

type TaskHandler struct {
	repo repository.Repository
}

func NewTaskHandler(repo repository.Repository) *TaskHandler {
	return &TaskHandler{
		repo: repo,
	}
}

// CreateTask handles task creation
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Get userID from JWT token
	task.UserID = 1 // Temporary

	if err := h.repo.CreateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GetTask handles getting a single task
func (h *TaskHandler) GetTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := h.repo.GetTaskByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// ListTasks handles listing all tasks for a user
func (h *TaskHandler) ListTasks(c *gin.Context) {
	// TODO: Get userID from JWT token
	userID := uint(1) // Temporary

	tasks, err := h.repo.ListTasks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
