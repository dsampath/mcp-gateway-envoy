package controller

import (
	"strings"
	"testing"

	"github.com/djsam/mcp-gateway-envoy/internal/config"
)

func TestRenderManifests(t *testing.T) {
	cfg := &config.Config{
		APIVersion: "mcp.envoy.io/v1alpha1",
		Kind:       "GatewayConfig",
		Gateway: config.Gateway{
			Name:       "mcp-gateway",
			ListenAddr: ":8080",
		},
		Auth: config.AuthDefaults{RequireAuth: true},
		Servers: []config.Server{
			{Name: "s1", Transport: "http", URL: "http://example"},
		},
		Routes: []config.Route{
			{Name: "r1", Path: "/mcp", Server: "s1", Policy: config.RoutePolicy{TimeoutMs: 1000}},
		},
	}

	manifest, err := RenderManifests(cfg, "mcp-gateway", "example/image:latest")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	text := string(manifest)
	for _, needle := range []string{"kind: Namespace", "kind: Deployment", "kind: MCPRoute", "kind: MCPAuthPolicy"} {
		if !strings.Contains(text, needle) {
			t.Fatalf("manifest missing %q", needle)
		}
	}
}
