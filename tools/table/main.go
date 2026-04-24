package table

import "github.com/modelcontextprotocol/go-sdk/mcp"

func Register(server *mcp.Server) {
	mcp.AddTool(server, listTables, listTablesHandler)
	mcp.AddTool(server, getTableSchema, getTableSchemaHandler)
	mcp.AddTool(server, describeTableColumns, describeTableColumnsHandler)
}
