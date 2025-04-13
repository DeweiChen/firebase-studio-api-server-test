package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dw/firebase-studio-api-server-test/task/models"
	"github.com/dw/firebase-studio-api-server-test/task/repositories"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() (*gin.Engine, *repositories.InMemoryTaskRepository) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	repo := repositories.NewInMemoryTaskRepository()
	handler := NewTaskHandler(repo)

	router.POST("/tasks", handler.CreateTask)
	router.GET("/tasks", handler.GetAllTasks)
	router.GET("/tasks/:id", handler.GetTask)
	router.PUT("/tasks/:id", handler.UpdateTask)
	router.DELETE("/tasks/:id", handler.DeleteTask)

	return router, repo
}

func TestTaskHandlers(t *testing.T) {
	router, repo := setupTestRouter()

	// Test CreateTask
	t.Run("CreateTask", func(t *testing.T) {
		taskReq := models.TaskRequest{
			Title:       "Test Task",
			Description: "Test Description",
			Status:      "pending",
		}

		jsonData, _ := json.Marshal(taskReq)
		req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Task
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, taskReq.Title, response.Title)
		assert.Equal(t, taskReq.Description, response.Description)
		assert.Equal(t, taskReq.Status, response.Status)
		assert.NotEmpty(t, response.ID)
	})

	// Test GetAllTasks
	t.Run("GetAllTasks", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/tasks", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []models.Task
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)
	})

	// Test GetTask
	t.Run("GetTask", func(t *testing.T) {
		tasks, _ := repo.GetAll()
		taskID := tasks[0].ID

		req, _ := http.NewRequest("GET", "/tasks/"+taskID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Task
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, taskID, response.ID)
	})

	// Test UpdateTask
	t.Run("UpdateTask", func(t *testing.T) {
		tasks, _ := repo.GetAll()
		taskID := tasks[0].ID

		updateReq := models.TaskRequest{
			Title:       "Updated Task",
			Description: "Updated Description",
			Status:      "completed",
		}

		jsonData, _ := json.Marshal(updateReq)
		req, _ := http.NewRequest("PUT", "/tasks/"+taskID, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Task
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, updateReq.Title, response.Title)
		assert.Equal(t, updateReq.Description, response.Description)
		assert.Equal(t, updateReq.Status, response.Status)
	})

	// Test DeleteTask
	t.Run("DeleteTask", func(t *testing.T) {
		tasks, _ := repo.GetAll()
		taskID := tasks[0].ID

		req, _ := http.NewRequest("DELETE", "/tasks/"+taskID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify task is deleted
		task, _ := repo.GetByID(taskID)
		assert.Nil(t, task)
	})
}
