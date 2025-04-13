package repositories

import (
	"sync"
	"time"

	"github.com/dw/firebase-studio-api-server-test/task/models"
	"github.com/google/uuid"
)

// TaskRepository defines the interface for task operations
type TaskRepository interface {
	Create(task *models.Task) (*models.Task, error)
	GetByID(id string) (*models.Task, error)
	GetAll() ([]*models.Task, error)
	Update(task *models.Task) (*models.Task, error)
	Delete(id string) error
}

// InMemoryTaskRepository implements TaskRepository using in-memory storage
type InMemoryTaskRepository struct {
	tasks map[string]*models.Task
	mu    sync.RWMutex
}

// NewInMemoryTaskRepository creates a new instance of InMemoryTaskRepository
func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks: make(map[string]*models.Task),
	}
}

// Create adds a new task to the repository
func (r *InMemoryTaskRepository) Create(task *models.Task) (*models.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.ID = uuid.New().String()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	r.tasks[task.ID] = task
	return task, nil
}

// GetByID retrieves a task by its ID
func (r *InMemoryTaskRepository) GetByID(id string) (*models.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, nil
	}
	return task, nil
}

// GetAll retrieves all tasks
func (r *InMemoryTaskRepository) GetAll() ([]*models.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*models.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// Update modifies an existing task
func (r *InMemoryTaskRepository) Update(task *models.Task) (*models.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		return nil, nil
	}

	task.UpdatedAt = time.Now()
	r.tasks[task.ID] = task
	return task, nil
}

// Delete removes a task by its ID
func (r *InMemoryTaskRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return nil
	}

	delete(r.tasks, id)
	return nil
}
