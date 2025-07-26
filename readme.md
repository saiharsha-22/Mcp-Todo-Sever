# Jira TODO MCP Server

A [Model Con Protocol (MCP)](https://modelconprotocol.org/) server implemented in Go for managing Jira TODOs.  
This server exposes a set of MCP tools to add, complete, regress, and list Jira tasks, supporting integration with MCP clients like Claude Desktop.

---

## Features

- Add Jira TODO items with severity, due date, and description.
- Mark Jira TODOs as completed or regressed.
- List pending, completed, and regressed Jira TODOs.
- Data persistence using a JSON file stored on the local filesystem.
- Supports both stdio mode (default) and HTTP streaming mode for MCP protocol.

---
## Project Structure
<pre> 
todo-app/
│
├── cmd/
│   └── jira-todo/           # main executable package
│       └── main.go
│
├── internal/
│   ├── models/
│   │   └── todo.go          # TodoItem struct, todo file path, load/save logic
│   │
│   ├── tools/
│   │   ├── add_jira.go       # AddJira tool implementation + prompt
│   │   ├── complete_jira.go  # CompleteJira tool + prompt
│   │   ├── regress_jira.go   # RegressJira tool + prompt
│   │   ├── list_todos.go     # ListTodos tool + prompt
│   │
│   └── util/
│       └── logger.go         # logStartupInfo and utility functions
│
├── go.mod
└── go.sum
</pre>

---

## Prerequisites

- Go (version 1.20+ recommended) installed on your system: [golang.org](https://golang.org/dl/)
- MCP client (e.g., **Claude Desktop**) configured to launch this server over stdio.

---
## Installation
1. Clone this repository:
git clone https://github.com/yourusername/jira-todo-mcp.git
cd jira-todo-mcp/cmd/jira-todo
2. Build the server binary:
go build -o jira-todo
3. Make sure the binary has execution permission:
chmod +x jira-todo
---
## Running the Server
### Default (stdio mode) - recommended for Claude Desktop
Simply run:
./jira-todo
This starts the MCP server communicating over standard input/output, suitable for MCP clients that launch it as a subprocess (like Claude Desktop).

### HTTP streaming mode
You can also run the server with HTTP streaming mode for clients that support it:
./jira-todo -http=:12345

This opens a network port (e.g., 12345) for MCP JSON-RPC streaming over HTTP.
---
## Configuration & Data Storage
- The server stores todos in JSON format in a per-user location:
  - By default, this file is `${HOME}/todo.json` on Unix-based systems including macOS.
  - If the home directory cannot be found, it falls back to `/tmp/todo.json`.
- Logs startup info including current working directory and todo file path.
---
## Integrating with Claude Desktop
1. Configure Claude Desktop to launch this MCP server by editing (or creating) `claude_desktop_config.json`:
<pre>
{
  "mcpServers": {
    "jira-todo-mcp-server": {
      "command": "PATH_OF_PROJECT_BUILD",
      "args": [],
      "env": {}
    }
  }
}
</pre>
3. Restart Claude Desktop after configuration.
4. Open the Tools sidebar; the Jira TODO tools (`add_jira`, `complete_jira`, `list_todos`, etc.) should be visible and ready for use.
---

## Toolset Overview

| Tool Name       | Description                       |
|-----------------|---------------------------------|
| add_jira        | Add a Jira TODO item             |
| complete_jira   | Mark a Jira TODO as completed    |
| regress_jira    | Mark a Jira TODO as regressed    |
| list_todos      | List pending Jira TODOs          |
| list_completed  | List completed Jira TODOs        |
| list_regressed  | List regressed Jira TODOs        |

Each tool has an associated prompt registered to assist the MCP client in presenting usage examples.
---
## Logging & Debugging
- Server logs are output to `stderr` with detailed startup info and tool invocation logs.
- On startup, the server logs the current working directory and the path to the todo JSON file.
- Use logs to troubleshoot file permission issues or MCP protocol communication.
---

## Contribution & Development

- Contributions and issues are welcome! Please open GitHub issues or pull requests.
- Follow Go best practices for code organization (see project structure).
- Remember to call `flag.Parse()` before registering or running the server.
---
## Acknowledgements
- Uses the [Model Con Protocol Go SDK](https://github.com/modelconprotocol/go-sdk).
- Inspired by ideas from Claude and MCP-based tool integrations.
---
*Happy coding & MCP tooling!*
