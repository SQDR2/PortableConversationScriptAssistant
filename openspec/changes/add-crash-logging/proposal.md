## Why

应用在 Windows 环境下会出现"自行关闭"的问题，但开发者没有 Windows 运行环境，无法现场调试。目前应用没有持久化的错误日志机制，所有 `runtime.LogErrorf` 输出只在 Wails 开发者工具中可见，一旦应用崩溃退出，错误信息随之丢失。需要一套文件日志系统来捕获和持久化运行时错误，便于远程排查崩溃问题。

## What Changes

- 新增日志工具模块，在应用运行目录下创建 `logs/` 目录
- 每次错误事件生成一个独立的日志文件，文件名包含时间戳和错误类型
- 日志文件内容包括：时间、错误信息、堆栈跟踪（stack trace）、运行上下文
- 在 `main()` 函数添加全局 panic 恢复器（defer recover），捕获未处理的 panic 并写入日志
- 在 `checkTarget()` goroutine 内添加 panic 恢复，防止轮询崩溃导致应用退出
- 在 Windows API 调用和启动阶段增强错误日志记录
- 日志模块全平台通用，不依赖 Windows 特定 API

## Capabilities

### New Capabilities
- `crash-logging`: 运行时错误日志系统，负责日志目录管理、错误文件写入、panic 恢复和堆栈跟踪捕获

### Modified Capabilities
- `window-follow`: 在窗口跟随轮询 (`checkTarget`) 中集成 panic 恢复和错误日志记录

## Impact

- **新增文件**: `backend/utils/logger.go` — 日志工具模块（全平台通用）
- **修改文件**: `main.go` — 添加全局 panic 恢复和日志初始化
- **修改文件**: `backend/services/window_service.go` — `checkTarget()` 中添加 panic 恢复
- **运行时目录**: 应用运行目录下新增 `logs/` 子目录
- **依赖**: 仅使用 Go 标准库（`os`, `runtime/debug`, `time`, `fmt`, `path/filepath`），无外部依赖
