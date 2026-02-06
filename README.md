# MCP Gateway on Envoy

Open-source MCP Gateway built on Envoy/Envoy Gateway, optimized for fast developer onboarding and production-ready deployment in minutes.

## Current Status

- `Phase`: Phase 1 runtime baseline
- `Implemented`: config schema/validation, CLI, runtime HTTP gateway, manifest render/apply, local Envoy compose stack
- `Next`: Envoy Gateway CRD apply and deeper transport/auth integrations

## Quick Start (Local in Minutes)

```bash
cd /Users/djsam/codex/mcp-gateway-envoy
docker compose -f deploy/local/docker-compose.yml up --build
```

Test through Envoy (`:10000`):

```bash
curl -H "Authorization: Bearer dev-token" http://localhost:10000/mcp/weather
```

Health checks:

```bash
curl http://localhost:18080/healthz
curl http://localhost:18080/readyz
```

## CLI

```bash
go run ./cmd/gateway init --output gateway.yaml
go run ./cmd/gateway validate --file gateway.yaml
go run ./cmd/gateway plan --file gateway.yaml
go run ./cmd/gateway render --file gateway.yaml --namespace mcp-gateway --output manifests.yaml
go run ./cmd/gateway apply --file gateway.yaml --namespace mcp-gateway --dry-run
go run ./cmd/gateway serve --file gateway.yaml
```

## Repository Layout

- `docs/research.md`: OSS landscape and feature analysis
- `docs/requirements.md`: product and technical requirements
- `docs/architecture.md`: architecture decisions (ADRs)
- `docs/roadmap.md`: delivery phases and milestones
- `cmd/gateway`: CLI entrypoint
- `internal/config`: config schema, loading, validation, template
- `internal/controller`: resource planning + Kubernetes manifest rendering
- `internal/runtime`: local server runtime + kubectl apply integration
- `deploy/examples`: sample gateway config
- `deploy/local`: local Docker Compose + Envoy config

## Near-term Deliverables

1. Reconciler that applies Envoy Gateway CRDs directly.
2. stdio bridge runtime process manager.
3. Helm chart for production bootstrap.
