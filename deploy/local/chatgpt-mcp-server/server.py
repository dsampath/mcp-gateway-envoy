from mcp.server.fastmcp import FastMCP

mcp = FastMCP(
    "Local Docs MCP",
    instructions=(
        "Local development MCP server with search/fetch tools for ChatGPT connector testing."
    ),
    json_response=True,
    host="0.0.0.0",
    port=8000,
    streamable_http_path="/mcp",
)

DOCS = [
    {
        "id": "getting-started",
        "title": "Getting Started",
        "text": "Run docker compose for local gateway and connect ChatGPT via the tunneled /mcp endpoint.",
    },
    {
        "id": "deploy",
        "title": "Deploy",
        "text": "Render manifests with gateway render and apply with gateway apply.",
    },
    {
        "id": "debug",
        "title": "Debugging",
        "text": "Set GATEWAY_LOG_BODIES=true to inspect request and response payloads in gateway logs.",
    },
]


@mcp.tool()
def search(query: str) -> list[dict[str, str]]:
    """Search local documentation snippets by keyword."""
    q = query.strip().lower()
    if not q:
        return []

    matches: list[dict[str, str]] = []
    for doc in DOCS:
        hay = f"{doc['title']} {doc['text']}".lower()
        if q in hay:
            matches.append(
                {
                    "id": doc["id"],
                    "title": doc["title"],
                    "snippet": doc["text"][:160],
                }
            )
    return matches[:8]


@mcp.tool()
def fetch(id: str) -> dict[str, str]:
    """Fetch a full document by id."""
    for doc in DOCS:
        if doc["id"] == id:
            return {
                "id": doc["id"],
                "title": doc["title"],
                "text": doc["text"],
            }
    return {"error": f"document not found: {id}"}


if __name__ == "__main__":
    mcp.run(transport="streamable-http")
