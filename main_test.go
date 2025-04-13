package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dw/firebase-studio-api-server-test/task/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	setupRoutes(r)
	return r
}

func TestTaskManager(t *testing.T) {
	router := setupRouter()

	// Test Create Task
	t.Run("Create Task", func(t *testing.T) {
		taskReq := models.TaskRequest{
			Title:       "Integration Test Task",
			Description: "Testing task creation through main app",
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

		// Store task ID for subsequent tests
		taskID := response.ID

		// Test Get Task
		t.Run("Get Task", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/tasks/"+taskID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var getResponse models.Task
			err := json.Unmarshal(w.Body.Bytes(), &getResponse)
			assert.NoError(t, err)
			assert.Equal(t, response.ID, getResponse.ID)
			assert.Equal(t, response.Title, getResponse.Title)
		})

		// Test Update Task
		t.Run("Update Task", func(t *testing.T) {
			updateReq := models.TaskRequest{
				Title:       "Updated Integration Test Task",
				Description: "Updated task description",
				Status:      "completed",
			}

			jsonData, _ := json.Marshal(updateReq)
			req, _ := http.NewRequest("PUT", "/tasks/"+taskID, bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var updateResponse models.Task
			err := json.Unmarshal(w.Body.Bytes(), &updateResponse)
			assert.NoError(t, err)
			assert.Equal(t, updateReq.Title, updateResponse.Title)
			assert.Equal(t, updateReq.Description, updateResponse.Description)
			assert.Equal(t, updateReq.Status, updateResponse.Status)
		})

		// Test Get All Tasks
		t.Run("Get All Tasks", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/tasks", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var tasks []models.Task
			err := json.Unmarshal(w.Body.Bytes(), &tasks)
			assert.NoError(t, err)
			assert.NotEmpty(t, tasks)
			assert.GreaterOrEqual(t, len(tasks), 1)
		})

		// Test Delete Task
		t.Run("Delete Task", func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", "/tasks/"+taskID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNoContent, w.Code)

			// Verify task is deleted by trying to get it
			req, _ = http.NewRequest("GET", "/tasks/"+taskID, nil)
			w = httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)
		})
	})
}

func TestHealthEndpoint(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "ok", response["status"])
}
