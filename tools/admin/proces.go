package admin

import (
	"context"
	"database/sql"
	"errors"

	"github.com/betterde/mysql-mcp-server/internal/journal"
	"github.com/betterde/mysql-mcp-server/internal/mysql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var process = &mcp.Tool{
	Name:        "admin/server/queries",
	Title:       "Show Running Queries",
	Description: "Show current non-idle MySQL sessions and running queries.",
}

type ProcessInput struct {
	Limit int `json:"limit,omitempty" jsonschema:"maximum number of running queries to return, defaults to 100 and cannot exceed 1000"`
}

type Process struct {
	ID      uint64 `json:"id"`
	User    string `json:"user"`
	Host    string `json:"host"`
	DB      string `json:"db,omitempty"`
	Command string `json:"command"`
	Time    int    `json:"time"`
	State   string `json:"state,omitempty"`
	Info    string `json:"info,omitempty"`
}

type ProcessOutput struct {
	Processes []Process `json:"processes"`
	Count     int       `json:"count"`
}

func processHandler(ctx context.Context, _ *mcp.CallToolRequest, input ProcessInput) (*mcp.CallToolResult, ProcessOutput, error) {
	if mysql.Conn == nil {
		return nil, ProcessOutput{}, errors.New("mysql connection is not initialized")
	}

	rows, err := mysql.Conn.QueryContext(ctx, `
		SELECT ID, USER, HOST, DB, COMMAND, TIME, STATE, INFO
		FROM information_schema.PROCESSLIST
		WHERE COMMAND <> 'Sleep'
		ORDER BY TIME DESC, ID
	`)
	if err != nil {
		return nil, ProcessOutput{}, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			journal.Logger.Error(err.Error())
		}
	}(rows)

	limit := normalizeAdminLimit(input.Limit)
	processes := make([]Process, 0)
	for rows.Next() {
		if len(processes) >= limit {
			break
		}

		var process Process
		var db sql.NullString
		var state sql.NullString
		var info sql.NullString
		if err := rows.Scan(
			&process.ID,
			&process.User,
			&process.Host,
			&db,
			&process.Command,
			&process.Time,
			&state,
			&info,
		); err != nil {
			return nil, ProcessOutput{}, err
		}

		process.DB = nullStringValue(db)
		process.State = nullStringValue(state)
		process.Info = nullStringValue(info)
		processes = append(processes, process)
	}

	if err := rows.Err(); err != nil {
		return nil, ProcessOutput{}, err
	}

	return &mcp.CallToolResult{}, ProcessOutput{
		Processes: processes,
		Count:     len(processes),
	}, nil
}
