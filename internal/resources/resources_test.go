package resources

import (
	"context"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestItemByIDHandlerExtractsID(t *testing.T) {
	req := mcp.ReadResourceRequest{}
	req.Params.URI = "item://42"
	req.Params.Arguments = map[string]any{"id": "42"}

	contents, err := itemByIDHandler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if len(contents) != 1 {
		t.Fatalf("expected 1 content, got %d", len(contents))
	}
	text, ok := contents[0].(mcp.TextResourceContents)
	if !ok {
		t.Fatalf("content is not text: %T", contents[0])
	}
	if !strings.Contains(text.Text, `"id": "42"`) {
		t.Fatalf("expected id 42 in body, got %q", text.Text)
	}
	if strings.Contains(text.Text, "item://") {
		t.Fatalf("scheme leaked into body: %q", text.Text)
	}
}

func TestItemByIDHandlerMissingID(t *testing.T) {
	req := mcp.ReadResourceRequest{}
	req.Params.URI = "item://"
	if _, err := itemByIDHandler(context.Background(), req); err == nil {
		t.Fatal("expected error for missing id, got nil")
	}
}
