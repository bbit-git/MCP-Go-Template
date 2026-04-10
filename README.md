# mcp-go-template

A minimal [Model Context Protocol](https://modelcontextprotocol.io/) server
template written in Go, built on [`mark3labs/mcp-go`](https://github.com/mark3labs/mcp-go).

It ships with small, working examples of the three MCP primitives so you can
copy-paste-rename your way to a real server:

| Primitive | File                              | Examples                       |
|-----------|-----------------------------------|--------------------------------|
| Tools     | `internal/tools/tools.go`         | `echo`, `add`                  |
| Resources | `internal/resources/resources.go` | `info://server`, `item://{id}` |
| Prompts   | `internal/prompts/prompts.go`     | `code-review`, `summarise`     |

Both the `stdio` and `sse` transports are wired up in `main.go`.

---

## Quick start

First, rename the Go module to your own path (the template ships under a
placeholder):

```sh
go mod edit -module github.com/you/your-server
# then update the three imports in main.go to match
```

```sh
# build
make build

# run over stdio (what most MCP clients use)
make run-stdio

# run over SSE on :8080
make run-sse

# run tests
make test
```

Point your MCP client at the `mcp-server` binary (stdio) or at
`http://localhost:8080` (SSE).

### SSE flags

```
-transport sse
-addr      :8080         # host:port to bind
-base-url  ""            # optional public URL advertised to clients;
                         # defaults to http://<addr> with localhost fallback
                         # when the bind host is empty, 0.0.0.0, or ::
```

---

## Repository layout

```
.
├── main.go                     # entrypoint + transport wiring
├── internal/
│   ├── tools/tools.go          # tool definitions and handlers
│   ├── resources/resources.go  # static resources and URI templates
│   └── prompts/prompts.go      # parameterised prompt templates
├── Makefile
└── go.mod
```

Each `internal/*` package exposes a single `Register(*server.MCPServer)`
function called from `main.go`. Adding a new primitive is a matter of adding
a definition, a handler, and one line to `Register`.

---

## Tutorial: extending the template

The three sections below walk through adding one tool, one resource, and one
prompt. The pattern is always: **schema → handler → register → test.**

### 1. Add a new tool

Say you want a `multiply` tool that takes two numbers.

**1a. Define the tool schema** in `internal/tools/tools.go`:

```go
func multiplyTool() mcp.Tool {
    return mcp.NewTool("multiply",
        mcp.WithDescription("Multiply two numbers and return the product"),
        mcp.WithNumber("a",
            mcp.Required(),
            mcp.Description("First factor"),
        ),
        mcp.WithNumber("b",
            mcp.Required(),
            mcp.Description("Second factor"),
        ),
    )
}
```

**1b. Write the handler.** Handlers return a `*mcp.CallToolResult`; use
`NewToolResultError` for user-facing validation errors and return a real
`error` only for unexpected failures.

```go
func multiplyHandler(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    a, err := req.RequireFloat("a")
    if err != nil {
        return mcp.NewToolResultError(err.Error()), nil
    }
    b, err := req.RequireFloat("b")
    if err != nil {
        return mcp.NewToolResultError(err.Error()), nil
    }
    return mcp.NewToolResultText(strconv.FormatFloat(a*b, 'f', -1, 64)), nil
}
```

**1c. Register it** inside `Register`:

```go
func Register(s *server.MCPServer) {
    s.AddTool(echoTool(), echoHandler)
    s.AddTool(addTool(), addHandler)
    s.AddTool(multiplyTool(), multiplyHandler) // <-- new
}
```

**1d. Test it** in `internal/tools/tools_test.go`:

```go
func TestMultiplyHandler(t *testing.T) {
    got := callTool(t, multiplyHandler, map[string]any{"a": 4.0, "b": 2.5})
    if got != "10" {
        t.Fatalf("multiply(4, 2.5): got %q, want %q", got, "10")
    }
}
```

Available argument helpers on `CallToolRequest` include `RequireString`,
`RequireInt`, `RequireFloat`, `RequireBool`, plus their `*Slice` variants and
non-`Require` (optional) forms.

### 2. Add a new resource

Resources come in two flavours:

- **Static** — a single fixed URI (e.g. `config://app`).
- **Template** — a URI with variables (e.g. `user://{id}`); `mark3labs/mcp-go`
  matches the incoming URI against the template and puts the matched
  variables in `req.Params.Arguments`.

**Static example** — add a `config://app` resource:

```go
s.AddResource(
    mcp.NewResource(
        "config://app",
        "App config",
        mcp.WithResourceDescription("Current application configuration"),
        mcp.WithMIMEType("application/json"),
    ),
    func(_ context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
        return []mcp.ResourceContents{
            mcp.TextResourceContents{
                URI:      "config://app",
                MIMEType: "application/json",
                Text:     `{"feature_x": true}`,
            },
        }, nil
    },
)
```

**Template example** — add a `user://{id}` resource. Note how `id` is read
from `req.Params.Arguments`, **not** parsed from `req.Params.URI`:

```go
s.AddResourceTemplate(
    mcp.NewResourceTemplate(
        "user://{id}",
        "User by ID",
        mcp.WithTemplateDescription("Fetch a user record by ID"),
        mcp.WithTemplateMIMEType("application/json"),
    ),
    func(_ context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
        id, _ := req.Params.Arguments["id"].(string)
        if id == "" {
            return nil, fmt.Errorf("missing id in URI %q", req.Params.URI)
        }
        // ... look up the user ...
        return []mcp.ResourceContents{
            mcp.TextResourceContents{
                URI:      req.Params.URI,
                MIMEType: "application/json",
                Text:     fmt.Sprintf(`{"id": %q}`, id),
            },
        }, nil
    },
)
```

For binary payloads, return `mcp.BlobResourceContents` instead of
`TextResourceContents`.

### 3. Add a new prompt

Prompts are parameterised message templates the client can fetch and send to
the model.

```go
func translatePrompt() mcp.Prompt {
    return mcp.NewPrompt("translate",
        mcp.WithPromptDescription("Translate text into a target language"),
        mcp.WithArgument("text",
            mcp.ArgumentDescription("Text to translate"),
            mcp.RequiredArgument(),
        ),
        mcp.WithArgument("target",
            mcp.ArgumentDescription("Target language (default: English)"),
        ),
    )
}

func translateHandler(_ context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
    text := req.Params.Arguments["text"]
    target := req.Params.Arguments["target"]
    if target == "" {
        target = "English"
    }

    content := fmt.Sprintf("Translate the following text into %s:\n\n%s", target, text)

    return mcp.NewGetPromptResult(
        "Translate prompt",
        []mcp.PromptMessage{
            mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(content)),
        },
    ), nil
}
```

Register it alongside the others in `prompts.Register`.

### 4. Use context and external services

Every handler receives a `context.Context`. Propagate it to outbound HTTP
calls, database queries, and any other blocking work so clients can cancel
in-flight requests.

```go
func weatherHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    city, err := req.RequireString("city")
    if err != nil {
        return mcp.NewToolResultError(err.Error()), nil
    }
    data, err := fetchWeather(ctx, city) // pass ctx through
    if err != nil {
        return nil, err
    }
    return mcp.NewToolResultText(data), nil
}
```

If your handler needs shared state (database handle, HTTP client, config),
the cleanest pattern is to define a struct with methods instead of package-
level functions:

```go
type WeatherService struct{ http *http.Client }

func (w *WeatherService) Register(s *server.MCPServer) {
    s.AddTool(weatherTool(), w.handle)
}

func (w *WeatherService) handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // ... use w.http ...
}
```

Then wire it up from `main.go`:

```go
weather := &tools.WeatherService{http: &http.Client{Timeout: 5 * time.Second}}
weather.Register(s)
```

### 5. Error handling conventions

- **User input is wrong** → `mcp.NewToolResultError("...")`, return `nil` error.
  The client sees a tool error but the session stays healthy.
- **Something exploded server-side** → return a real `error`. The framework
  converts it into a JSON-RPC error for the client.
- **Validation at the boundary only** — once a handler has valid inputs,
  trust internal helpers; don't re-check every layer.

### 6. Running against a real client

- **stdio**: most desktop MCP clients (Claude Desktop, editors) expect a
  binary they can spawn. Point them at `./mcp-server`. No flags needed.
- **SSE**: useful for remote servers or debugging with `curl`. Run
  `make run-sse` and connect to `http://localhost:8080`.

---

## Dependencies

- [`github.com/mark3labs/mcp-go`](https://github.com/mark3labs/mcp-go) — MCP
  server framework (stdio and SSE transports, schema builders, request types).

Run `make tidy` after adding new imports.

---

## License

Use it however you like — it's a template.
