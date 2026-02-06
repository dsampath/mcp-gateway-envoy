# Requirements v0.2

Date: 2026-02-06

## Product Objectives

- Let developers route to MCP servers in under 10 minutes.
- Support local dev, single-node deploy, and Kubernetes production.
- Provide secure-by-default gateway behavior with minimal config.

## Scope Decisions (Reviewed)

- Control path: Envoy Gateway CRD-driven mapping in v1.
- Security default: deny-by-default (`requireAuth: true`).
- Tenant model: single-tenant in v1; multi-tenant in post-MVP.
- Transport v1: streamable HTTP + stdio bridge; SSE optional compatibility mode.

## Personas

- App developer: wants fast local setup and stable endpoint.
- Platform engineer: wants policy control, observability, and straightforward rollout.
- Security engineer: wants centralized auth, auditability, least privilege.

## Functional Requirements

1. Server Registry
- Register MCP servers via declarative YAML config.
- MVP supports static definitions only.
- Health/readiness state surfaced per server.

2. Transport Interop
- Route streamable HTTP MCP traffic.
- Bridge stdio-based servers with adapter contract.
- Optional SSE compatibility mode for client fallback.

3. Routing and Policy
- MCP-aware route objects mapped to Envoy Gateway resources.
- Per-route auth policy with `apiKey`, `jwt`, `none`.
- Route-level timeout/retry/rate limits.

4. Security
- TLS termination at gateway.
- Upstream TLS support.
- Secret references for credentials/tokens.
- Audit log events for auth and route access decisions.

5. Developer Experience
- Single command local startup (`docker compose up`) in MVP.
- `gateway init` to scaffold config + sample servers.
- `gateway validate` for config lint and semantic checks.
- Actionable error output with fix hints.

6. Operations
- Prometheus metrics for traffic, latency, errors.
- OpenTelemetry traces.
- Structured logs with correlation IDs.
- `/healthz` and `/readyz` endpoints.

7. Deployment
- Local Docker Compose profile.
- Kubernetes Helm chart.
- Example manifests for quick POC environments.

## Non-Functional Requirements

- Reliability: no single failing upstream degrades unrelated routes.
- Performance: p95 gateway overhead target under 20ms in baseline environment.
- Security baseline: unauthenticated routes must be explicit (`auth.type: none`).
- Portability: macOS/Linux for dev and Kubernetes for production.
- Maintainability: modular extension points for auth, registry, transport adapters.

## MVP Scope (Must Have)

- Static MCP server registry config (`gateway.yaml`).
- Envoy Gateway-backed route generation for MCP endpoints.
- Route auth with API key and JWT.
- stdio adapter integration contract.
- Docker Compose one-command startup.
- Basic metrics/logging/health endpoints.

## Post-MVP (Should Have)

- Dynamic registry providers (GitOps/remote catalog).
- Multi-tenant policy model and RBAC.
- Web admin/API for runtime management.
- Advanced policy plugins, quotas, and chargeback.

## Acceptance Criteria (Phase 0)

- `gateway init` creates a valid config template.
- `gateway validate` passes for template and fails on invalid references.
- `gateway plan` outputs MCP route/auth/backend resources from config.
- Architecture decisions documented in `docs/architecture.md`.
