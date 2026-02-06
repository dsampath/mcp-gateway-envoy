# Build Roadmap

Date: 2026-02-06

## Phase 0: Foundation (1 week)

- Confirm architecture (Envoy Gateway CRD-driven control path).
- Define config schema for server registry and route policy.
- Set up CI skeleton (lint, unit tests, image build).

Exit criteria:
- Architecture doc approved.
- Config schema validated with example files.

## Phase 1: MVP Runtime (2-3 weeks)

- Implement gateway controller that translates config to Envoy Gateway resources.
- Implement static registry loader and config validator.
- Add API key/JWT auth policy wiring.
- Add stdio adapter launcher contract.

Exit criteria:
- End-to-end request from MCP client through Envoy to sample server.
- Local compose startup under 10 minutes from clean machine.

## Phase 2: Deploy in Minutes UX (1-2 weeks)

- Ship `gateway init`, `gateway validate`, `gateway up` CLI commands.
- Publish Docker Compose + Helm quickstarts.
- Add first-run diagnostics and troubleshooting guide.

Exit criteria:
- New user can run quickstart and pass health checks in <=10 minutes.

## Phase 3: Production Hardening (2 weeks)

- Expand metrics/tracing coverage and dashboards.
- Add policy tests, failure injection tests, perf baseline.
- Add upgrade path and config migration checks.

Exit criteria:
- Reliability and performance targets met in staging benchmarks.

## MVP Definition of Done

- Open-source repo with docs and runnable examples.
- One-command local run and Helm-based K8s deploy.
- Auth, TLS, observability, and basic policy controls functional.
