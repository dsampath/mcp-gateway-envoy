# MCP Gateway on Envoy

Open-source MCP Gateway built on Envoy/Envoy Gateway, optimized for fast developer onboarding and production-ready deployment in minutes.

## Current Status

- `Phase`: Phase 0 foundation in progress
- `Implemented`: config schema + validation, CLI bootstrap, resource planning skeleton
- `Next`: runtime wiring to apply generated resources + local compose quickstart

## Why This Project

- Drop-in MCP gateway for local dev and production
- Minimal setup (`docker compose up` local, Helm on Kubernetes)
- Secure-by-default routing and auth policy
- Transport interoperability (HTTP streamable + stdio bridge)
- Strong operator experience (validation, metrics, traces, health)

## Quick Start (Current)

```bash
go run ./cmd/gateway init --output gateway.yaml
go run ./cmd/gateway validate --file gateway.yaml
go run ./cmd/gateway plan --file gateway.yaml
```

## Repository Layout

- `docs/research.md`: OSS landscape and feature analysis
- `docs/requirements.md`: product and technical requirements
- `docs/architecture.md`: architecture decisions (ADRs)
- `docs/roadmap.md`: delivery phases and milestones
- `cmd/gateway`: CLI entrypoint
- `internal/config`: config schema, loading, validation, template
- `internal/controller`: resource plan generation (phase-0)
- `deploy/examples`: example gateway config

## Near-term Deliverables

1. Implement runtime reconciler and Envoy Gateway resource apply path.
2. Ship one-command local deployment profile.
3. Publish Helm chart and docs for production bootstrap.
