# CICD2Jenkins

这是一个按生产化分层方式组织的 Go 博客系统后端骨架，当前已经实现：

- 用户登录
- JWT 鉴权
- 两种角色权限控制
- 文章 CRUD
- 文章内容落到 Elasticsearch

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
│   ├── domain                  # 领域模型
│   ├── repository              # 仓储接口与实现
│   ├── service                 # 业务服务
│   └── transport/httpapi       # HTTP 路由、处理器、中间件
├── docker-compose.yml          # 本地 ES
└── Makefile
```

## 快速开始

### 1. 启动 Elasticsearch

```bash
docker compose up -d elasticsearch
```

### 2. 准备环境变量

```bash
cp configs/local.env.example .env
```

如果你用的是 `zsh` / `bash`，可以先导入环境变量：

```bash
set -a
source .env
set +a
```

### 3. 启动服务

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
    "tags": ["go", "elasticsearch"],
    "published": true
  }'
```

## 常用命令

```bash
make run
make test
make fmt
make wire
make es-up
```

## 当前模块范围

这次先按你的要求实现了最核心的两块：

- 用户模块：当前仅包含登录和身份识别
- 文章模块：支持文章 CRUD，内容存储在 Elasticsearch

后续建议优先补充的模块：

1. 分类与标签模块
2. 评论模块
3. 文件上传模块
4. 操作日志与审计模块
5. 后台仪表盘模块
6. 配置中心与多环境部署模块

## 新需求设计说明

根据你的最新要求，已补充架构与存储拆分设计文档：`docs/architecture.md`。

当前应用装配采用 `google/wire` 生成依赖注入代码，组合根位于 `internal/app`。
