package admin

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/betterde/mysql-mcp-server/intenal/journal"
	"github.com/betterde/mysql-mcp-server/intenal/mysql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var health = &mcp.Tool{
	Name:        "admin/server/health",
	Title:       "Server Health",
	Description: "Show current MySQL server health status and key runtime metrics.",
}

type HealthInput struct {
}

type HealthOutput struct {
	Healthy            bool              `json:"healthy"`
	Version            string            `json:"version,omitempty"`
	VersionComment     string            `json:"version_comment,omitempty"`
	ReadOnly           bool              `json:"read_only"`
	UptimeSeconds      int64             `json:"uptime_seconds"`
	ThreadsConnected   int64             `json:"threads_connected"`
	ThreadsRunning     int64             `json:"threads_running"`
	Connections        int64             `json:"connections"`
	MaxConnections     int64             `json:"max_connections"`
	MaxUsedConnections int64             `json:"max_used_connections"`
	Questions          int64             `json:"questions"`
	SlowQueries        int64             `json:"slow_queries"`
	AbortedConnects    int64             `json:"aborted_connects"`
	Status             map[string]string `json:"status"`
	Variables          map[string]string `json:"variables"`
}

func healthHandler(ctx context.Context, _ *mcp.CallToolRequest, _ HealthInput) (*mcp.CallToolResult, HealthOutput, error) {
	if mysql.Conn == nil {
		return nil, HealthOutput{}, errors.New("mysql connection is not initialized")
	}
	if err := mysql.Conn.PingContext(ctx); err != nil {
		return nil, HealthOutput{}, err
	}

	status, err := fetchNameValueRows(ctx, "SHOW GLOBAL STATUS WHERE Variable_name IN ('Uptime', 'Threads_connected', 'Threads_running', 'Connections', 'Max_used_connections', 'Questions', 'Slow_queries', 'Aborted_connects')")
	if err != nil {
		return nil, HealthOutput{}, err
	}

	variables, err := fetchNameValueRows(ctx, "SHOW VARIABLES WHERE Variable_name IN ('version', 'version_comment', 'read_only', 'super_read_only', 'max_connections')")
	if err != nil {
		return nil, HealthOutput{}, err
	}

	readOnly := variables["read_only"] == "ON" || variables["super_read_only"] == "ON"
	output := HealthOutput{
		Healthy:            true,
		Version:            variables["version"],
		VersionComment:     variables["version_comment"],
		ReadOnly:           readOnly,
		UptimeSeconds:      parseInt64(status["Uptime"]),
		ThreadsConnected:   parseInt64(status["Threads_connected"]),
		ThreadsRunning:     parseInt64(status["Threads_running"]),
		Connections:        parseInt64(status["Connections"]),
		MaxConnections:     parseInt64(variables["max_connections"]),
		MaxUsedConnections: parseInt64(status["Max_used_connections"]),
		Questions:          parseInt64(status["Questions"]),
		SlowQueries:        parseInt64(status["Slow_queries"]),
		AbortedConnects:    parseInt64(status["Aborted_connects"]),
		Status:             status,
		Variables:          variables,
	}

	return &mcp.CallToolResult{}, output, nil
}

func fetchNameValueRows(ctx context.Context, query string) (map[string]string, error) {
	rows, err := mysql.Conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			journal.Logger.Error(err.Error())
		}
	}(rows)

	values := make(map[string]string)
	for rows.Next() {
		var name string
		var value string
		if err := rows.Scan(&name, &value); err != nil {
			return nil, err
		}
		values[name] = value
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return values, nil
}

func parseInt64(value string) int64 {
	number, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return number
}
