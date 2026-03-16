## 存储与模块设计（按你的最新要求）

### 1) 图片上传模块（MinIO）
- 图片二进制文件存储到 MinIO（对象存储）。
- 文章、评论、用户头像等只在业务库中保存对象 key / URL，不保存二进制。

### 2) 文章模块（ES + MySQL 混合存储）
- **Elasticsearch**：只存文章正文内容（`content`）与搜索相关字段（可选分词字段）。
- **MySQL**：存文章元数据（标题、摘要、作者、创建时间、更新时间、发布状态、分类、标签关系等）。

### 3) 分类/标签模块（MySQL）
- 分类表、标签表及文章标签关联表都放 MySQL。

### 4) 评论模块（MySQL）
- 评论主数据（文章ID、用户ID、内容、层级、状态、创建时间）放 MySQL。

### 5) 日志模块（ES）
- 操作日志/审计日志写入 Elasticsearch，便于检索与分析。

---

## 依赖注入（DI）实现说明

本项目现在采用 **Google Wire 编译期依赖注入**，并使用 `internal/app` 作为组合根（Composition Root）。

### 具体做法
1. 在 `internal/app/providers.go` 中声明基础设施 provider，例如 ES client、Article Repository、Seed Users、Auth Service、HTTP Server。
2. 在 `internal/app/wire.go` 中通过 `wire.NewSet(...)` 描述仓储、服务、HTTP 路由的装配关系，并对接口做 `wire.Bind(...)` 绑定。
3. 使用 `wire` 生成 `internal/app/wire_gen.go`，由生成代码负责把 Repository -> Service -> Router -> Server 串起来。
4. `internal/app/server.go` 只保留稳定的 `NewServer(cfg)` 入口，对外隐藏具体注入细节。

### 优点
- 依赖关系显式、可读性高。
- 单元测试容易替换依赖（你现在的 `*_test.go` 就是通过 stub repo 来完成）。
- 装配代码由编译期生成，减少手写样板代码，同时保持类型安全。
- 当后续接入 MinIO、MySQL、多仓储、日志链路时，可以继续把新 provider 纳入同一个 Wire Set 统一管理。
