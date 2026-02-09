## 1. 日志工具模块

- [x] 1.1 创建 `backend/utils/logger.go`，实现日志目录初始化函数 `InitLogDir()`，在运行目录下创建 `logs/` 子目录，失败时输出到 `stderr` 但不阻止启动
- [x] 1.2 实现错误日志写入函数 `LogError(errorType string, err interface{})`，生成独立日志文件 `{YYYY-MM-DD_HH-mm-ss}_{errorType}.log`，内容包含时间戳、错误类型、错误信息和 `runtime/debug.Stack()` 堆栈跟踪

## 2. main 函数集成

- [x] 2.1 在 `main.go` 的 `main()` 函数最开头调用 `utils.InitLogDir()` 初始化日志目录
- [x] 2.2 在 `main()` 函数中添加 `defer recover()` 全局 panic 恢复，捕获 panic 后调用 `utils.LogError("panic", r)` 写入日志
- [x] 2.3 在 `OnStartup` 回调中 `db.InitDB` 失败时调用 `utils.LogError("startup_db", err)` 记录错误

## 3. window_service 集成

- [x] 3.1 在 `backend/services/window_service.go` 的 `checkTarget()` 方法开头添加 `defer recover()` panic 恢复，捕获后调用 `utils.LogError("check_target_panic", r)`，不终止轮询

## 4. 验证

- [x] 4.1 在 Linux/macOS 下测试日志目录创建和日志文件写入
- [x] 4.2 验证模拟 panic 场景下日志文件正确生成且应用不崩溃（`checkTarget` 场景）
