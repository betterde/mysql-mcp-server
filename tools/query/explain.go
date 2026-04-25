package query

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/betterde/mysql-mcp-server/internal/mysql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var explainQuery = &mcp.Tool{
	Name:        "query/explain",
	Title:       "Explain Query",
	Description: "Run EXPLAIN for a read-only MySQL query and return the execution plan.",
}

type ExplainInput struct {
	Query  string `json:"query" jsonschema:"the read-only SQL query to explain"`
	Format string `json:"format,omitempty" jsonschema:"optional EXPLAIN format, one of TRADITIONAL, TREE, or JSON"`
}

type ExplainOutput struct {
	Query  string      `json:"query"`
	Format string      `json:"format,omitempty"`
	Plan   QueryResult `json:"plan"`
}

func explainQueryHandler(ctx context.Context, _ *mcp.CallToolRequest, input ExplainInput) (*mcp.CallToolResult, ExplainOutput, error) {
	if mysql.Conn == nil {
		return nil, ExplainOutput{}, errors.New("mysql connection is not initialized")
	}

	query, err := normalizeReadOnlyQuery(input.Query)
	if err != nil {
		return nil, ExplainOutput{}, err
	}

	format, err := normalizeExplainFormat(input.Format)
	if err != nil {
		return nil, ExplainOutput{}, err
	}

	explainSQL := "EXPLAIN " + query
	if format != "" {
		explainSQL = fmt.Sprintf("EXPLAIN FORMAT=%s %s", format, query)
	}

	rows, err := mysql.Conn.QueryContext(ctx, explainSQL)
	if err != nil {
		return nil, ExplainOutput{}, err
	}
	defer rows.Close()

	plan, err := scanRows(rows, maxSelectLimit)
	if err != nil {
		return nil, ExplainOutput{}, err
	}

	return &mcp.CallToolResult{}, ExplainOutput{
		Query:  query,
		Format: format,
		Plan:   plan,
	}, nil
}

func normalizeExplainFormat(format string) (string, error) {
	format = strings.ToUpper(strings.TrimSpace(format))
	switch format {
	case "":
		return "", nil
	case "TRADITIONAL", "TREE", "JSON":
		return format, nil
	default:
		return "", fmt.Errorf("unsupported explain format %q", format)
	}
}
