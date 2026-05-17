package resources

import (
	"context"

	"github.com/betterde/mysql-mcp-server/internal/handler"
	"github.com/betterde/mysql-mcp-server/internal/journal"
	"github.com/betterde/mysql-mcp-server/internal/mysql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ddlResult struct {
	Database string `json:"database"`
	Table    string `json:"table"`
	DDL      string `json:"ddl"`
}

func registerDDL() {
	s.AddResourceTemplate(&mcp.ResourceTemplate{
		URITemplate: "mysql://{database}/{table}/ddl",
		Name:        "Table DDL",
		Title:       "Table Schema DDL",
		Description: "Get the CREATE TABLE DDL for a table.",
		MIMEType:    "application/json",
	}, ddlHandler)
}

func ddlHandler(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	if !checkConn() {
		return nil, mcp.ResourceNotFoundError(req.Params.URI)
	}

	database, table, err := extractDatabaseAndTable(req.Params.URI)
	if err != nil {
		return nil, err
	}

	rows, err := mysql.Conn.QueryContext(ctx, "SHOW CREATE TABLE `"+database+"`.`"+table+"`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result ddlResult
	if rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName, &result.DDL); err != nil {
			return nil, err
		}
		result.Database = database
		result.Table = table
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	journal.Logger.Sugar().Debug("Resource accessed", "uri", req.Params.URI)
	return handler.ResourceResult(result)
}
