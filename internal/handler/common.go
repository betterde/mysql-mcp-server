package handler

import (
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ResourceResult serializes data as indented JSON and returns a ReadResourceResult.
func ResourceResult(data any) (*mcp.ReadResourceResult, error) {
	text, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, err
	}
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{{
			MIMEType: "application/json",
			Text:     string(text),
		}},
	}, nil
}
