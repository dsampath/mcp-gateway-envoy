# OSS MCP Gateway Landscape (Initial)

Date: 2026-02-06

## Selection Criteria

- Open-source availability and active ecosystem presence
- Relevance to MCP gateway/proxy behavior
- Signals for production usage: deployment model, security, operations, extensibility

## Projects Reviewed

### 1) Envoy AI Gateway (Envoy ecosystem)

- Source: https://github.com/envoyproxy/ai-gateway
- Why relevant: Envoy-native AI/MCP gateway trajectory with Kubernetes-first operations.
- Notable signals:
  - Quickstart supports `kind`, Helm install, and sample gateway deployment from README.
  - Envoy Gateway v0.5 release notes call out MCP-specific APIs (`MCPRoute`, `MCPAuthPolicy`) and backend TLS support for MCP routes.
- Takeaway for us:
  - Build around Envoy Gateway APIs instead of bespoke L7 routing logic.
  - Prioritize MCP policy surface compatible with Envoy Gateway concepts.

### 2) Microsoft MCP Gateway

- Source: https://github.com/microsoft/mcp-gateway
- Why relevant: Comprehensive OSS MCP gateway implementation and reference architecture.
- Notable signals:
  - Explicit positioning as an MCP Gateway service.
  - README highlights broad plugin/component model and deployment patterns.
- Takeaway for us:
  - Treat extensibility as a first-class concern (auth providers, registries, transforms, policy plugins).
  - Ship a clear local-to-cloud promotion path.

### 3) IBM MCP Context Forge

- Source: https://github.com/IBM/mcp-context-forge
- Why relevant: MCP gateway plus context engineering platform pattern.
- Notable signals:
  - Multi-tenant enterprise orientation with administration and management features.
  - Strong packaging/docs posture for self-hosting workflows.
- Takeaway for us:
  - Include tenant and workspace boundaries early in data model.
  - Separate “simple default mode” from enterprise controls to preserve UX.

### 4) Supergateway (transport bridge)

- Source: https://github.com/supercorp-ai/supergateway
- Why relevant: Focused bridge for stdio MCP servers to networked MCP transports.
- Notable signals:
  - Purpose-built transport conversion for MCP server connectivity.
  - Practical for local dev and legacy server onboarding.
- Takeaway for us:
  - MVP should include robust transport adaptation path (stdio <-> HTTP streamable/SSE).

## Shared Patterns Across Popular Projects

- Quick local startup path is essential for adoption.
- Configuration must be declarative with sane defaults.
- Security is expected out of the box (authn/authz, TLS, token handling).
- Operators need production telemetry without extra glue.
- Teams want central server registry/discovery and policy enforcement.

## Gaps/Opportunity

- “Deploy in minutes” with Envoy-native control-plane UX is still fragmented.
- Opportunity: a focused Envoy-based MCP gateway that starts tiny, scales cleanly, and keeps policy/ops integrated from day one.

## Source Links

- Envoy AI Gateway README: https://raw.githubusercontent.com/envoyproxy/ai-gateway/main/README.md
- Envoy Gateway v0.5 release notes (`MCPRoute`, `MCPAuthPolicy`): https://raw.githubusercontent.com/envoyproxy/gateway/main/site/content/en/news/releases/v0.5/v0.5.0.md
- Microsoft MCP Gateway README: https://raw.githubusercontent.com/microsoft/mcp-gateway/main/README.md
- IBM MCP Context Forge README: https://raw.githubusercontent.com/IBM/mcp-context-forge/main/README.md
- Supergateway README: https://raw.githubusercontent.com/supercorp-ai/supergateway/main/README.md
