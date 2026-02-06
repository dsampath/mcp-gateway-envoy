package controller

import (
	"fmt"
	"strings"

	"github.com/djsam/mcp-gateway-envoy/internal/config"
	"gopkg.in/yaml.v3"
)

// RenderManifests renders Kubernetes YAML for gateway runtime + MCP route policy resources.
func RenderManifests(cfg *config.Config, namespace string, image string) ([]byte, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	if strings.TrimSpace(namespace) == "" {
		namespace = "default"
	}
	if strings.TrimSpace(image) == "" {
		image = "ghcr.io/dsampath/mcp-gateway-envoy:latest"
	}

	docs := make([]map[string]any, 0, 4+len(cfg.Routes)*2)
	docs = append(docs,
		namespaceDoc(namespace),
		configMapDoc(cfg, namespace),
		deploymentDoc(cfg, namespace, image),
		serviceDoc(cfg, namespace),
	)

	for _, r := range cfg.Routes {
		docs = append(docs, routeDoc(r, namespace))
		if kind := routeAuthKind(cfg, r); kind != "none" {
			docs = append(docs, authPolicyDoc(r, namespace, kind))
		}
	}

	var b strings.Builder
	for i, d := range docs {
		if i > 0 {
			b.WriteString("---\n")
		}
		y, err := yaml.Marshal(d)
		if err != nil {
			return nil, fmt.Errorf("marshal manifest: %w", err)
		}
		b.Write(y)
	}
	return []byte(b.String()), nil
}

func namespaceDoc(namespace string) map[string]any {
	return map[string]any{
		"apiVersion": "v1",
		"kind":       "Namespace",
		"metadata": map[string]any{
			"name": namespace,
		},
	}
}

func configMapDoc(cfg *config.Config, namespace string) map[string]any {
	return map[string]any{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata": map[string]any{
			"name":      cfg.Gateway.Name + "-config",
			"namespace": namespace,
		},
		"data": map[string]any{
			"gateway.yaml": configToYAML(cfg),
		},
	}
}

func deploymentDoc(cfg *config.Config, namespace, image string) map[string]any {
	return map[string]any{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata": map[string]any{
			"name":      cfg.Gateway.Name,
			"namespace": namespace,
			"labels": map[string]any{
				"app": cfg.Gateway.Name,
			},
		},
		"spec": map[string]any{
			"replicas": 1,
			"selector": map[string]any{
				"matchLabels": map[string]any{
					"app": cfg.Gateway.Name,
				},
			},
			"template": map[string]any{
				"metadata": map[string]any{
					"labels": map[string]any{
						"app": cfg.Gateway.Name,
					},
				},
				"spec": map[string]any{
					"containers": []map[string]any{
						{
							"name":  "gateway",
							"image": image,
							"args":  []string{"serve", "--file", "/etc/mcp-gateway/gateway.yaml"},
							"ports": []map[string]any{{"name": "http", "containerPort": parsePort(cfg.Gateway.ListenAddr)}},
							"volumeMounts": []map[string]any{
								{
									"name":      "config",
									"mountPath": "/etc/mcp-gateway",
								},
							},
						},
					},
					"volumes": []map[string]any{
						{
							"name": "config",
							"configMap": map[string]any{
								"name": cfg.Gateway.Name + "-config",
							},
						},
					},
				},
			},
		},
	}
}

func serviceDoc(cfg *config.Config, namespace string) map[string]any {
	port := parsePort(cfg.Gateway.ListenAddr)
	return map[string]any{
		"apiVersion": "v1",
		"kind":       "Service",
		"metadata": map[string]any{
			"name":      cfg.Gateway.Name,
			"namespace": namespace,
		},
		"spec": map[string]any{
			"selector": map[string]any{"app": cfg.Gateway.Name},
			"ports": []map[string]any{
				{
					"name":       "http",
					"port":       port,
					"targetPort": port,
				},
			},
		},
	}
}

func routeDoc(route config.Route, namespace string) map[string]any {
	return map[string]any{
		"apiVersion": "mcp.envoy.io/v1alpha1",
		"kind":       "MCPRoute",
		"metadata": map[string]any{
			"name":      route.Name,
			"namespace": namespace,
		},
		"spec": map[string]any{
			"path":       route.Path,
			"backendRef": route.Server,
			"policy": map[string]any{
				"timeoutMs":    route.Policy.TimeoutMs,
				"retryCount":   route.Policy.RetryCount,
				"rateLimitRps": route.Policy.RateLimitRPS,
			},
		},
	}
}

func authPolicyDoc(route config.Route, namespace, authType string) map[string]any {
	return map[string]any{
		"apiVersion": "mcp.envoy.io/v1alpha1",
		"kind":       "MCPAuthPolicy",
		"metadata": map[string]any{
			"name":      route.Name + "-auth",
			"namespace": namespace,
		},
		"spec": map[string]any{
			"targetRoute": route.Name,
			"type":        authType,
		},
	}
}

func configToYAML(cfg *config.Config) string {
	b, _ := yaml.Marshal(cfg)
	return string(b)
}

func parsePort(listenAddr string) int {
	parts := strings.Split(strings.TrimSpace(listenAddr), ":")
	if len(parts) == 0 {
		return 8080
	}
	p := parts[len(parts)-1]
	if p == "" {
		return 8080
	}
	var port int
	_, err := fmt.Sscanf(p, "%d", &port)
	if err != nil || port <= 0 {
		return 8080
	}
	return port
}
