package config

import (
	"fmt"
	"strings"
)

// Validate applies schema and semantic validation rules.
func (c Config) Validate() error {
	if c.APIVersion != "mcp.envoy.io/v1alpha1" {
		return fmt.Errorf("apiVersion must be mcp.envoy.io/v1alpha1")
	}
	if c.Kind != "GatewayConfig" {
		return fmt.Errorf("kind must be GatewayConfig")
	}
	if strings.TrimSpace(c.Gateway.Name) == "" {
		return fmt.Errorf("gateway.name is required")
	}
	if strings.TrimSpace(c.Gateway.ListenAddr) == "" {
		return fmt.Errorf("gateway.listenAddr is required")
	}
	if len(c.Servers) == 0 {
		return fmt.Errorf("servers must include at least one server")
	}
	if len(c.Routes) == 0 {
		return fmt.Errorf("routes must include at least one route")
	}

	seenServers := map[string]struct{}{}
	for _, s := range c.Servers {
		if strings.TrimSpace(s.Name) == "" {
			return fmt.Errorf("servers[].name is required")
		}
		if _, ok := seenServers[s.Name]; ok {
			return fmt.Errorf("duplicate server name: %s", s.Name)
		}
		seenServers[s.Name] = struct{}{}

		switch s.Transport {
		case "http":
			if strings.TrimSpace(s.URL) == "" {
				return fmt.Errorf("server %q transport http requires url", s.Name)
			}
		case "stdio":
			if strings.TrimSpace(s.Command) == "" {
				return fmt.Errorf("server %q transport stdio requires command", s.Name)
			}
		default:
			return fmt.Errorf("server %q has unsupported transport %q", s.Name, s.Transport)
		}
	}

	seenRoutes := map[string]struct{}{}
	for _, r := range c.Routes {
		if strings.TrimSpace(r.Name) == "" {
			return fmt.Errorf("routes[].name is required")
		}
		if _, ok := seenRoutes[r.Name]; ok {
			return fmt.Errorf("duplicate route name: %s", r.Name)
		}
		seenRoutes[r.Name] = struct{}{}

		if strings.TrimSpace(r.Path) == "" || !strings.HasPrefix(r.Path, "/") {
			return fmt.Errorf("route %q path must start with '/'", r.Name)
		}
		if _, ok := seenServers[r.Server]; !ok {
			return fmt.Errorf("route %q references unknown server %q", r.Name, r.Server)
		}
		if r.Policy.TimeoutMs < 0 || r.Policy.RetryCount < 0 || r.Policy.RateLimitRPS < 0 {
			return fmt.Errorf("route %q policy values must be >= 0", r.Name)
		}
		if r.Auth != nil {
			switch r.Auth.Type {
			case "apiKey":
				if strings.TrimSpace(r.Auth.HeaderName) == "" || len(r.Auth.APIKeys) == 0 {
					return fmt.Errorf("route %q apiKey auth requires headerName and apiKeys", r.Name)
				}
			case "jwt":
				if strings.TrimSpace(r.Auth.Issuer) == "" || strings.TrimSpace(r.Auth.Audience) == "" {
					return fmt.Errorf("route %q jwt auth requires issuer and audience", r.Name)
				}
			case "none":
			default:
				return fmt.Errorf("route %q auth type must be apiKey, jwt, or none", r.Name)
			}
		}
	}

	return nil
}
