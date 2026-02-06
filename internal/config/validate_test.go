package config

import "testing"

func TestValidateTemplate(t *testing.T) {
	cfg, err := LoadFile("../../deploy/examples/gateway.yaml")
	if err != nil {
		t.Fatalf("expected template config to validate, got error: %v", err)
	}
	if len(cfg.Servers) != 2 || len(cfg.Routes) != 2 {
		t.Fatalf("unexpected counts: servers=%d routes=%d", len(cfg.Servers), len(cfg.Routes))
	}
}

func TestValidateRouteUnknownServer(t *testing.T) {
	cfg := Config{
		APIVersion: "mcp.envoy.io/v1alpha1",
		Kind:       "GatewayConfig",
		Gateway: Gateway{
			Name:       "gw",
			ListenAddr: ":8080",
		},
		Servers: []Server{{Name: "a", Transport: "http", URL: "http://example"}},
		Routes:  []Route{{Name: "r1", Path: "/mcp", Server: "missing"}},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for unknown server")
	}
}
