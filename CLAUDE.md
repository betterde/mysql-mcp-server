# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & run

```bash
go build -o mysql-mcp-server .
go vet ./...
./mysql-mcp-server serve        # start the MCP streamable HTTP server
./mysql-mcp-server serve -v      # force DEBUG log level
./mysql-mcp-server --config .config.yaml serve
```

There is no test suite or Makefile in this repository.

## Architecture

```
main.go                     → calls cmd.Execute()
cmd/
  root.go                   → Cobra root command, Viper config init
  serve.go                  → "serve" subcommand: creates MCP server, registers tools, starts HTTP
config/config.go            → Viper-based Config struct (DSN, HTTP, Logging, ReadOnly)
global/context.go           → shared context.Context + CancelFunc
internal/
  build/info.go             → ldflags-injected build vars (Name, Version, Commit, Build)
  journal/logger.go         → uber-go/zap logger setup, JSON output to stdout
  middleware/logging.go     → MCP receiving middleware that logs method calls and timing
  mysql/driver.go           → single sql.Conn from DSN, initiated once at startup
tools/
  register.go               → calls Register() for each tool category
  table/                    → tables/list, tables/schema, tables/columns
  query/                    → query/select (read-only enforcement, limit clamping), query/explain
  admin/                    → admin/server/health, process list, connections, version, threads
  database/                 → database/list
```

## Key patterns

- **MCP SDK:** Uses `github.com/modelcontextprotocol/go-sdk`. Servers are `mcp.Server`, tools are registered with `mcp.AddTool(server, toolDef, handler)`. The streamable HTTP handler is created with `mcp.NewStreamableHTTPHandler`.
- **Handler signature:** Every MCP tool handler follows the same triple-return pattern: `func handler(ctx, *mcp.CallToolRequest, input InputType) (*mcp.CallToolResult, OutputType, error)`.
- **Tool naming:** Forward-slash namespaced names, e.g. `query/select`, `admin/server/health`, `tables/list`.
- **Read-only query enforcement:** `query/select` and `query/explain` only accept queries starting with `SELECT`, `SHOW`, `DESCRIBE`, `DESC`, or `WITH`. Multi-statement queries (containing `;`) are rejected.
- **Config:** Viper loads from `.config.yaml` → `.env` (via `subosito/gotenv`) → env vars prefixed `MYSQL_MCP_SERVER_` (dots replaced with underscores). Defaults: listen on `0.0.0.0:8080`, log level `DEBUG`, read_only `true`.
- **Logging:** Structured JSON via `zap`. The logger is created in `journal.InitLogger()` and reused globally as `journal.Logger`.
- **Startup flow:** `main()` → `cmd.Execute()` → Cobra kicks off `initConfig()` (logger + config) → `serve` runs `start()` which initializes the MySQL connection, registers tools, and starts the HTTP server.

## Docker

Docker CI triggers on version tags (`v*`). Multi-arch build (`linux/amd64`, `linux/arm64`) pushes to `ghcr.io`. Build args inject version/build/commit for ldflags.
