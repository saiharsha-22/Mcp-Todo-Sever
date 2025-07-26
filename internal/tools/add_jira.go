package tools

import (
	"context"
	"fmt"
	"time"
	"todo-app/internal/models"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type AddJiraArgs struct {
	JiraID      string `json:"jira_id"`
	Severity    string `json:"severity"`
	DueDate     string `json:"due_date"`
	Description string `json:"description"`
}

func AddJira(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[AddJiraArgs]) (*mcp.CallToolResultFor[struct{}], error) {
	todos, err := models.LoadTodos(ctx)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	item := models.TodoItem{
		JiraID:      params.Arguments.JiraID,
		Status:      "pending",
		Severity:    params.Arguments.Severity,
		DueDate:     params.Arguments.DueDate,
		Description: params.Arguments.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	todos = append(todos, item)
	if err := models.SaveTodos(ctx, todos); err != nil {
		return nil, err
	}
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Jira %s added to your todo list!", params.Arguments.JiraID)}},
	}, nil
}

func PromptAddJira(ctx context.Context, ss *mcp.ServerSession, params *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "Add a Jira todo to your list.",
		Messages: []*mcp.PromptMessage{
			{Role: "user", Content: &mcp.TextContent{Text: "Add a new Jira item with severity, due date, and a short description."}},
		},
	}, nil
}
