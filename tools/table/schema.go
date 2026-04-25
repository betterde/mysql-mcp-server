package table

import (
	"context"
	"errors"

	"github.com/betterde/mysql-mcp-server/internal/mysql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var getTableSchema = &mcp.Tool{
	Name:        "table/schema",
	Title:       "Get Table Schema",
	Description: "Get the CREATE TABLE statement and column metadata for a MySQL table.",
}

type GetTableSchemaInput struct {
	Database string `json:"database" jsonschema:"the database schema name"`
	Table    string `json:"table" jsonschema:"the table name"`
}

type GetTableSchemaOutput struct {
	Database        string   `json:"database"`
	Table           string   `json:"table"`
	CreateStatement string   `json:"create_statement"`
	Columns         []Column `json:"columns"`
}

func getTableSchemaHandler(ctx context.Context, req *mcp.CallToolRequest, input GetTableSchemaInput) (*mcp.CallToolResult, GetTableSchemaOutput, error) {
	if mysql.Conn == nil {
		return nil, GetTableSchemaOutput{}, errors.New("mysql connection is not initialized")
	}
	if input.Database == "" {
		return nil, GetTableSchemaOutput{}, errors.New("database is required")
	}
	if input.Table == "" {
		return nil, GetTableSchemaOutput{}, errors.New("table is required")
	}

	columnsResult, columnsOutput, err := describeTableColumnsHandler(ctx, req, DescribeTableColumnsInput{
		Database: input.Database,
		Table:    input.Table,
	})
	if err != nil {
		return nil, GetTableSchemaOutput{}, err
	}
	if columnsResult != nil && columnsResult.IsError {
		return columnsResult, GetTableSchemaOutput{}, nil
	}

	var tableName string
	var createStatement string
	if err := mysql.Conn.QueryRowContext(ctx, "SHOW CREATE TABLE "+quoteIdentifier(input.Database)+"."+quoteIdentifier(input.Table)).
		Scan(&tableName, &createStatement); err != nil {
		return nil, GetTableSchemaOutput{}, err
	}

	return &mcp.CallToolResult{}, GetTableSchemaOutput{
		Database:        input.Database,
		Table:           tableName,
		CreateStatement: createStatement,
		Columns:         columnsOutput.Columns,
	}, nil
}

func quoteIdentifier(identifier string) string {
	quoted := "`"
	for _, r := range identifier {
		if r == '`' {
			quoted += "``"
			continue
		}
		quoted += string(r)
	}
	return quoted + "`"
}
