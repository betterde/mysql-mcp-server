package query

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/betterde/mysql-mcp-server/internal/mysql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	defaultSelectLimit = 100
	maxSelectLimit     = 1000
)

var selectQuery = &mcp.Tool{
	Name:        "query/select",
	Title:       "Execute Select Query",
	Description: "Execute a read-only SELECT, SHOW, DESCRIBE, or WITH query against MySQL and return rows.",
}

type SelectInput struct {
	Query string `json:"query" jsonschema:"the read-only SQL query to execute"`
	Limit int    `json:"limit,omitempty" jsonschema:"maximum number of rows to return, defaults to 100 and cannot exceed 1000"`
}

type QueryResult struct {
	Columns []string         `json:"columns"`
	Rows    []map[string]any `json:"rows"`
	Count   int              `json:"count"`
}

type SelectOutput struct {
	Query  string      `json:"query"`
	Limit  int         `json:"limit"`
	Result QueryResult `json:"result"`
}

func selectQueryHandler(ctx context.Context, _ *mcp.CallToolRequest, input SelectInput) (*mcp.CallToolResult, SelectOutput, error) {
	if mysql.Conn == nil {
		return nil, SelectOutput{}, errors.New("mysql connection is not initialized")
	}
	query, err := normalizeReadOnlyQuery(input.Query)
	if err != nil {
		return nil, SelectOutput{}, err
	}

	limit := normalizeLimit(input.Limit)
	rows, err := mysql.Conn.QueryContext(ctx, query)
	if err != nil {
		return nil, SelectOutput{}, err
	}
	defer rows.Close()

	result, err := scanRows(rows, limit)
	if err != nil {
		return nil, SelectOutput{}, err
	}

	return &mcp.CallToolResult{}, SelectOutput{
		Query:  query,
		Limit:  limit,
		Result: result,
	}, nil
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return defaultSelectLimit
	}
	if limit > maxSelectLimit {
		return maxSelectLimit
	}
	return limit
}

func normalizeReadOnlyQuery(query string) (string, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return "", errors.New("query is required")
	}
	if strings.Contains(query, ";") {
		return "", errors.New("multiple statements are not allowed")
	}

	keyword := firstKeyword(query)
	switch keyword {
	case "select", "show", "describe", "desc", "with":
		return query, nil
	default:
		return "", fmt.Errorf("only read-only queries are allowed, got %q", keyword)
	}
}

func firstKeyword(query string) string {
	fields := strings.Fields(strings.TrimLeft(query, " \t\r\n("))
	if len(fields) == 0 {
		return ""
	}
	return strings.ToLower(fields[0])
}

func scanRows(rows *sql.Rows, limit int) (QueryResult, error) {
	columns, err := rows.Columns()
	if err != nil {
		return QueryResult{}, err
	}

	result := QueryResult{
		Columns: columns,
		Rows:    make([]map[string]any, 0),
	}

	for rows.Next() {
		if result.Count >= limit {
			break
		}

		values := make([]any, len(columns))
		pointers := make([]any, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return QueryResult{}, err
		}

		row := make(map[string]any, len(columns))
		for i, column := range columns {
			row[column] = normalizeValue(values[i])
		}
		result.Rows = append(result.Rows, row)
		result.Count++
	}

	if err := rows.Err(); err != nil {
		return QueryResult{}, err
	}

	return result, nil
}

func normalizeValue(value any) any {
	switch v := value.(type) {
	case nil:
		return nil
	case []byte:
		return string(v)
	case time.Time:
		return v.Format(time.RFC3339Nano)
	case int64:
		return v
	case uint64:
		return v
	case float64:
		return v
	case bool:
		return v
	default:
		return fmt.Sprint(v)
	}
}
