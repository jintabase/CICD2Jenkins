## 当前重构结果（Gin + GORM）

### 1) 当前分层
- `model`：定义 `User`、`Article`、`Actor` 等核心实体，并承载 GORM 映射。
- `repo`：直接依赖 `model` 与 GORM，承载数据库 CRUD。
- `logic`：承载登录、鉴权、文章 CRUD、角色校验等业务逻辑。
- `service`：承载 Gin handler，请求解析与响应编排放在这一层。

### 2) 当前存储
- 默认数据库使用 SQLite，零配置即可本地启动。
- 通过环境变量可以切换到 MySQL，GORM 初始化逻辑已经预留。
- 用户与文章都走统一的 GORM Repository，不再依赖 Elasticsearch 或内存仓储。

### 3) 当前启动流程
1. 加载配置与种子用户。
2. 初始化 GORM 数据库连接。
3. 自动迁移 `users`、`articles` 表结构。
4. 写入或更新默认管理员/读者账号。
5. 组装 Gin 路由、鉴权中间件和 HTTP Server。

## 依赖注入（DI）实现说明

本项目现在采用 **Google Wire 编译期依赖注入**，并使用 `internal/app` 作为组合根（Composition Root）。

### 具体做法
1. 在 `internal/app/providers.go` 中声明基础设施 provider，例如 GORM DB、Seed Users、Repository、Logic、HTTP Server。
2. 在 `internal/app/wire.go` 中通过 `wire.NewSet(...)` 描述仓储、逻辑、服务、HTTP 路由的装配关系。
3. 使用 `wire` 生成 `internal/app/wire_gen.go`，由生成代码负责把 Repository -> Logic -> Service -> Router -> Server 串起来。
4. `internal/app/server.go` 只保留稳定的 `NewServer(cfg)` 入口，对外隐藏具体注入细节。

### 优点
- 依赖关系显式、可读性高。
- 单元测试容易替换依赖。
- 装配代码由编译期生成，减少手写样板代码，同时保持类型安全。
- 当后续接入对象存储、搜索引擎、多仓储、日志链路时，可以继续把新 provider 纳入同一个 Wire Set 统一管理。

## 仓储层现在的约定

当前不再把仓储层拆成 `internal/repo` 接口 + `internal/repo/gormrepo` 实现两部分，而是直接保留一个 `internal/repo`：

1. `repo` 层直接依赖 `model` 的表结构与 GORM，负责增删改查。
2. `logic` 层依赖这个仓储层，处理鉴权、参数校验、角色判断、业务规则。
3. `service` 层只负责 HTTP 协议适配，对外暴露接口。

这样目录会更直观，也更符合这个项目当前以 MySQL/GORM 为主的数据访问方式。
