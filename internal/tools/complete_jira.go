package tools

import (
	"context"
	"fmt"
	"log"
	"time"
	"todo-app/internal/models"
	"todo-app/internal/util"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type JiraIDArg struct {
	JiraID string `json:"jira_id"`
}

func CompleteJira(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[JiraIDArg]) (*mcp.CallToolResultFor[struct{}], error) {
	todos, err := models.LoadTodos(ctx)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	found := false
	for i, t := range todos {
		if t.JiraID == params.Arguments.JiraID {
			todos[i].Status = "completed"
			todos[i].UpdatedAt = now
			found = true
			break
		}
	}
	if err := models.SaveTodos(ctx, todos); err != nil {
		return nil, err
	}
	msg := "Jira not found!"
	if found {
		msg = fmt.Sprintf("Jira %s marked completed!", params.Arguments.JiraID)
	}
	return &mcp.CallToolResultFor[struct{}]{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil
}

func PromptCompleteJira(ctx context.Context, ss *mcp.ServerSession, params *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "Mark a Jira as completed.",
		Messages: []*mcp.PromptMessage{
			{Role: "user", Content: &mcp.TextContent{Text: "Mark JIRA-123 as completed."}},
		},
	}, nil
}

func ListCompleted(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[struct{}]) (*mcp.CallToolResultFor[struct{}], error) {
	log.Println("[ListCompleted] tool called")
	todos, err := models.LoadTodos(ctx)
	if err != nil {
		log.Printf("[ListCompleted] error loading todos: %v", err)
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error loading todo list: " + err.Error()},
			},
			IsError: true,
		}, nil
	}

	var lines []string
	for _, todo := range todos {
		if todo.Status == "completed" {
			lines = append(lines, fmt.Sprintf("%s (completed at %s)", todo.JiraID, todo.UpdatedAt.Format("2006-01-02")))
		}
	}

	if len(lines) == 0 {
		lines = append(lines, "No completed JIRAs!")
	}

	resultText := "Completed JIRAs:\n" + util.JoinLines(lines)
	log.Println("[ListCompleted] returning result:", resultText)

	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultText},
		},
	}, nil
}

func PromptListCompleted(ctx context.Context, ss *mcp.ServerSession, params *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "List completed Jira tasks.",
		Messages: []*mcp.PromptMessage{
			{Role: "user", Content: &mcp.TextContent{Text: "Show completed Jiras."}},
		},
	}, nil
}
