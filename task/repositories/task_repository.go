package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/dw/firebase-studio-api-server-test/task/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
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

// RedisTaskRepository implements TaskRepository using Redis storage
type RedisTaskRepository struct {
	client *redis.Client
	ctx    context.Context
	mu     sync.RWMutex
}

// NewRedisTaskRepository creates a new instance of RedisTaskRepository
func NewRedisTaskRepository() (*RedisTaskRepository, error) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379" // Default fallback
	}
	fmt.Println("redisURL", redisURL)

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	ctx := context.Background()

	// Test the connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisTaskRepository{
		client: client,
		ctx:    ctx,
	}, nil
}

// Create adds a new task to Redis
func (r *RedisTaskRepository) Create(task *models.Task) (*models.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.ID = uuid.New().String()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	// Convert task to JSON
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	// Store task in Redis hash
	err = r.client.HSet(r.ctx, "tasks", task.ID, taskJSON).Err()
	if err != nil {
		return nil, err
	}

	return task, nil
}

// GetByID retrieves a task by its ID from Redis
func (r *RedisTaskRepository) GetByID(id string) (*models.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	taskJSON, err := r.client.HGet(r.ctx, "tasks", id).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var task models.Task
	err = json.Unmarshal([]byte(taskJSON), &task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

// GetAll retrieves all tasks from Redis
func (r *RedisTaskRepository) GetAll() ([]*models.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasksMap, err := r.client.HGetAll(r.ctx, "tasks").Result()
	if err != nil {
		return nil, err
	}

	tasks := make([]*models.Task, 0, len(tasksMap))
	for _, taskJSON := range tasksMap {
		var task models.Task
		err := json.Unmarshal([]byte(taskJSON), &task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

// Update modifies an existing task in Redis
func (r *RedisTaskRepository) Update(task *models.Task) (*models.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if task exists
	exists, err := r.client.HExists(r.ctx, "tasks", task.ID).Result()
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	task.UpdatedAt = time.Now()

	// Convert task to JSON
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	// Update task in Redis
	err = r.client.HSet(r.ctx, "tasks", task.ID, taskJSON).Err()
	if err != nil {
		return nil, err
	}

	return task, nil
}

// Delete removes a task by its ID from Redis
func (r *RedisTaskRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if task exists
	exists, err := r.client.HExists(r.ctx, "tasks", id).Result()
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	// Delete task from Redis
	return r.client.HDel(r.ctx, "tasks", id).Err()
}
