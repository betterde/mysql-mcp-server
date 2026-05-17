package query

import (
	"context"
	"errors"
	"strings"

	"github.com/betterde/mysql-mcp-server/internal/mysql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var executeQuery = &mcp.Tool{
	Name:        "query/execute",
	Title:       "Execute Write Query",
	Description: "Execute an INSERT, UPDATE, DELETE, DDL, or other statement against MySQL. Only available when read_only mode is disabled.",
}

type ExecuteInput struct {
	Query string `json:"query" jsonschema:"the SQL statement to execute"`
}

type ExecuteOutput struct {
	Query        string       `json:"query"`
	RowsAffected int64        `json:"rows_affected"`
	LastInsertID int64        `json:"last_insert_id,omitempty"`
	Result       *QueryResult `json:"result,omitempty"`
}

func executeQueryHandler(ctx context.Context, _ *mcp.CallToolRequest, input ExecuteInput) (*mcp.CallToolResult, ExecuteOutput, error) {
	if mysql.Conn == nil {
		return nil, ExecuteOutput{}, errors.New("mysql connection is not initialized")
	}

	query := strings.TrimSpace(input.Query)
	if query == "" {
		return nil, ExecuteOutput{}, errors.New("query is required")
	}
	if strings.Contains(query, ";") {
		return nil, ExecuteOutput{}, errors.New("multiple statements are not allowed")
	}

	keyword := strings.ToLower(strings.Fields(strings.TrimLeft(query, " \t\r\n("))[0])
	if keyword == "use" {
		return nil, ExecuteOutput{}, errors.New("USE statement is not allowed")
	}

	output := ExecuteOutput{Query: query}

	execResult, err := mysql.Conn.ExecContext(ctx, query)
	if err != nil {
		return nil, ExecuteOutput{}, err
	}

	affected, _ := execResult.RowsAffected()
	output.RowsAffected = affected
	if lastID, err := execResult.LastInsertId(); err == nil {
		output.LastInsertID = lastID
	}

	return &mcp.CallToolResult{}, output, nil
}
