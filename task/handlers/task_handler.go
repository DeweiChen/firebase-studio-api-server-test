package handlers

import (
	"net/http"

	"github.com/dw/firebase-studio-api-server-test/task/models"
	"github.com/dw/firebase-studio-api-server-test/task/repositories"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	repo repositories.TaskRepository
}

func NewTaskHandler(repo repositories.TaskRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

// CreateTask handles the creation of a new task
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var taskReq models.TaskRequest
	if err := c.ShouldBindJSON(&taskReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := &models.Task{
		Title:       taskReq.Title,
		Description: taskReq.Description,
		Status:      taskReq.Status,
	}

	createdTask, err := h.repo.Create(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, createdTask)
}

// GetTask handles retrieving a task by ID
func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get task"})
		return
	}

	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// GetAllTasks handles retrieving all tasks
func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	tasks, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// UpdateTask handles updating an existing task
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	existingTask, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get task"})
		return
	}

	if existingTask == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var taskReq models.TaskRequest
	if err := c.ShouldBindJSON(&taskReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingTask.Title = taskReq.Title
	existingTask.Description = taskReq.Description
	existingTask.Status = taskReq.Status

	updatedTask, err := h.repo.Update(existingTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

// DeleteTask handles deleting a task
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	err := h.repo.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.Status(http.StatusNoContent)
}
