## 1. 项目初始化 (Project Setup)

- [x] 1.1 使用 Wails CLI 初始化 Go + Vue + Quasar 项目结构 (`wails init`).
- [x] 1.2 配置项目基础依赖 (Go mod, npm packages).
- [x] 1.3 搭建基础的目录结构 (backend, frontend).

## 2. 后端核心开发 - 窗口跟随 (Backend - Window Tracker)

- [x] 2.1 实现 Go 端的 Window Utils (调用 User32.dll/WinAPI 获取窗口坐标).
- [x] 2.2 实现 Linux 下的 Window Utils (使用 XLib/EWMH 获取窗口坐标).
- [x] 2.3 开发 `WindowService`，提供 `GetStartApps` 接口供前端选择目标.
- [x] 2.4 实现后台轮询机制 (Ticker)，定期向前端发送目标窗口的位置更新事件.

## 3. 后端核心开发 - 话术管理 (Backend - Script Management)

- [x] 3.1 集成 GORM 和 SQLite.
- [x] 3.2 定义 `Script` 数据模型 (ID, Content, Tags, etc.).
- [x] 3.3 实现 `ScriptService` 的 CRUD 方法 (Create, Update, Delete, List).
- [x] 3.4 实现文件导入逻辑 (解析 txt 文本并批量存入 DB).
- [x] 3.5 实现基于 SQLite FTS 的模糊搜索功能.

## 4. 前端界面开发 - 基础框架 (Frontend - UI Foundation)

- [x] 4.1 初始化 Quasar 布局配置 (Layout, Theme).
- [x] 4.2 实现 Glassmorphism 全局样式 (Backdrop filter).
- [x] 4.3 开发 "选择目标应用" 的弹窗界面.
- [x] 4.4 对接 Wails Runtime，实现窗口自身的位置控制 (SetPosition, Hide, Show).

## 5. 前端界面开发 - 话术功能 (Frontend - Script Features)

- [x] 5.1 开发话术列表组件 (Virtual Scroll).
- [x] 5.2 实现搜索栏组件 (Search Input).
- [x] 5.3 实现话术的增删改交互 (CRUD Dialogs).
- [x] 5.4 开发话术导入预览界面 (File Parse & Preview).

## 6. 整合与优化 (Integration & Polish)

- [x] 6.1 窗口跟随的偏移量校准 (Offset Calibration).
- [x] 6.2 性能优化 (CPU/Memory usage check).
- [x] 6.3 Windows/Linux 跨平台构建测试.
