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

func RegressJira(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[JiraIDArg]) (*mcp.CallToolResultFor[struct{}], error) {
	todos, err := models.LoadTodos(ctx)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	found := false
	for i, t := range todos {
		if t.JiraID == params.Arguments.JiraID {
			todos[i].Status = "regressed"
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
		msg = fmt.Sprintf("Jira %s marked regressed!", params.Arguments.JiraID)
	}
	return &mcp.CallToolResultFor[struct{}]{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil
}

func PromptRegressJira(ctx context.Context, ss *mcp.ServerSession, params *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "Mark a Jira as regressed.",
		Messages: []*mcp.PromptMessage{
			{Role: "user", Content: &mcp.TextContent{Text: "Mark JIRA-456 as regressed."}},
		},
	}, nil
}

func ListRegressed(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[struct{}]) (*mcp.CallToolResultFor[struct{}], error) {
	log.Println("[ListRegressed] tool called")
	todos, err := models.LoadTodos(ctx)
	if err != nil {
		log.Printf("[ListRegressed] error loading todos: %v", err)
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error loading todo list: " + err.Error()},
			},
			IsError: true,
		}, nil
	}

	var lines []string
	for _, todo := range todos {
		if todo.Status == "regressed" {
			lines = append(lines, fmt.Sprintf("%s: %s", todo.JiraID, todo.Description))
		}
	}

	if len(lines) == 0 {
		lines = append(lines, "No regressed JIRAs!")
	}

	resultText := "Regressed JIRAs:\n" + util.JoinLines(lines)
	log.Println("[ListRegressed] returning result:", resultText)

	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultText},
		},
	}, nil
}

func PromptListRegressed(ctx context.Context, ss *mcp.ServerSession, params *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "List regressed Jira tasks.",
		Messages: []*mcp.PromptMessage{
			{Role: "user", Content: &mcp.TextContent{Text: "Show regressed Jiras."}},
		},
	}, nil
}
