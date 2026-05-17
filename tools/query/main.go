package query

import (
	"github.com/betterde/mysql-mcp-server/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Register(server *mcp.Server) {
	mcp.AddTool(server, selectQuery, selectQueryHandler)
	mcp.AddTool(server, explainQuery, explainQueryHandler)
	if !config.Conf.ReadOnly {
		mcp.AddTool(server, executeQuery, executeQueryHandler)
	}
}
