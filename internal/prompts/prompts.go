// Package prompts registers MCP prompt templates.
// Prompts are reusable, parameterised instruction templates returned to the client.
package prompts

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Register adds all prompts to the MCP server.
func Register(s *server.MCPServer) {
	s.AddPrompt(reviewPrompt(), reviewHandler)
	s.AddPrompt(summarisePrompt(), summariseHandler)
}

// --- code-review prompt ---

func reviewPrompt() mcp.Prompt {
	return mcp.NewPrompt("code-review",
		mcp.WithPromptDescription("Generate a thorough code-review prompt for a given snippet"),
		mcp.WithArgument("language",
			mcp.ArgumentDescription("Programming language of the snippet"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("code",
			mcp.ArgumentDescription("The code snippet to review"),
			mcp.RequiredArgument(),
		),
	)
}

func reviewHandler(_ context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	language := req.Params.Arguments["language"]
	code := req.Params.Arguments["code"]

	text := fmt.Sprintf(
		"You are an expert %s developer. Review the following code for correctness, style, "+
			"security issues, and performance. Provide specific, actionable feedback.\n\n```%s\n%s\n```",
		language, language, code,
	)

	return mcp.NewGetPromptResult(
		"Code review prompt",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(text)),
		},
	), nil
}

// --- summarise prompt ---

func summarisePrompt() mcp.Prompt {
	return mcp.NewPrompt("summarise",
		mcp.WithPromptDescription("Summarise a block of text in a given number of sentences"),
		mcp.WithArgument("text",
			mcp.ArgumentDescription("Text to summarise"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("sentences",
			mcp.ArgumentDescription("Target number of sentences (default: 3)"),
		),
	)
}

func summariseHandler(_ context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	text := req.Params.Arguments["text"]
	sentences := req.Params.Arguments["sentences"]
	if sentences == "" {
		sentences = "3"
	}

	content := fmt.Sprintf(
		"Summarise the following text in exactly %s sentence(s):\n\n%s",
		sentences, text,
	)

	return mcp.NewGetPromptResult(
		"Summarise prompt",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(content)),
		},
	), nil
}
