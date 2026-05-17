# MySQL MCP Server

A lightweight MCP server written in Go for interacting with MySQL databases.

---

## Features

- Built with Go

- Connect to MySQL using standard DSN

- Expose database operations through MCP (Model Context Protocol)

- Lightweight and easy to deploy

- Suitable for AI tools, automation workflows, and internal integrations

- Easy to extend with custom tools and permission controls

- **v2 new**: Schema and metadata operations split into MCP Resources with URI-based access

---

## Use Cases

- Allow AI assistants to query MySQL in a controlled way

- Build internal database operation tools

- Expose schema inspection and query capabilities to automation systems

- Simplify integration between MySQL and MCP-compatible clients

---

## Installation

### Go install

```bash
go install github.com/betterde/mysql-mcp-server@latest
```

### Docker

```bash
docker run -d --name mysql-mcp-server ghcr.io/betterde/mysql-mcp-server:latest
```

### Build from source

```bash
git clone https://github.com/betterde/mysql-mcp-server.git

cd mysql-mcp-server

go build -o mysql-mcp-server .
```

---

## Configuration

Config is loaded from `.config.yaml`, with `.env` and environment variable overrides.

### Config file

```yaml
dsn: user:pass@tcp(127.0.0.1:3306)
http:
  listen: 0.0.0.0:8080
logging:
  level: DEBUG
read_only: true
```

### Environment variables

Environment variables use the `MYSQL_MCP_SERVER_` prefix, with dots replaced by underscores:

```bash
export MYSQL_MCP_SERVER_DSN="user:pass@tcp(127.0.0.1:3306)"
export MYSQL_MCP_SERVER_READ_ONLY=true
export MYSQL_MCP_SERVER_LOGGING_LEVEL=DEBUG
```

### Config options

| Option | Env variable | Default | Description |
|--------|-------------|---------|-------------|
| `dsn` | `MYSQL_MCP_SERVER_DSN` | — | MySQL connection string |
| `http.listen` | `MYSQL_MCP_SERVER_HTTP_LISTEN` | `0.0.0.0:8080` | HTTP server listen address |
| `logging.level` | `MYSQL_MCP_SERVER_LOGGING_LEVEL` | `DEBUG` | Log level: DEBUG/INFO/WARN/ERROR |
| `read_only` | `MYSQL_MCP_SERVER_READ_ONLY` | `true` | Read-only mode; when `false`, write tool is exposed |

---

## Starting the server

```bash
# Default config
./mysql-mcp-server serve

# Custom config file
./mysql-mcp-server --config .config.yaml serve

# Force DEBUG log level
./mysql-mcp-server serve -v
```

---

## MCP Resources (v2)

Access database metadata and server information via URIs:

| URI | Description |
|-----|-------------|
| `mysql://databases` | List all databases |
| `mysql://server/version` | Server version info |
| `mysql://server/variables` | Global system variables |
| `mysql://server/status` | Global status variables |
| `mysql://{database}/tables` | List tables in a database |
| `mysql://{database}/{table}/columns` | Column definitions for a table |
| `mysql://{database}/{table}/ddl` | Table DDL (CREATE TABLE statement) |

---

## MCP Tools

| Tool name | Description |
|-----------|-------------|
| `query/select` | Execute a read-only query (SELECT/SHOW/DESCRIBE/DESC/WITH) |
| `query/explain` | Show query execution plan |
| `query/execute` | Execute write operations (INSERT/UPDATE/DELETE/DDL), only available when `read_only: false` |
| `admin/health` | Server health check and key runtime metrics |
| `admin/processes` | Current process list |
| `admin/connections` | Current connection details |
| `admin/threads` | Thread state information |

---

## Running tests

```bash
go test ./...
go test -v ./tools/query/
go test -v ./resources/
go test -v ./internal/handler/
```

## Tech stack

- **Go** 1.26+
- **MySQL Driver**: `github.com/go-sql-driver/mysql`
- **MCP SDK**: `github.com/modelcontextprotocol/go-sdk`
- **CLI**: `github.com/spf13/cobra` + `github.com/spf13/viper`
- **Logging**: `go.uber.org/zap`

## License

MIT License
