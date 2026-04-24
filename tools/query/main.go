package query

import "github.com/modelcontextprotocol/go-sdk/mcp"

func Register(server *mcp.Server) {
	mcp.AddTool(server, selectQuery, selectQueryHandler)
	mcp.AddTool(server, explainQuery, explainQueryHandler)
}
