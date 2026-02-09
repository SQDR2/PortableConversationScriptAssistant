## Why

在部分 Windows 系统上，sidekick 在运行一段时间后会自动退出。日志显示触发了 Go 的致命错误 `fatal error: too many callback functions`，导致进程直接崩溃退出。

同时，当前构建/运行环境中数据库初始化可能失败（`CGO_ENABLED=0` 时 `go-sqlite3` 无法工作），进而导致脚本/分类接口在运行期出现空指针错误，影响核心功能可用性。

## What Changes

- 修复 Windows 平台窗口跟随轮询中的 callback 泄漏：避免在高频轮询路径中反复调用 `syscall.NewCallback`，防止触发 `too many callback functions` 崩溃。
- 调整 Windows 平台窗口定位/尺寸同步实现：优先使用已缓存的窗口句柄（hwnd）进行 Win32 API 调用，减少按标题枚举窗口等高开销/高风险操作。
- 修复数据库初始化在 `CGO_ENABLED=0` 场景下不可用的问题：确保应用在 Windows 上以一致方式可用（开发/发布构建均能初始化 DB），避免后续服务调用出现 nil pointer。
- 补充启动阶段错误日志记录：当 DB 初始化失败或窗口跟随逻辑出现异常时，记录可诊断的错误日志，便于定位用户环境差异。

## Capabilities

### New Capabilities

（无）

### Modified Capabilities

- `window-follow`: 约束实现细节以提升稳定性——Windows 上禁止在轮询热路径中无限创建回调；窗口定位应以 hwnd 为主，避免频繁按标题枚举窗口导致资源耗尽。
- `script-management`: 运行时必须能稳定读取脚本数据；当 DB 初始化失败时不得导致脚本接口崩溃（应返回可诊断错误并记录日志）。
- `script-categories`: 同上；当 DB 初始化失败时不得导致分类接口崩溃（应返回可诊断错误并记录日志）。
- `crash-logging`: 强化启动阶段错误日志记录覆盖（DB 初始化失败、窗口跟随异常）。

## Impact

- 受影响代码：
  - `backend/utils/window_windows.go`：Win32 窗口枚举/定位相关函数（避免重复 `syscall.NewCallback`）
  - `backend/services/window_service.go`：轮询与对齐逻辑（使用缓存 hwnd、降低枚举频率/次数）
  - `backend/db/db.go` / DB 驱动选择与构建参数：确保 Windows 下 DB 可初始化
  - 相关服务（CategoryService/ScriptService）在 DB 不可用时的错误处理路径
- 受影响平台：以 Windows 为主；其他平台应保持行为不变。
- 风险：涉及窗口跟随与数据层初始化，需验证常见 DPI/窗口状态场景，以及开发/发布构建方式（CGO on/off）。
