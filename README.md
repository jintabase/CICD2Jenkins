# CICD2Jenkins

这是一个 Go 项目初始化模板。

## 快速开始

```bash
go run .
```

## 推送到云端仓库（GitHub 示例）

1. 在 GitHub 创建一个空仓库（例如：`CICD2Jenkins`）。
2. 关联远程仓库：

```bash
git remote add origin <你的仓库URL>
```

3. 推送当前分支：

```bash
git push -u origin $(git branch --show-current)
```
