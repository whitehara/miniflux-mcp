# Miniflux MCP Server

A Model Context Protocol (MCP) server for interacting with [Miniflux](https://miniflux.app/) RSS reader.
This is a fork of [tssujt/miniflux-mcp](https://github.com/tssujt/miniflux-mcp) that adds **Streamable HTTP transport** support, making it usable as a remote MCP server behind `sigbit/mcp-auth-proxy`.

## Features

- **Feed Management**: List, create, refresh, and delete RSS/Atom feeds
- **Entry Operations**: Read entries, update status (read/unread/removed), bookmark, save
- **Category Management**: List, create, update, and delete categories
- **User Management**: List, create, and delete users
- **Flexible Authentication**: API key or username/password
- **Dual Transport**: stdio (local) or Streamable HTTP (remote) via `MCP_HTTP_PORT`

## Environment Variables

| Variable | Description | Required |
|---|---|---|
| `MINIFLUX_URL` | Your Miniflux instance URL | Yes |
| `MINIFLUX_API_KEY` | API key for authentication | Yes* |
| `MINIFLUX_USERNAME` | Username for basic auth | Yes* |
| `MINIFLUX_PASSWORD` | Password for basic auth | Yes* |
| `MCP_HTTP_PORT` | Port to listen on in HTTP mode | No |

*Either `MINIFLUX_API_KEY` or both `MINIFLUX_USERNAME` and `MINIFLUX_PASSWORD` must be set.

## Getting a Miniflux API Key

1. Log into your Miniflux instance
2. Go to **Settings → API Keys**
3. Create a new API key and copy the generated value

## Usage

### stdio mode (local / Claude Desktop)

The default mode. Suitable for local use with Claude Desktop or any stdio-based MCP client.

```bash
docker run -i --rm \
  -e MINIFLUX_URL=https://miniflux.example.com \
  -e MINIFLUX_API_KEY=your_api_key_here \
  ghcr.io/whitehara/miniflux-mcp:latest
```

To integrate with **Claude Desktop**, add the following to your configuration file:

```json
{
  "mcpServers": {
    "miniflux": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-e", "MINIFLUX_URL",
        "-e", "MINIFLUX_API_KEY",
        "ghcr.io/whitehara/miniflux-mcp:latest"
      ],
      "env": {
        "MINIFLUX_URL": "https://your-miniflux-instance.com",
        "MINIFLUX_API_KEY": "your_api_key_here"
      }
    }
  }
}
```

### HTTP mode (remote / mcp-auth-proxy)

Set `MCP_HTTP_PORT` to switch to Streamable HTTP transport. The MCP endpoint is exposed at `/mcp`.

```bash
docker run -d \
  -e MINIFLUX_URL=https://miniflux.example.com \
  -e MINIFLUX_API_KEY=your_api_key_here \
  -e MCP_HTTP_PORT=3000 \
  -p 3000:3000 \
  ghcr.io/whitehara/miniflux-mcp:latest
```

MCP endpoint: `http://<host>:3000/mcp`

This mode is designed to work behind [`sigbit/mcp-auth-proxy`](https://github.com/sigbit/mcp-auth-proxy) for OIDC authentication (e.g., Cloudflare Access). The proxy forwards requests to port 3000 on the `/mcp` path.

Example stack with mcp-auth-proxy:

```yaml
services:
  miniflux-mcp:
    image: ghcr.io/whitehara/miniflux-mcp:latest
    environment:
      MINIFLUX_URL: https://miniflux.example.com
      MINIFLUX_API_KEY: your_api_key_here
      MCP_HTTP_PORT: "3000"

  mcp-auth-proxy:
    image: ghcr.io/sigbit/mcp-auth-proxy:latest
    environment:
      TARGET_URL: http://miniflux-mcp:3000/mcp
      # ... OIDC config
    ports:
      - "8080:8080"
```

## Available Tools

The server provides **40+ tools** covering all Miniflux API functionality. See the [Miniflux API Reference](https://miniflux.app/docs/api.html#go-client) for details.

### Feed Management (10 tools)
- `get_feeds` — Get all RSS/Atom feeds
- `get_feed` — Get a specific feed by ID
- `create_feed` — Add a new RSS/Atom feed
- `delete_feed` — Delete a specific feed
- `refresh_feed` — Manually refresh a specific feed
- `refresh_all_feeds` — Refresh all feeds
- `get_feed_entries` — Get entries from a specific feed
- `get_feed_entry` — Get a specific entry from a feed
- `get_feed_icon` — Get the icon of a specific feed
- `mark_feed_as_read` — Mark all entries in a feed as read

### Entry Management (8 tools)
- `get_entries` — Get entries with optional filtering
- `get_entry` — Get a specific entry by ID
- `update_entry_status` — Update entry status (read/unread/removed)
- `toggle_bookmark` — Toggle bookmark status of an entry
- `save_entry` — Save an entry
- `fetch_original_content` — Fetch original content of an entry
- `mark_all_as_read` — Mark all entries as read for a user
- `get_category_entry` — Get a specific entry from a category

### Category Management (8 tools)
- `get_categories` — Get all feed categories
- `create_category` — Create a new category
- `update_category` — Update a category title
- `delete_category` — Delete a category
- `get_category_feeds` — Get all feeds in a specific category
- `get_category_entries` — Get all entries in a specific category
- `mark_category_as_read` — Mark all entries in a category as read
- `refresh_category` — Refresh all feeds in a category

### User Management (6 tools)
- `get_users` — Get all users
- `get_me` — Get current user information
- `get_user_by_id` — Get a specific user by ID
- `get_user_by_username` — Get a specific user by username
- `create_user` — Create a new user
- `delete_user` — Delete a user

### System & Utility (7 tools)
- `get_version` — Get Miniflux version information
- `healthcheck` — Perform a health check
- `fetch_counters` — Fetch feed counters
- `discover` — Discover feeds from a URL
- `export` — Export feeds as OPML
- `flush_history` — Flush the read history

### API Key Management (3 tools)
- `get_api_keys` — Get all API keys
- `create_api_key` — Create a new API key
- `delete_api_key` — Delete an API key

### Icons & Media (2 tools)
- `get_icon` — Get an icon by ID
- `get_enclosure` — Get an enclosure by ID

## License

MIT License — see [LICENSE](LICENSE) for details.
