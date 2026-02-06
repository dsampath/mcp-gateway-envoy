# Build Roadmap

Date: 2026-02-06

## Phase 0: Foundation (Completed)

- Architecture decisions documented.
- Config schema and validation implemented.
- CLI scaffold (`init`, `validate`, `plan`) implemented.

## Phase 1: Runtime Baseline (Completed)

- Implemented local gateway runtime (`serve`).
- Implemented manifest rendering (`render`) and apply path (`apply`).
- Added local deploy stack with Envoy fronting gateway (`deploy/local/docker-compose.yml`).

## Phase 2: Envoy Gateway CRD Reconciler (In Progress)

- Replace placeholder MCP CR generation with Envoy Gateway CRD-native resources.
- Add state reconciliation loop and drift detection.
- Integrate Kubernetes API-based apply path as primary (kubectl fallback optional).

## Phase 3: Production Hardening

- Expand metrics/tracing coverage and dashboards.
- Add policy tests, failure injection tests, perf baseline.
- Add upgrade path and config migration checks.

## Definition of Done for v1.0

- Open-source repo with docs and runnable examples.
- One-command local run and Helm-based K8s deploy.
- Auth, TLS, observability, and policy controls functional.
