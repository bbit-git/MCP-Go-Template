package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/andrzej/mcp-go-template/internal/prompts"
	"github.com/andrzej/mcp-go-template/internal/resources"
	"github.com/andrzej/mcp-go-template/internal/tools"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	transport := flag.String("transport", "stdio", "Transport type: stdio or sse")
	addr := flag.String("addr", ":8080", "Address to listen on, host:port (SSE transport only)")
	baseURL := flag.String("base-url", "", "Public base URL advertised by the SSE server (defaults to http://<addr>)")
	flag.Parse()

	s := server.NewMCPServer(
		"mcp-go-template",
		"0.1.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
	)

	tools.Register(s)
	resources.Register(s)
	prompts.Register(s)

	switch *transport {
	case "stdio":
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("stdio server error: %v", err)
		}
	case "sse":
		url := *baseURL
		if url == "" {
			host, port, err := net.SplitHostPort(*addr)
			if err != nil {
				log.Fatalf("invalid -addr %q: %v", *addr, err)
			}
			if host == "" || host == "0.0.0.0" || host == "::" {
				host = "localhost"
			}
			url = fmt.Sprintf("http://%s", net.JoinHostPort(host, port))
		}
		sseServer := server.NewSSEServer(s, server.WithBaseURL(url))
		log.Printf("SSE server listening on %s (base URL %s)", *addr, url)
		if err := sseServer.Start(*addr); err != nil {
			log.Fatalf("SSE server error: %v", err)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown transport %q (use stdio or sse)\n", *transport)
		os.Exit(1)
	}
}
