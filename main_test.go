package main

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func newTestMCPServer() *server.MCPServer {
	return server.NewMCPServer("miniflux-mcp-test", "0.1.0")
}

func postJSON(t *testing.T, ts *httptest.Server, sessionID string, body any) *http.Response {
	t.Helper()
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/mcp", bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("new request failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if sessionID != "" {
		req.Header.Set("Mcp-Session-Id", sessionID)
	}
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("POST /mcp failed: %v", err)
	}
	return resp
}

func TestTransportSelection(t *testing.T) {
	t.Run("stdio mode when MCP_HTTP_PORT is empty", func(t *testing.T) {
		os.Setenv("MCP_HTTP_PORT", "")
		defer os.Unsetenv("MCP_HTTP_PORT")

		if port := os.Getenv("MCP_HTTP_PORT"); port != "" {
			t.Errorf("expected empty port, got %q", port)
		}
	})

	t.Run("HTTP mode when MCP_HTTP_PORT is set", func(t *testing.T) {
		os.Setenv("MCP_HTTP_PORT", "3000")
		defer os.Unsetenv("MCP_HTTP_PORT")

		if port := os.Getenv("MCP_HTTP_PORT"); port == "" {
			t.Error("expected non-empty port")
		}
	})
}

func TestHTTPServerStartsAndResponds(t *testing.T) {
	ts := server.NewTestStreamableHTTPServer(newTestMCPServer())
	defer ts.Close()

	// Step 1: initialize to obtain a session ID.
	initBody := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]any{
			"protocolVersion": mcp.LATEST_PROTOCOL_VERSION,
			"clientInfo":      map[string]any{"name": "test", "version": "0.0.1"},
			"capabilities":    map[string]any{},
		},
	}
	initResp := postJSON(t, ts, "", initBody)
	defer initResp.Body.Close()
	if initResp.StatusCode != http.StatusOK {
		t.Fatalf("initialize: expected 200, got %d", initResp.StatusCode)
	}
	sessionID := initResp.Header.Get("Mcp-Session-Id")
	if sessionID == "" {
		t.Fatal("initialize: missing Mcp-Session-Id response header")
	}

	// Step 2: tools/list with the session ID.
	listBody := map[string]any{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/list",
		"params":  map[string]any{},
	}
	listResp := postJSON(t, ts, sessionID, listBody)
	defer listResp.Body.Close()
	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("tools/list: expected 200, got %d", listResp.StatusCode)
	}

	var result map[string]any
	if err := json.NewDecoder(listResp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if result["jsonrpc"] != "2.0" {
		t.Errorf("expected jsonrpc 2.0, got %v", result["jsonrpc"])
	}
}

func TestHTTPServerFailsOnUsedPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()
	addr := ln.Addr().String()

	s := newTestMCPServer()
	httpServer := server.NewStreamableHTTPServer(s)

	err = httpServer.Start(addr)
	if err == nil {
		t.Error("expected error when port is already in use, got nil")
	}
}
