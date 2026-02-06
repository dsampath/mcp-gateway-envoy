package runtime

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/djsam/mcp-gateway-envoy/internal/config"
)

func TestHTTPProxyRoute(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("upstream-ok"))
	}))
	defer upstream.Close()

	cfg := &config.Config{
		APIVersion: "mcp.envoy.io/v1alpha1",
		Kind:       "GatewayConfig",
		Gateway:    config.Gateway{Name: "gw", ListenAddr: ":0"},
		Auth:       config.AuthDefaults{RequireAuth: false},
		Servers:    []config.Server{{Name: "s1", Transport: "http", URL: upstream.URL}},
		Routes:     []config.Route{{Name: "r1", Path: "/mcp", Server: "s1"}},
	}
	s := NewServer(cfg)

	req := httptest.NewRequest(http.MethodGet, "/mcp", nil)
	rr := httptest.NewRecorder()
	s.handleRequest(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body, _ := io.ReadAll(rr.Body)
	if string(body) != "upstream-ok" {
		t.Fatalf("unexpected body: %s", string(body))
	}
}

func TestAPIKeyAuthDenied(t *testing.T) {
	cfg := &config.Config{
		APIVersion: "mcp.envoy.io/v1alpha1",
		Kind:       "GatewayConfig",
		Gateway:    config.Gateway{Name: "gw", ListenAddr: ":0"},
		Auth:       config.AuthDefaults{RequireAuth: true},
		Servers:    []config.Server{{Name: "s1", Transport: "stdio", Command: "echo"}},
		Routes: []config.Route{{
			Name:   "r1",
			Path:   "/mcp",
			Server: "s1",
			Auth: &config.RouteAuth{
				Type:       "apiKey",
				HeaderName: "X-API-Key",
				APIKeys:    []string{"secret"},
			},
		}},
	}
	s := NewServer(cfg)

	req := httptest.NewRequest(http.MethodGet, "/mcp", nil)
	rr := httptest.NewRecorder()
	s.handleRequest(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}
