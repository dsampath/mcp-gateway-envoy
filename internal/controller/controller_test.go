package controller

import (
	"testing"

	"github.com/djsam/mcp-gateway-envoy/internal/config"
)

func TestBuildResources(t *testing.T) {
	cfg := &config.Config{
		APIVersion: "mcp.envoy.io/v1alpha1",
		Kind:       "GatewayConfig",
		Gateway: config.Gateway{
			Name:       "gw",
			ListenAddr: ":8080",
		},
		Auth: config.AuthDefaults{RequireAuth: true},
		Servers: []config.Server{
			{Name: "s1", Transport: "http", URL: "http://example"},
		},
		Routes: []config.Route{
			{Name: "r1", Path: "/mcp", Server: "s1", Policy: config.RoutePolicy{TimeoutMs: 10}},
		},
	}
	resources, err := BuildResources(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resources) != 3 {
		t.Fatalf("expected 3 resources, got %d", len(resources))
	}
}
