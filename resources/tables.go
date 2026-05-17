package resources

import (
	"context"
	"fmt"
	"net/url"

	"github.com/betterde/mysql-mcp-server/internal/handler"
	"github.com/betterde/mysql-mcp-server/internal/journal"
	"github.com/betterde/mysql-mcp-server/internal/mysql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type tableEntry struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Comment string `json:"comment,omitempty"`
}

type tablesResult struct {
	Database string       `json:"database"`
	Tables   []tableEntry `json:"tables"`
}

func registerTables() {
	s.AddResourceTemplate(&mcp.ResourceTemplate{
		URITemplate: "mysql://{database}/tables",
		Name:        "Tables",
		Title:       "Table List",
		Description: "List all tables in a database.",
		MIMEType:    "application/json",
	}, tablesHandler)
}

func tablesHandler(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	if !checkConn() {
		return nil, mcp.ResourceNotFoundError(req.Params.URI)
	}

	database, err := extractDatabase(req.Params.URI)
	if err != nil {
		return nil, err
	}

	rows, err := mysql.Conn.QueryContext(ctx, `
		SELECT TABLE_NAME, TABLE_TYPE, COALESCE(TABLE_COMMENT, '')
		FROM information_schema.TABLES
		WHERE TABLE_SCHEMA = ?
		ORDER BY TABLE_NAME
	`, database)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := tablesResult{
		Database: database,
		Tables:   make([]tableEntry, 0),
	}
	for rows.Next() {
		var t tableEntry
		if err := rows.Scan(&t.Name, &t.Type, &t.Comment); err != nil {
			return nil, err
		}
		result.Tables = append(result.Tables, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	journal.Logger.Sugar().Debug("Resource accessed", "uri", req.Params.URI)
	return handler.ResourceResult(result)
}

func extractDatabase(rawURI string) (string, error) {
	uri, err := url.Parse(rawURI)
	if err != nil {
		return "", err
	}
	if uri.Host == "" {
		return "", fmt.Errorf("invalid URI %q: missing database name", rawURI)
	}
	return uri.Host, nil
}
