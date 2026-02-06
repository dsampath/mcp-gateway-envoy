package controller

import (
	"fmt"
	"strings"

	"github.com/djsam/mcp-gateway-envoy/internal/config"
)

// Resource is a lightweight phase-0 representation of generated resources.
type Resource struct {
	Kind string         `json:"kind"`
	Name string         `json:"name"`
	Spec map[string]any `json:"spec"`
}

// BuildResources translates gateway config into Envoy Gateway-oriented resources.
func BuildResources(cfg *config.Config) ([]Resource, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	resources := make([]Resource, 0, len(cfg.Routes)*3)
	for _, route := range cfg.Routes {
		resources = append(resources,
			Resource{
				Kind: "BackendRef",
				Name: route.Name + "-backend",
				Spec: map[string]any{
					"server":    route.Server,
					"transport": lookupServerTransport(cfg, route.Server),
				},
			},
			Resource{
				Kind: "MCPRoute",
				Name: route.Name,
				Spec: map[string]any{
					"path":         route.Path,
					"backendRef":   route.Name + "-backend",
					"timeoutMs":    route.Policy.TimeoutMs,
					"retryCount":   route.Policy.RetryCount,
					"rateLimitRps": route.Policy.RateLimitRPS,
				},
			},
		)

		if authKind := routeAuthKind(cfg, route); authKind != "none" {
			resources = append(resources, Resource{
				Kind: "MCPAuthPolicy",
				Name: route.Name + "-auth",
				Spec: map[string]any{
					"route": route.Name,
					"type":  authKind,
				},
			})
		}
	}

	return resources, nil
}

func lookupServerTransport(cfg *config.Config, serverName string) string {
	for _, s := range cfg.Servers {
		if s.Name == serverName {
			return s.Transport
		}
	}
	return "unknown"
}

func routeAuthKind(cfg *config.Config, route config.Route) string {
	if route.Auth == nil {
		if cfg.Auth.RequireAuth {
			return "apiKey"
		}
		return "none"
	}
	kind := strings.TrimSpace(route.Auth.Type)
	if kind == "" {
		if cfg.Auth.RequireAuth {
			return "apiKey"
		}
		return "none"
	}
	return kind
}
