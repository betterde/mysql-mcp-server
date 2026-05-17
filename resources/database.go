package resources

import (
	"context"

	"github.com/betterde/mysql-mcp-server/internal/handler"
	"github.com/betterde/mysql-mcp-server/internal/journal"
	"github.com/betterde/mysql-mcp-server/internal/mysql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type databaseEntry struct {
	Name string `json:"name"`
}

type databasesResult struct {
	Databases []databaseEntry `json:"databases"`
}

func registerDatabases() {
	s.AddResource(&mcp.Resource{
		URI:         "mysql://databases",
		Name:        "Databases",
		Title:       "Database List",
		Description: "List all databases on the MySQL server.",
		MIMEType:    "application/json",
	}, databasesHandler)
}

func databasesHandler(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	if !checkConn() {
		return nil, mcp.ResourceNotFoundError(req.Params.URI)
	}

	rows, err := mysql.Conn.QueryContext(ctx, "SHOW DATABASES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := databasesResult{Databases: make([]databaseEntry, 0)}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		result.Databases = append(result.Databases, databaseEntry{Name: name})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	journal.Logger.Sugar().Debug("Resource accessed", "uri", req.Params.URI)
	return handler.ResourceResult(result)
}
