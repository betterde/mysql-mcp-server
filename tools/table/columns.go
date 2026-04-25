package table

import (
	"context"
	"database/sql"
	"errors"

	"github.com/betterde/mysql-mcp-server/internal/journal"
	"github.com/betterde/mysql-mcp-server/internal/mysql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var describeTableColumns = &mcp.Tool{
	Name:        "table/columns",
	Title:       "Describe Table Columns",
	Description: "List column metadata for a MySQL table.",
}

type DescribeTableColumnsInput struct {
	Database string `json:"database" jsonschema:"the database schema name"`
	Table    string `json:"table" jsonschema:"the table name"`
}

type Column struct {
	Name             string `json:"name"`
	Ordinal          int    `json:"ordinal"`
	Type             string `json:"type"`
	DataType         string `json:"data_type"`
	Nullable         bool   `json:"nullable"`
	Key              string `json:"key,omitempty"`
	Default          string `json:"default,omitempty"`
	Extra            string `json:"extra,omitempty"`
	Comment          string `json:"comment,omitempty"`
	CharacterSet     string `json:"character_set,omitempty"`
	Collation        string `json:"collation,omitempty"`
	NumericPrecision string `json:"numeric_precision,omitempty"`
	NumericScale     string `json:"numeric_scale,omitempty"`
}

type DescribeTableColumnsOutput struct {
	Database string   `json:"database"`
	Table    string   `json:"table"`
	Columns  []Column `json:"columns"`
}

func describeTableColumnsHandler(ctx context.Context, _ *mcp.CallToolRequest, input DescribeTableColumnsInput) (*mcp.CallToolResult, DescribeTableColumnsOutput, error) {
	if mysql.Conn == nil {
		return nil, DescribeTableColumnsOutput{}, errors.New("mysql connection is not initialized")
	}
	if input.Database == "" {
		return nil, DescribeTableColumnsOutput{}, errors.New("database is required")
	}
	if input.Table == "" {
		return nil, DescribeTableColumnsOutput{}, errors.New("table is required")
	}

	rows, err := mysql.Conn.QueryContext(ctx, `
		SELECT
			COLUMN_NAME,
			ORDINAL_POSITION,
			COLUMN_TYPE,
			DATA_TYPE,
			IS_NULLABLE,
			COALESCE(COLUMN_KEY, ''),
			COALESCE(COLUMN_DEFAULT, ''),
			COALESCE(EXTRA, ''),
			COALESCE(COLUMN_COMMENT, ''),
			COALESCE(CHARACTER_SET_NAME, ''),
			COALESCE(COLLATION_NAME, ''),
			COALESCE(CAST(NUMERIC_PRECISION AS CHAR), ''),
			COALESCE(CAST(NUMERIC_SCALE AS CHAR), '')
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
		ORDER BY ORDINAL_POSITION
	`, input.Database, input.Table)
	if err != nil {
		return nil, DescribeTableColumnsOutput{}, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			journal.Logger.Error(err.Error())
		}
	}(rows)

	output := DescribeTableColumnsOutput{
		Database: input.Database,
		Table:    input.Table,
		Columns:  make([]Column, 0),
	}

	for rows.Next() {
		var column Column
		var nullable string
		if err := rows.Scan(
			&column.Name,
			&column.Ordinal,
			&column.Type,
			&column.DataType,
			&nullable,
			&column.Key,
			&column.Default,
			&column.Extra,
			&column.Comment,
			&column.CharacterSet,
			&column.Collation,
			&column.NumericPrecision,
			&column.NumericScale,
		); err != nil {
			return nil, DescribeTableColumnsOutput{}, err
		}
		column.Nullable = nullable == "YES"
		output.Columns = append(output.Columns, column)
	}

	if err := rows.Err(); err != nil {
		return nil, DescribeTableColumnsOutput{}, err
	}

	return &mcp.CallToolResult{}, output, nil
}
