package resources

import (
	"github.com/betterde/mysql-mcp-server/internal/mysql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	s *mcp.Server
)

func Register(server *mcp.Server) {
	s = server
	registerDatabases()
	registerTables()
	registerColumns()
	registerDDL()
	registerServer()
}

func checkConn() bool {
	return mysql.Conn != nil
}
