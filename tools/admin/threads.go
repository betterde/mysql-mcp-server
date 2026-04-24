package admin

import (
	"context"
	"database/sql"
	"errors"

	"github.com/betterde/mysql-mcp-server/intenal/journal"
	"github.com/betterde/mysql-mcp-server/intenal/mysql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var threads = &mcp.Tool{
	Name:        "admin/server/threads",
	Title:       "Show Threads",
	Description: "Show MySQL thread status counters and current thread summary.",
}

type ThreadsInput struct {
}

type ThreadsOutput struct {
	ThreadsConnected int64             `json:"threads_connected"`
	ThreadsRunning   int64             `json:"threads_running"`
	ThreadsCached    int64             `json:"threads_cached"`
	ThreadsCreated   int64             `json:"threads_created"`
	ByCommand        map[string]int64  `json:"by_command"`
	ByUser           map[string]int64  `json:"by_user"`
	Status           map[string]string `json:"status"`
}

func threadsHandler(ctx context.Context, _ *mcp.CallToolRequest, _ ThreadsInput) (*mcp.CallToolResult, ThreadsOutput, error) {
	if mysql.Conn == nil {
		return nil, ThreadsOutput{}, errors.New("mysql connection is not initialized")
	}

	status, err := fetchNameValueRows(ctx, "SHOW GLOBAL STATUS WHERE Variable_name IN ('Threads_connected', 'Threads_running', 'Threads_cached', 'Threads_created')")
	if err != nil {
		return nil, ThreadsOutput{}, err
	}

	byCommand, err := fetchThreadCounts(ctx, "COMMAND")
	if err != nil {
		return nil, ThreadsOutput{}, err
	}

	byUser, err := fetchThreadCounts(ctx, "USER")
	if err != nil {
		return nil, ThreadsOutput{}, err
	}

	return &mcp.CallToolResult{}, ThreadsOutput{
		ThreadsConnected: parseInt64(status["Threads_connected"]),
		ThreadsRunning:   parseInt64(status["Threads_running"]),
		ThreadsCached:    parseInt64(status["Threads_cached"]),
		ThreadsCreated:   parseInt64(status["Threads_created"]),
		ByCommand:        byCommand,
		ByUser:           byUser,
		Status:           status,
	}, nil
}

func fetchThreadCounts(ctx context.Context, column string) (map[string]int64, error) {
	rows, err := mysql.Conn.QueryContext(ctx, "SELECT "+column+", COUNT(*) FROM information_schema.PROCESSLIST GROUP BY "+column+" ORDER BY "+column)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			journal.Logger.Error(err.Error())
		}
	}(rows)

	counts := make(map[string]int64)
	for rows.Next() {
		var name string
		var count int64
		if err := rows.Scan(&name, &count); err != nil {
			return nil, err
		}
		counts[name] = count
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return counts, nil
}
