# MySQL MCP Server

一个基于 Go 语言开发的轻量级 MCP 服务器，用于与 MySQL 数据库进行交互。

---

## 功能特性

- 使用 Go 语言构建

- 通过标准 DSN 连接 MySQL 数据库

- 通过 MCP（Model Context Protocol）暴露数据库操作能力

- 轻量级、易于部署

- 适用于 AI 工具、自动化工作流和内部集成场景

- 易于通过自定义 Tools 和权限控制进行扩展

- **v2 新增**：将 schema 和元数据操作拆分为 MCP Resources，支持通过 URI 访问

---

## 使用场景

- 允许 AI 助手以可控的方式查询 MySQL 数据库

- 构建内部数据库运维工具

- 向自动化系统暴露 schema 检查和查询能力

- 简化 MySQL 与 MCP 兼容客户端之间的集成

---

## 安装方式

### Go 安装

```bash
go install github.com/betterde/mysql-mcp-server@latest
```

### Docker

```bash
docker run -d --name mysql-mcp-server ghcr.io/betterde/mysql-mcp-server:latest
```

### 从源码构建

```bash
git clone https://github.com/betterde/mysql-mcp-server.git

cd mysql-mcp-server

go build -o mysql-mcp-server .
```

---

## 配置说明

通过 `.config.yaml` 文件配置，支持 `.env` 和环境变量覆盖。

### 配置文件

```yaml
dsn: user:pass@tcp(127.0.0.1:3306)
http:
  listen: 0.0.0.0:8080
logging:
  level: DEBUG
read_only: true
```

### 环境变量

环境变量使用 `MYSQL_MCP_SERVER_` 前缀，点号分隔转为下划线：

```bash
export MYSQL_MCP_SERVER_DSN="user:pass@tcp(127.0.0.1:3306)"
export MYSQL_MCP_SERVER_READ_ONLY=true
export MYSQL_MCP_SERVER_LOGGING_LEVEL=DEBUG
```

### 配置选项

| 选项 | 环境变量 | 默认值 | 说明 |
|------|---------|--------|------|
| `dsn` | `MYSQL_MCP_SERVER_DSN` | — | MySQL 连接字符串 |
| `http.listen` | `MYSQL_MCP_SERVER_HTTP_LISTEN` | `0.0.0.0:8080` | HTTP 服务监听地址 |
| `logging.level` | `MYSQL_MCP_SERVER_LOGGING_LEVEL` | `DEBUG` | 日志级别：DEBUG/INFO/WARN/ERROR |
| `read_only` | `MYSQL_MCP_SERVER_READ_ONLY` | `true` | 只读模式，为 `false` 时开放写操作 tool |

---

## 启动服务

```bash
# 默认配置
./mysql-mcp-server serve

# 指定配置文件
./mysql-mcp-server --config .config.yaml serve

# 强制 DEBUG 日志级别
./mysql-mcp-server serve -v
```

---

## MCP Resources (v2)

通过 URI 访问数据库元数据和服务器信息：

| URI | 说明 |
|-----|------|
| `mysql://databases` | 列出所有数据库 |
| `mysql://server/version` | 服务器版本信息 |
| `mysql://server/variables` | 全局系统变量 |
| `mysql://server/status` | 全局状态指标 |
| `mysql://{database}/tables` | 列出数据库中的表 |
| `mysql://{database}/{table}/columns` | 表的列定义 |
| `mysql://{database}/{table}/ddl` | 表的 DDL（建表语句） |

---

## MCP Tools

| Tool 名称 | 说明 |
|-----------|------|
| `query/select` | 执行只读查询（SELECT/SHOW/DESCRIBE/DESC/WITH） |
| `query/explain` | 显示查询执行计划 |
| `query/execute` | 执行写操作（INSERT/UPDATE/DELETE/DDL），仅在 `read_only: false` 时可用 |
| `admin/health` | 服务器健康检查和关键运行时指标 |
| `admin/processes` | 当前进程列表 |
| `admin/connections` | 当前连接详情 |
| `admin/threads` | 线程状态信息 |

---

## 运行测试

```bash
go test ./...
go test -v ./tools/query/
go test -v ./resources/
go test -v ./internal/handler/
```

## 技术栈

- **Go** 1.26+
- **MySQL Driver**: `github.com/go-sql-driver/mysql`
- **MCP SDK**: `github.com/modelcontextprotocol/go-sdk`
- **CLI**: `github.com/spf13/cobra` + `github.com/spf13/viper`
- **Logging**: `go.uber.org/zap`

## 项目许可证

MIT License
