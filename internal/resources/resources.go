// Package resources registers MCP resource handlers.
// Resources expose data that can be read into LLM context.
package resources

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Register adds all resources to the MCP server.
func Register(s *server.MCPServer) {
	// Static resource: always the same content
	s.AddResource(
		mcp.NewResource(
			"info://server",
			"Server Info",
			mcp.WithResourceDescription("Basic information about this MCP server"),
			mcp.WithMIMEType("text/plain"),
		),
		serverInfoHandler,
	)

	// Resource template: URI contains a variable (like a REST path param)
	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"item://{id}",
			"Item by ID",
			mcp.WithTemplateDescription("Fetch an item by its ID"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		itemByIDHandler,
	)
}

func serverInfoHandler(_ context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "info://server",
			MIMEType: "text/plain",
			Text:     "mcp-go-template v0.1.0 — replace me with your server description.",
		},
	}, nil
}

func itemByIDHandler(_ context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// mark3labs matches the URI against the template and places each
	// named variable (e.g. {id}) in req.Params.Arguments.
	id, _ := req.Params.Arguments["id"].(string)
	if id == "" {
		return nil, fmt.Errorf("missing id in URI %q", req.Params.URI)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      req.Params.URI,
			MIMEType: "application/json",
			Text:     fmt.Sprintf(`{"id": %q, "name": "Example item %s"}`, id, id),
		},
	}, nil
}
