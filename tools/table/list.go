package table

import (
	"context"
	"errors"

	"github.com/betterde/mysql-mcp-server/internal/mysql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var listTables = &mcp.Tool{
	Name:        "tables/list",
	Title:       "List Tables",
	Description: "List tables in a MySQL database schema.",
}

type ListTablesInput struct {
	Database string `json:"database" jsonschema:"the database schema name"`
}

type Table struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Comment string `json:"comment,omitempty"`
}

type ListTablesOutput struct {
	Database string  `json:"database"`
	Tables   []Table `json:"tables"`
}

func listTablesHandler(ctx context.Context, _ *mcp.CallToolRequest, input ListTablesInput) (*mcp.CallToolResult, ListTablesOutput, error) {
	if mysql.Conn == nil {
		return nil, ListTablesOutput{}, errors.New("mysql connection is not initialized")
	}
	if input.Database == "" {
		return nil, ListTablesOutput{}, errors.New("database is required")
	}

	rows, err := mysql.Conn.QueryContext(ctx, `
		SELECT TABLE_NAME, TABLE_TYPE, COALESCE(TABLE_COMMENT, '')
		FROM information_schema.TABLES
		WHERE TABLE_SCHEMA = ?
		ORDER BY TABLE_NAME
	`, input.Database)
	if err != nil {
		return nil, ListTablesOutput{}, err
	}
	defer rows.Close()

	output := ListTablesOutput{
		Database: input.Database,
		Tables:   make([]Table, 0),
	}

	for rows.Next() {
		var table Table
		if err := rows.Scan(&table.Name, &table.Type, &table.Comment); err != nil {
			return nil, ListTablesOutput{}, err
		}
		output.Tables = append(output.Tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, ListTablesOutput{}, err
	}

	return &mcp.CallToolResult{}, output, nil
}
