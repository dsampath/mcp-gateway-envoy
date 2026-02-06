# ChatGPT Connection Setup (Local MCP via Gateway)

Date: 2026-02-06

## Important Current Limitation

As of February 6, 2026, OpenAI's MCP apps/connectors documentation states custom MCP apps are available on ChatGPT web and are not available on mobile. In practice this means you should do connector setup and validation in ChatGPT web first.

- OpenAI help article (MCP apps in ChatGPT): https://help.openai.com/en/articles/11487775-connectors-in-chatgpt
- OpenAI product blog (connectors + custom MCP): https://openai.com/index/new-tools-and-features-in-the-responses-api/

If your macOS client does not show the MCP connector UI yet, this is expected for now.

## Local Stack

This repo includes a local MCP server with `search` and `fetch` tools, routed through your gateway and Envoy.

Run:

```bash
cd /Users/djsam/codex/mcp-gateway-envoy
docker compose -f deploy/local/chatgpt-compose.yml up --build
```

Endpoints:

- MCP through Envoy: `http://localhost:10000/mcp`
- Gateway health: `http://localhost:18080/healthz`
- Raw local MCP server (debug): `http://localhost:18110/mcp`

## Make It Reachable by ChatGPT (Remote URL Required)

ChatGPT connectors need a publicly reachable URL. For local testing, use a tunnel.

Example with Cloudflare tunnel:

```bash
cloudflared tunnel --url http://localhost:10000
```

Use the generated `https://<random>.trycloudflare.com/mcp` URL in ChatGPT connector setup.

## Configure in ChatGPT (Web)

1. Open ChatGPT on the web.
2. Add a custom MCP connector/app.
3. Server URL: `https://<your-tunnel-host>/mcp`
4. Authentication: none (for local testing).
5. Save and test.

## Observe Prompts and Responses Through Gateway

Gateway request/response payload logging is enabled in `deploy/local/chatgpt-compose.yml` via:

- `GATEWAY_LOG_BODIES=true`

Tail logs:

```bash
docker compose -f deploy/local/chatgpt-compose.yml logs -f gateway
```

You will see lines like:

- `mcp_request ... body=...`
- `mcp_response ... body=...`

## Security Note

This local profile runs without auth for fast iteration. Before any real deployment, switch to authenticated routes and TLS-only public exposure.
