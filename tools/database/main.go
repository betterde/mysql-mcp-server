package database

import "github.com/modelcontextprotocol/go-sdk/mcp"

func Register(server *mcp.Server) {
	mcp.AddTool(server, database, handler)
}
