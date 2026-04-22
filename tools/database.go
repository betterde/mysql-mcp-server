package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var database = &mcp.Tool{
	Name:        "Database Tool",
	Description: "Can list databases",
}

type Input struct {
}

type Output struct {
}

func handler(ctx context.Context, req *mcp.CallToolRequest, input Input) (*mcp.CallToolResult, Output, error) {
	return &mcp.CallToolResult{}, Output{}, nil
}

func DatabaseRegister(server *mcp.Server) {
	mcp.AddTool(server, database, handler)
}
