package tools

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func callTool(t *testing.T, handler func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error), args map[string]any) string {
	t.Helper()
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	res, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if res.IsError {
		t.Fatalf("tool returned error result: %+v", res.Content)
	}
	if len(res.Content) == 0 {
		t.Fatalf("empty content")
	}
	text, ok := res.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("content is not text: %T", res.Content[0])
	}
	return text.Text
}

func TestEchoHandler(t *testing.T) {
	got := callTool(t, echoHandler, map[string]any{"message": "hello"})
	if got != "hello" {
		t.Fatalf("echo: got %q, want %q", got, "hello")
	}
}

func TestAddHandlerIntegers(t *testing.T) {
	got := callTool(t, addHandler, map[string]any{"a": 2.0, "b": 3.0})
	if got != "5" {
		t.Fatalf("add(2,3): got %q, want %q", got, "5")
	}
}

func TestAddHandlerFloats(t *testing.T) {
	got := callTool(t, addHandler, map[string]any{"a": 1.5, "b": 2.25})
	if got != "3.75" {
		t.Fatalf("add(1.5,2.25): got %q, want %q", got, "3.75")
	}
}
