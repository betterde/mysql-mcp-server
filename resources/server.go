package resources

import (
	"context"
	"database/sql"

	"github.com/betterde/mysql-mcp-server/internal/handler"
	"github.com/betterde/mysql-mcp-server/internal/journal"
	"github.com/betterde/mysql-mcp-server/internal/mysql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type serverEntry struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func registerServer() {
	s.AddResource(&mcp.Resource{
		URI:         "mysql://server/version",
		Name:        "Server Version",
		Title:       "MySQL Server Version",
		Description: "Show the MySQL server version information.",
		MIMEType:    "application/json",
	}, versionHandler)

	s.AddResource(&mcp.Resource{
		URI:         "mysql://server/variables",
		Name:        "Server Variables",
		Title:       "Global System Variables",
		Description: "List all global system variables of the MySQL server.",
		MIMEType:    "application/json",
	}, variablesHandler)

	s.AddResource(&mcp.Resource{
		URI:         "mysql://server/status",
		Name:        "Server Status",
		Title:       "Global Status Variables",
		Description: "List all global status variables of the MySQL server.",
		MIMEType:    "application/json",
	}, statusHandler)
}

func versionHandler(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	if !checkConn() {
		return nil, mcp.ResourceNotFoundError(req.Params.URI)
	}

	rows, err := mysql.Conn.QueryContext(ctx, "SELECT VERSION()")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var version string
	if rows.Next() {
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
	}

	journal.Logger.Sugar().Debug("Resource accessed", "uri", req.Params.URI)
	return handler.ResourceResult(map[string]string{"version": version})
}

func variablesHandler(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	if !checkConn() {
		return nil, mcp.ResourceNotFoundError(req.Params.URI)
	}

	entries, err := fetchServerEntries(ctx, "SHOW GLOBAL VARIABLES")
	if err != nil {
		return nil, err
	}

	journal.Logger.Sugar().Debug("Resource accessed", "uri", req.Params.URI)
	return handler.ResourceResult(entries)
}

func statusHandler(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	if !checkConn() {
		return nil, mcp.ResourceNotFoundError(req.Params.URI)
	}

	entries, err := fetchServerEntries(ctx, "SHOW GLOBAL STATUS")
	if err != nil {
		return nil, err
	}

	journal.Logger.Sugar().Debug("Resource accessed", "uri", req.Params.URI)
	return handler.ResourceResult(entries)
}

func fetchServerEntries(ctx context.Context, query string) ([]serverEntry, error) {
	rows, err := mysql.Conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if cerr := rows.Close(); cerr != nil {
			journal.Logger.Error(cerr.Error())
		}
	}(rows)

	entries := make([]serverEntry, 0)
	for rows.Next() {
		var name, value string
		if err := rows.Scan(&name, &value); err != nil {
			return nil, err
		}
		entries = append(entries, serverEntry{Name: name, Value: value})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
