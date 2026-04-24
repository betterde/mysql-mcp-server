package admin

import "github.com/modelcontextprotocol/go-sdk/mcp"

func Register(server *mcp.Server) {
	mcp.AddTool(server, health, healthHandler)
	mcp.AddTool(server, process, processHandler)
	mcp.AddTool(server, connections, connectionsHandler)
	mcp.AddTool(server, version, versionHandler)
	mcp.AddTool(server, threads, threadsHandler)
}
