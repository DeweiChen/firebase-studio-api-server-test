package repositories

import (
	"testing"
	"time"

	"github.com/dw/firebase-studio-api-server-test/task/models"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryTaskRepository(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	// Test Create
	task := &models.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "pending",
	}

	createdTask, err := repo.Create(task)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdTask.ID)
	assert.Equal(t, task.Title, createdTask.Title)
	assert.Equal(t, task.Description, createdTask.Description)
	assert.Equal(t, task.Status, createdTask.Status)

	// Test GetByID
	retrievedTask, err := repo.GetByID(createdTask.ID)
	assert.NoError(t, err)
	assert.Equal(t, createdTask, retrievedTask)

	// Test GetAll
	tasks, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, createdTask, tasks[0])

	// Test Update
	updatedTask := &models.Task{
		ID:          createdTask.ID,
		Title:       "Updated Task",
		Description: "Updated Description",
		Status:      "completed",
		CreatedAt:   createdTask.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	updated, err := repo.Update(updatedTask)
	assert.NoError(t, err)
	assert.Equal(t, updatedTask.Title, updated.Title)
	assert.Equal(t, updatedTask.Description, updated.Description)
	assert.Equal(t, updatedTask.Status, updated.Status)

	// Test Delete
	err = repo.Delete(createdTask.ID)
	assert.NoError(t, err)

	// Verify task is deleted
	deletedTask, err := repo.GetByID(createdTask.ID)
	assert.NoError(t, err)
	assert.Nil(t, deletedTask)

	// Test GetByID with non-existent task
	nonExistentTask, err := repo.GetByID("non-existent-id")
	assert.NoError(t, err)
	assert.Nil(t, nonExistentTask)
}
