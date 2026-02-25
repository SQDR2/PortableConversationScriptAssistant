# Tech Stack

## Overall Architecture

- 类型：Desktop App（Wails）
- 模式：Go 后端 + Vue 3/TypeScript 前端 + 本地 SQLite

## Backend

- 语言：Go
- 框架：Wails v2
- ORM：GORM
- 数据库：SQLite（含 FTS5）
- 关键能力：窗口联动（Linux XGB/EWMH、Windows Win32）、翻译服务集成（腾讯云 SDK）

## Frontend

- 框架：Vue 3（Composition API）
- 语言：TypeScript
- UI：Quasar Framework
- 构建：Vite
- 测试：Vitest（按项目约定）

## Tooling

- 包管理：pnpm（frontend）
- 版本控制：Git
- 运行方式：`wails dev` / `wails build`
