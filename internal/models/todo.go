package models

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
	"time"
)

var mu sync.Mutex

type TodoItem struct {
	JiraID      string    `json:"jira_id"`
	Status      string    `json:"status"`
	Severity    string    `json:"severity"`
	DueDate     string    `json:"due_date"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var TodoFilePath string

func init() {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		TodoFilePath = "/tmp/todo.json"
	} else {
		TodoFilePath = home + "/todo.json"
	}
}

func LoadTodos(ctx context.Context) ([]TodoItem, error) {
	mu.Lock()
	defer mu.Unlock()

	f, err := os.Open(TodoFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Printf("File %s does not exist, returning empty todo list\n", TodoFilePath)
			return []TodoItem{}, nil
		}
		log.Printf("Error opening file %s: %v\n", TodoFilePath, err)
		return nil, err
	}
	defer f.Close()

	var todos []TodoItem
	decoder := json.NewDecoder(f)
	if err = decoder.Decode(&todos); err != nil {
		log.Printf("Failed to decode JSON from %s: %v\n", TodoFilePath, err)
		return nil, err
	}

	log.Printf("Loaded %d todos from %s\n", len(todos), TodoFilePath)
	return todos, nil
}

func SaveTodos(ctx context.Context, todos []TodoItem) error {
	mu.Lock()
	defer mu.Unlock()

	b, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(TodoFilePath, b, 0644)
}
