package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/betterde/mysql-mcp-server/intenal/journal"
	"github.com/betterde/mysql-mcp-server/intenal/mysql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var database = &mcp.Tool{
	Name:        "database/list",
	Title:       "List Databases",
	Description: "List all databases available on the configured MySQL server.",
}

type Input struct {
}

type Output struct {
	Databases []string `json:"databases"`
}

func handler(ctx context.Context, _ *mcp.CallToolRequest, _ Input) (*mcp.CallToolResult, Output, error) {
	if mysql.Conn == nil {
		return nil, Output{}, errors.New("mysql connection is not initialized")
	}

	rows, err := mysql.Conn.QueryContext(ctx, "SHOW DATABASES")
	if err != nil {
		return nil, Output{}, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			journal.Logger.Panic(err.Error())
		}
	}(rows)

	output := Output{
		Databases: make([]string, 0),
	}

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, Output{}, err
		}

		output.Databases = append(output.Databases, name)
	}

	if err := rows.Err(); err != nil {
		return nil, Output{}, err
	}

	return &mcp.CallToolResult{}, output, nil
}
