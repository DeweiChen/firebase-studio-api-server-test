package handlers

import (
	"net/http"

	"github.com/dw/firebase-studio-api-server-test/task/models"
	"github.com/dw/firebase-studio-api-server-test/task/repositories"
	"github.com/gin-gonic/gin"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

type TaskHandler struct {
	repo repositories.TaskRepository
}

func NewTaskHandler(repo repositories.TaskRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

// @Summary      Create a new task
// @Description  Create a new task with the provided details
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        task  body      models.TaskRequest  true  "Task details"
// @Success      201   {object}  models.Task
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var taskReq models.TaskRequest
	if err := c.ShouldBindJSON(&taskReq); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	task := &models.Task{
		Title:       taskReq.Title,
		Description: taskReq.Description,
		Status:      taskReq.Status,
	}

	createdTask, err := h.repo.Create(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, createdTask)
}

// @Summary      Get a task by ID
// @Description  Get a specific task by its ID
// @Tags         tasks
// @Produce      json
// @Param        id   path      string  true  "Task ID"
// @Success      200  {object}  models.Task
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tasks/{id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get task"})
		return
	}

	if task == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// @Summary      Get all tasks
// @Description  Get a list of all tasks
// @Tags         tasks
// @Produce      json
// @Success      200  {array}   models.Task
// @Failure      500  {object}  ErrorResponse
// @Router       /tasks [get]
func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	tasks, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// @Summary      Update a task
// @Description  Update an existing task with new details
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id    path      string           true  "Task ID"
// @Param        task  body      models.TaskRequest  true  "Updated task details"
// @Success      200   {object}  models.Task
// @Failure      400   {object}  ErrorResponse
// @Failure      404   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	existingTask, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get task"})
		return
	}

	if existingTask == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Task not found"})
		return
	}

	var taskReq models.TaskRequest
	if err := c.ShouldBindJSON(&taskReq); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	existingTask.Title = taskReq.Title
	existingTask.Description = taskReq.Description
	existingTask.Status = taskReq.Status

	updatedTask, err := h.repo.Update(existingTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

// @Summary      Delete a task
// @Description  Delete a task by its ID
// @Tags         tasks
// @Param        id   path      string  true  "Task ID"
// @Success      204  "No Content"
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	err := h.repo.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete task"})
		return
	}

	c.Status(http.StatusNoContent)
}
