package admin

import (
	"context"
	"errors"

	"github.com/betterde/mysql-mcp-server/intenal/mysql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var version = &mcp.Tool{
	Name:        "admin/server/version",
	Title:       "Server Version",
	Description: "Show MySQL server version and build information.",
}

type VersionInput struct {
}

type VersionOutput struct {
	Version               string            `json:"version"`
	VersionComment        string            `json:"version_comment,omitempty"`
	VersionCompileMachine string            `json:"version_compile_machine,omitempty"`
	VersionCompileOS      string            `json:"version_compile_os,omitempty"`
	ProtocolVersion       string            `json:"protocol_version,omitempty"`
	Variables             map[string]string `json:"variables"`
}

func versionHandler(ctx context.Context, _ *mcp.CallToolRequest, _ VersionInput) (*mcp.CallToolResult, VersionOutput, error) {
	if mysql.Conn == nil {
		return nil, VersionOutput{}, errors.New("mysql connection is not initialized")
	}

	variables, err := fetchNameValueRows(ctx, "SHOW VARIABLES WHERE Variable_name IN ('version', 'version_comment', 'version_compile_machine', 'version_compile_os', 'protocol_version')")
	if err != nil {
		return nil, VersionOutput{}, err
	}

	return &mcp.CallToolResult{}, VersionOutput{
		Version:               variables["version"],
		VersionComment:        variables["version_comment"],
		VersionCompileMachine: variables["version_compile_machine"],
		VersionCompileOS:      variables["version_compile_os"],
		ProtocolVersion:       variables["protocol_version"],
		Variables:             variables,
	}, nil
}
