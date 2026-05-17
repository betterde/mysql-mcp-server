package tools

import (
	"github.com/betterde/mysql-mcp-server/tools/admin"
	"github.com/betterde/mysql-mcp-server/tools/query"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Register(server *mcp.Server) {
	query.Register(server)
	admin.Register(server)
}
