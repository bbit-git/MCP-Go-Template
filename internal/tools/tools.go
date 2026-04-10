// Package tools registers MCP tool handlers.
// Add your tool implementations here following the example below.
package tools

import (
	"context"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Register adds all tools to the MCP server.
func Register(s *server.MCPServer) {
	s.AddTool(echoTool(), echoHandler)
	s.AddTool(addTool(), addHandler)
}

// --- echo tool ---

func echoTool() mcp.Tool {
	return mcp.NewTool("echo",
		mcp.WithDescription("Echo a message back to the caller"),
		mcp.WithString("message",
			mcp.Required(),
			mcp.Description("The message to echo"),
		),
	)
}

func echoHandler(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	msg, err := req.RequireString("message")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(msg), nil
}

// --- add tool ---

func addTool() mcp.Tool {
	return mcp.NewTool("add",
		mcp.WithDescription("Add two numbers and return the sum"),
		mcp.WithNumber("a",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("b",
			mcp.Required(),
			mcp.Description("Second number"),
		),
	)
}

func addHandler(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a, err := req.RequireFloat("a")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	b, err := req.RequireFloat("b")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(strconv.FormatFloat(a+b, 'f', -1, 64)), nil
}
