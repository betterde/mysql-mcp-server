package tools

import (
	"github.com/betterde/mysql-mcp-server/tools/admin"
	"github.com/betterde/mysql-mcp-server/tools/database"
	"github.com/betterde/mysql-mcp-server/tools/query"
	"github.com/betterde/mysql-mcp-server/tools/table"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Register(server *mcp.Server) {
	table.Register(server)
	query.Register(server)
	admin.Register(server)
	database.Register(server)
}
