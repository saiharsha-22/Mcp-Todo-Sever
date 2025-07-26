package tools

import (
	"context"
	"fmt"
	"log"
	"todo-app/internal/models"
	"todo-app/internal/util"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ListTodos tool implementation
func ListTodos(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[struct{}]) (*mcp.CallToolResultFor[struct{}], error) {
	log.Println("[ListTodos] tool called")
	todos, err := models.LoadTodos(ctx)
	if err != nil {
		log.Printf("[ListTodos] Error loading todos: %v\n", err)
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error loading todo list: " + err.Error()},
			},
			IsError: true,
		}, nil
	}

	var lines []string
	for _, todo := range todos {
		if todo.Status == "pending" {
			lines = append(lines, fmt.Sprintf("%s [%s] (Due %s): %s", todo.JiraID, todo.Severity, todo.DueDate, todo.Description))
		}
	}
	if len(lines) == 0 {
		lines = append(lines, "No pending todos!")
	}

	resultText := "Pending TODOS:\n" + util.JoinLines(lines)
	log.Println("[ListTodos] returning result:", resultText)

	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{&mcp.TextContent{Text: resultText}},
	}, nil
}

// Prompt for ListTodos tool
func PromptListTodos(ctx context.Context, ss *mcp.ServerSession, params *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "List your pending Jira todos.",
		Messages: []*mcp.PromptMessage{
			{Role: "user", Content: &mcp.TextContent{Text: "Show my pending todos."}},
		},
	}, nil
}
