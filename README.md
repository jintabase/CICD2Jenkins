# CICD2Jenkins

这是一个按生产化分层方式组织的 Go 博客系统后端骨架，现已重构为 `gin + gorm` 技术栈，并实现：

- 用户登录
- JWT 鉴权
- 两种角色权限控制
- 文章 CRUD
- 用户与文章数据持久化到数据库

当前角色说明：

- `SUPER_ADMIN`：可以登录并执行文章新增、修改、删除、查询
- `USER`：可以登录并浏览文章，不能执行新增、修改、删除

## 项目结构

```text
.
├── cmd/blog-api                # 启动入口
├── configs                     # 本地环境变量示例
├── internal
│   ├── app                     # 应用装配
│   ├── apperrors               # 业务错误定义
│   ├── config                  # 配置加载
│   ├── logic                   # 业务逻辑层
│   ├── model                   # GORM 模型与核心实体
│   ├── repo                    # 仓储层，直接承载 GORM CRUD
│   ├── service                 # Gin 接口服务层
│   └── transport/httpapi       # HTTP 路由、中间件、响应工具
├── docker-compose.yml          # 可选本地 MySQL
└── Makefile
```

## 快速开始

### 1. 准备环境变量

```bash
cp configs/local.env.example .env
```

如果你用的是 `zsh` / `bash`，可以先导入环境变量：

```bash
set -a
source .env
set +a
```

默认使用 SQLite，本地直接运行就能启动，不需要额外数据库。

如果你想切换成 MySQL，可以执行：

```bash
docker compose up -d mysql
```

并把 `.env` 中的 `DB_DRIVER` / `DB_DSN` 改成 MySQL 配置。

### 2. 启动服务

```bash
go run ./cmd/blog-api
```

或使用：

```bash
make run
```

服务默认监听：

```text
http://localhost:8080
```

## 默认账号

默认内置了两个种子账号，便于本地开发联调：

| 用户名 | 密码 | 角色 |
| --- | --- | --- |
| `admin` | `Admin@123456` | `SUPER_ADMIN` |
| `reader` | `Reader@123456` | `USER` |

生产环境请务必通过环境变量覆盖默认密码。

## 接口示例

### 登录

```bash
curl --request POST 'http://localhost:8080/api/v1/auth/login' \
  --header 'Content-Type: application/json' \
  --data '{
    "username": "admin",
    "password": "Admin@123456"
  }'
```

### 获取当前用户

```bash
curl --request GET 'http://localhost:8080/api/v1/me' \
  --header 'Authorization: Bearer <token>'
```

### 查询文章列表

```bash
curl --request GET 'http://localhost:8080/api/v1/articles' \
  --header 'Authorization: Bearer <token>'
```

### 新增文章

```bash
curl --request POST 'http://localhost:8080/api/v1/articles' \
  --header 'Authorization: Bearer <token>' \
  --header 'Content-Type: application/json' \
  --data '{
    "title": "第一篇博客",
    "summary": "这是文章摘要",
    "content": "这是文章正文",
    "tags": ["go", "gorm"],
    "published": true
  }'
```

## 常用命令

```bash
make run
make test
make fmt
make wire
make mysql-up
```

## 当前技术栈说明

- HTTP 路由框架：`gin-gonic/gin`
- ORM：`gorm`
- 默认数据库：`SQLite`
- 可选数据库：`MySQL`

## 新需求设计说明

当前应用装配采用 `google/wire` 生成依赖注入代码，组合根位于 `internal/app`。重构后的架构说明见 `docs/architecture.md`。
