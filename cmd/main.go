package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"todo-app/internal/models"
	"todo-app/internal/tools"
	"todo-app/internal/util"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var httpAddr = flag.String("http", "", "If set, use streamable HTTP at this address instead of stdio")

func main() {
	flag.Parse()

	util.LogStartupInfo(models.TodoFilePath)

	server := mcp.NewServer(&mcp.Implementation{Name: "jira-todo"}, nil)

	// Register tools and prompts
	mcp.AddTool(server, &mcp.Tool{Name: "add_jira", Description: "Add a Jira todo to your list"}, tools.AddJira)
	mcp.AddTool(server, &mcp.Tool{Name: "complete_jira", Description: "Mark a Jira as completed"}, tools.CompleteJira)
	mcp.AddTool(server, &mcp.Tool{Name: "regress_jira", Description: "Mark a Jira as regressed"}, tools.RegressJira)
	mcp.AddTool(server, &mcp.Tool{Name: "list_todos", Description: "List pending todos"}, tools.ListTodos)
	mcp.AddTool(server, &mcp.Tool{Name: "list_completed", Description: "List completed Jiras"}, tools.ListCompleted)
	mcp.AddTool(server, &mcp.Tool{Name: "list_regressed", Description: "List regressed Jiras"}, tools.ListRegressed)

	server.AddPrompt(&mcp.Prompt{Name: "add_jira"}, tools.PromptAddJira)
	server.AddPrompt(&mcp.Prompt{Name: "complete_jira"}, tools.PromptCompleteJira)
	server.AddPrompt(&mcp.Prompt{Name: "regress_jira"}, tools.PromptRegressJira)
	server.AddPrompt(&mcp.Prompt{Name: "list_todos"}, tools.PromptListTodos)
	server.AddPrompt(&mcp.Prompt{Name: "list_completed"}, tools.PromptListCompleted)
	server.AddPrompt(&mcp.Prompt{Name: "list_regressed"}, tools.PromptListRegressed)

	if *httpAddr != "" {
		handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
			return server
		}, nil)
		log.Printf("MCP HTTP handler listening at %s", *httpAddr)
		if err := http.ListenAndServe(*httpAddr, handler); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	} else {
		transport := mcp.NewLoggingTransport(mcp.NewStdioTransport(), os.Stderr)
		log.Println("Starting Jira TODO MCP server (stdio mode)...")
		if err := server.Run(context.Background(), transport); err != nil {
			log.Fatalf("MCP server error: %v", err)
		}
	}
}
