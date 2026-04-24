package admin

import (
	"context"
	"database/sql"
	"errors"

	"github.com/betterde/mysql-mcp-server/intenal/journal"
	"github.com/betterde/mysql-mcp-server/intenal/mysql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var connections = &mcp.Tool{
	Name:        "admin/server/connections",
	Title:       "Show Connections",
	Description: "Show current MySQL client connections, including idle sessions.",
}

type ConnectionsInput struct {
	Limit int `json:"limit,omitempty" jsonschema:"maximum number of connections to return, defaults to 100 and cannot exceed 1000"`
}

type Connection struct {
	ID      uint64 `json:"id"`
	User    string `json:"user"`
	Host    string `json:"host"`
	DB      string `json:"db,omitempty"`
	Command string `json:"command"`
	Time    int    `json:"time"`
	State   string `json:"state,omitempty"`
	Info    string `json:"info,omitempty"`
}

type ConnectionsOutput struct {
	Connections []Connection `json:"connections"`
	Count       int          `json:"count"`
}

func connectionsHandler(ctx context.Context, _ *mcp.CallToolRequest, input ConnectionsInput) (*mcp.CallToolResult, ConnectionsOutput, error) {
	if mysql.Conn == nil {
		return nil, ConnectionsOutput{}, errors.New("mysql connection is not initialized")
	}

	rows, err := mysql.Conn.QueryContext(ctx, `
		SELECT ID, USER, HOST, DB, COMMAND, TIME, STATE, INFO
		FROM information_schema.PROCESSLIST
		ORDER BY ID
	`)
	if err != nil {
		return nil, ConnectionsOutput{}, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			journal.Logger.Error(err.Error())
		}
	}(rows)

	limit := normalizeAdminLimit(input.Limit)
	connections := make([]Connection, 0)
	for rows.Next() {
		if len(connections) >= limit {
			break
		}

		var connection Connection
		var db sql.NullString
		var state sql.NullString
		var info sql.NullString
		if err := rows.Scan(
			&connection.ID,
			&connection.User,
			&connection.Host,
			&db,
			&connection.Command,
			&connection.Time,
			&state,
			&info,
		); err != nil {
			return nil, ConnectionsOutput{}, err
		}

		connection.DB = nullStringValue(db)
		connection.State = nullStringValue(state)
		connection.Info = nullStringValue(info)
		connections = append(connections, connection)
	}

	if err := rows.Err(); err != nil {
		return nil, ConnectionsOutput{}, err
	}

	return &mcp.CallToolResult{}, ConnectionsOutput{
		Connections: connections,
		Count:       len(connections),
	}, nil
}
