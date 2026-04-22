# Introduction

A lightweight MCP server written in Go for interacting with MySQL databases.

---

## Features

- Built with Go

- Connect to MySQL using standard DSN

- Expose database operations through MCP

- Lightweight and easy to deploy

- Suitable for AI tools, automation workflows, and internal integrations

- Easy to extend with custom tools and permission controls

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
docker run -d --name mysql-mcp-server betterde/mysql-mcp-server:latest
```

### Build from source

```bash
git clone https://github.com/betterde/mysql-mcp-server.git

cd mysql-mcp-server

go build -o mysql-mcp-server .
```