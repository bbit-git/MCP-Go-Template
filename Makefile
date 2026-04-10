BINARY := mcp-server
GO     := go

.PHONY: build run-stdio run-sse test lint tidy clean

build:
	$(GO) build -buildvcs=false -o $(BINARY) .

run-stdio: build
	./$(BINARY) -transport stdio

run-sse: build
	./$(BINARY) -transport sse -addr :8080

test:
	$(GO) test -buildvcs=false ./...

lint:
	golangci-lint run ./...

tidy:
	$(GO) mod tidy

clean:
	rm -f $(BINARY)
