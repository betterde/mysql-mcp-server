package resources

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/betterde/mysql-mcp-server/internal/handler"
	"github.com/betterde/mysql-mcp-server/internal/journal"
	"github.com/betterde/mysql-mcp-server/internal/mysql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type columnEntry struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Nullable  string `json:"nullable"`
	Key       string `json:"key,omitempty"`
	Default   string `json:"default,omitempty"`
	Extra     string `json:"extra,omitempty"`
	Comment   string `json:"comment,omitempty"`
	Collation string `json:"collation,omitempty"`
}

type columnsResult struct {
	Database string        `json:"database"`
	Table    string        `json:"table"`
	Columns  []columnEntry `json:"columns"`
}

func registerColumns() {
	s.AddResourceTemplate(&mcp.ResourceTemplate{
		URITemplate: "mysql://{database}/{table}/columns",
		Name:        "Table Columns",
		Title:       "Column Definitions",
		Description: "List column definitions for a table.",
		MIMEType:    "application/json",
	}, columnsHandler)
}

func columnsHandler(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	if !checkConn() {
		return nil, mcp.ResourceNotFoundError(req.Params.URI)
	}

	database, table, err := extractDatabaseAndTable(req.Params.URI)
	if err != nil {
		return nil, err
	}

	rows, err := mysql.Conn.QueryContext(ctx, `
		SELECT
			COLUMN_NAME,
			COLUMN_TYPE,
			IS_NULLABLE,
			COALESCE(COLUMN_KEY, ''),
			COALESCE(COLUMN_DEFAULT, ''),
			COALESCE(EXTRA, ''),
			COALESCE(COLUMN_COMMENT, ''),
			COALESCE(COLLATION_NAME, '')
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
		ORDER BY ORDINAL_POSITION
	`, database, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := columnsResult{
		Database: database,
		Table:    table,
		Columns:  make([]columnEntry, 0),
	}
	for rows.Next() {
		var c columnEntry
		if err := rows.Scan(&c.Name, &c.Type, &c.Nullable, &c.Key, &c.Default, &c.Extra, &c.Comment, &c.Collation); err != nil {
			return nil, err
		}
		result.Columns = append(result.Columns, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	journal.Logger.Sugar().Debug("Resource accessed", "uri", req.Params.URI)
	return handler.ResourceResult(result)
}

func extractDatabaseAndTable(rawURI string) (string, string, error) {
	uri, err := url.Parse(rawURI)
	if err != nil {
		return "", "", err
	}
	if uri.Host == "" {
		return "", "", fmt.Errorf("invalid URI %q: missing database name", rawURI)
	}
	parts := strings.Split(strings.Trim(uri.Path, "/"), "/")
	if len(parts) < 1 || parts[0] == "" {
		return "", "", fmt.Errorf("invalid URI %q: missing table name", rawURI)
	}
	return uri.Host, parts[0], nil
}
