package tools

import (
	"github.com/betterde/mysql-mcp-server/tools/database"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Register(server *mcp.Server) {
	database.Register(server)
}
