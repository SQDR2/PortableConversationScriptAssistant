## Context

当前 Windows 平台窗口跟随由 `WindowService.StartPolling()` 以 50ms 周期轮询 `checkTarget()` 实现（20Hz）。在 v2.0.1 之后，为了在 Windows 上实现“视觉高度精确对齐”，实现改为绕过 Wails 的窗口管理 API，直接使用 Win32 `SetWindowPos`（参见 `window-follow` 规范）。

现状实现中，多处通过 `EnumWindows` 按标题查找窗口句柄，并在每次查找时创建新的 `syscall.NewCallback(...)`。由于 `checkTarget()` 高频调用这些函数，导致 Windows 进程不断创建 callback thunk，最终触发 Go runtime 的致命错误：`fatal error: too many callback functions`，进程直接退出，表现为“运行一会自动关闭”。

同时，数据库层使用 GORM 的 sqlite driver（底层依赖 `go-sqlite3`）。当构建参数为 `CGO_ENABLED=0` 时（常见于某些 CI/交叉编译或用户本地环境），sqlite driver 变为 stub，`db.InitDB` 失败，后续 Script/Category 相关服务仍启动并访问 `db.DB`，引发空指针错误。

## Goals / Non-Goals

**Goals:**

- 修复 Windows 平台运行一段时间后崩溃退出的问题，确保窗口跟随长期稳定运行。
- 保持 `window-follow` 的像素级对齐效果与现有交互不变（仍使用 Win32 `SetWindowPos`）。
- 让数据库初始化在 Windows 下“默认可用”：开发运行与发布构建不因 `CGO_ENABLED` 差异而失效。
- 当 DB 初始化失败时，应用应输出清晰、可诊断的错误并记录日志，服务层不得因 nil pointer 崩溃。

**Non-Goals:**

- 不重写窗口跟随架构（仍保持轮询机制与现有事件模型）。
- 不引入新的数据库类型或迁移到完全不同的数据存储（除非为了兼容 `CGO_ENABLED=0` 的 sqlite 实现替换为纯 Go 方案）。
- 不对前端 UI/交互做额外扩展。

## Decisions

### 决策 1：Windows 轮询热路径禁止创建新的 callback

**选择**：将 `EnumWindows` 的 callback 设计为“进程生命周期内复用”，或在轮询路径中完全避免 `EnumWindows`。

**理由**：`syscall.NewCallback` 会生成不可回收/数量受限的回调桩。在 20Hz 轮询下，任何按次创建都会在短时间内耗尽配额并触发 fatal。

**备选方案**：

- 继续按标题枚举窗口但降低频率（例如 500ms/1s）→ 只能延后崩溃，无法根治。
- 在每次调用后“释放” callback → Go 的 `NewCallback` 并无对等释放机制，无法实现。

### 决策 2：窗口定位/取宽高尽量使用缓存 hwnd，而不是按标题查找

**选择**：在 `WindowService` 内维护 sidekick 自身 hwnd（已存在 `sidekickHWID` 字段），并补充将其转换为 `syscall.Handle` / `uintptr` 后直接用于 Win32 API。

**理由**：

- sidekick 自身窗口句柄在进程内稳定，适合缓存。
- 避免在轮询中调用 `EnumWindows`，既减少 callback 风险，也提升性能。

**边界处理**：如果 hwnd 尚未获取到或失效（例如窗口重建），允许低频（例如启动时/异常时）重新查找一次，并且查找逻辑必须不引入 callback 泄漏（复用 callback 或低频路径创建）。

### 决策 3：`ForceMoveResizeWindowByTitle` 调整为基于 hwnd 的版本，并保留薄封装

**选择**：新增/迁移为 `ForceMoveResizeWindow(hwnd, x, y, w, h)` 这类 API；`ForceMoveResizeWindowByTitle` 仅作为“非轮询路径”兼容封装或调试工具。

**理由**：以 hwnd 为核心 API 可以从根源上切断“按标题枚举窗口”的依赖，让轮询路径无 callback 创建。

### 决策 4：DB 初始化在 `CGO_ENABLED=0` 下使用纯 Go sqlite driver

**选择**：将 sqlite driver 切换为纯 Go 实现（例如 `modernc.org/sqlite` 对应的 GORM driver 适配，或等价方案），保证在无 cgo 环境下依然可用。

**理由**：

- 用户日志已出现 `go-sqlite3 requires cgo`，这是“确定会发生”的环境差异。
- 纯 Go driver 能避免 Windows/CI 上配置 cgo 工具链的复杂性。

**备选方案**：

- 强制要求 `CGO_ENABLED=1` 并在构建脚本/README 中说明 → 对用户环境要求高，且不利于分发。
- 启动时发现 DB 不可用则禁用相关功能 → 会显著降低应用可用性，不符合“脚本管理/分类”为核心功能的定位。

### 决策 5：服务层对 DB 不可用做显式错误返回与日志记录

**选择**：当 `db.InitDB` 返回错误时，应用启动阶段记录错误日志；CategoryService/ScriptService 的对外方法在 DB 未就绪时返回可诊断错误，而不是触发 nil pointer。

**理由**：这是稳定性底线；即使 DB 初始化仍出现罕见失败，也不应把错误扩散成崩溃。

## Risks / Trade-offs

- **[风险] hwnd 缓存失效（窗口重建/标题变化）** → 通过检测 Win32 API 返回值与 `IsWindow`（如需要）在低频路径做一次性重取。
- **[风险] 纯 Go sqlite driver 在性能/兼容性上与 go-sqlite3 不完全一致** → 只要满足当前 CRUD 场景即可；启用 WAL/外键等 PRAGMA 需确认支持情况。
- **[风险] 变更触及窗口跟随关键路径** → 按 Windows 平台分支实现，其他平台保持不变；并在 DPI/最小化/恢复等场景回归验证。
- **[权衡] 保留 ByTitle 封装会继续包含 EnumWindows 逻辑** → 但必须确保不在高频轮询路径中调用，或内部 callback 复用到进程级。
