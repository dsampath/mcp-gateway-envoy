package config

const DefaultTemplateYAML = `apiVersion: mcp.envoy.io/v1alpha1
kind: GatewayConfig
gateway:
  name: mcp-gateway
  listenAddr: ":8080"
  adminAddr: ":9090"
  logLevel: info
auth:
  requireAuth: true
servers:
  - name: weather-http
    transport: http
    url: http://weather-mcp:8000
  - name: filesystem-local
    transport: stdio
    command: npx
    args: ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"]
routes:
  - name: weather
    path: /mcp/weather
    server: weather-http
    auth:
      type: jwt
      issuer: https://issuer.example.com
      audience: mcp-gateway
    policy:
      timeoutMs: 10000
      retryCount: 1
      rateLimitRps: 20
  - name: filesystem
    path: /mcp/fs
    server: filesystem-local
    auth:
      type: apiKey
      headerName: X-API-Key
      apiKeys: ["replace-me"]
    policy:
      timeoutMs: 15000
      retryCount: 0
      rateLimitRps: 10
`
