package prompts

import (
	"context"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func callPrompt(t *testing.T, handler func(context.Context, mcp.GetPromptRequest) (*mcp.GetPromptResult, error), args map[string]string) string {
	t.Helper()
	req := mcp.GetPromptRequest{}
	req.Params.Arguments = args
	res, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if len(res.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(res.Messages))
	}
	text, ok := res.Messages[0].Content.(mcp.TextContent)
	if !ok {
		t.Fatalf("content is not text: %T", res.Messages[0].Content)
	}
	return text.Text
}

func TestReviewHandlerInterpolates(t *testing.T) {
	got := callPrompt(t, reviewHandler, map[string]string{
		"language": "go",
		"code":     "fmt.Println(\"hi\")",
	})
	if !strings.Contains(got, "expert go developer") {
		t.Errorf("missing language in prompt: %q", got)
	}
	if !strings.Contains(got, `fmt.Println("hi")`) {
		t.Errorf("missing code in prompt: %q", got)
	}
}

func TestSummariseHandlerDefaultsSentences(t *testing.T) {
	got := callPrompt(t, summariseHandler, map[string]string{"text": "a long passage"})
	if !strings.Contains(got, "exactly 3 sentence") {
		t.Errorf("expected default of 3 sentences, got %q", got)
	}
	if !strings.Contains(got, "a long passage") {
		t.Errorf("missing text in prompt: %q", got)
	}
}

func TestSummariseHandlerExplicitSentences(t *testing.T) {
	got := callPrompt(t, summariseHandler, map[string]string{
		"text":      "body",
		"sentences": "5",
	})
	if !strings.Contains(got, "exactly 5 sentence") {
		t.Errorf("expected 5 sentences, got %q", got)
	}
}
