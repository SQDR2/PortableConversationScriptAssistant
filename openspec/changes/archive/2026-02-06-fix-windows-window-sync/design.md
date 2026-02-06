## Context

目前 `WindowProvider.GetWindowRect` 在 Windows 平台仅调用 `GetWindowRect`，这在 Windows 10/11 上返回的是包含不可见阴影的物理位置。同时，`WindowService` 的同步循环仅根据矩形（Rect）和图标状态（Iconic）的变化来决定是否更新，导致：

1. **视觉误差**：获取的 Rect 偏大（包含了左右和底部各 7-8 像素的半透明阴影区域）。
2. **层级丢失**：目标窗口被其他应用挡住后再被用户点开时，助手窗口无法识别出这种“视觉上变到最前”的变化，从而继续留在底层。

## Goals / Non-Goals

**Goals:**

- 在 Windows 平台实现厘米级的视觉对齐。
- 实现助手窗口与目标窗口的层位（Z-Order）同步。
- 解决从最小化恢复或从后台切换时的显示延迟。

**Non-Goals:**

- 实现完美的窗口吸附手感（拖拽时的动态平滑）。
- 处理所有窗口的层级关系（仅处理目标窗口与自身）。

## Decisions

### 1. 使用 DWM API 获取视觉矩形

- **Decision**: 在 `backend/utils/window_windows.go` 中动态加载 `dwmapi.dll`，使用 `DwmGetWindowAttribute` 获取 `DWMWA_EXTENDED_FRAME_BOUNDS`。
- **Rationale**: 这是 Windows 10+ 获取真实窗口边界的官方标准。它可以自动处理 DPI 缩放和阴影边框带来的干扰。
- **Alternative**: 手动根据 DPI 和主题硬编码修正值。但在不同系统配置下极易出错且维护成本高。

### 2. 引入前台窗口（Foreground）状态跟踪

- **Decision**: 在 `WindowProvider` 定义中增加 `GetForegroundHandle` 的封装。在 `WindowService` 的 `checkTarget` 逻辑中，每轮循环都检查当前系统前台窗口句柄。
- **Rationale**: 如果当前前台窗口是目标应用，且助手窗口不是当前前台窗口，则必须进行强制对齐和置顶。这打破了“仅位置变化才更新”的限制，解决了遮挡问题。

### 3. Always-on-top strategy

- **Decision**: 采用“瞬时提升，即刻释放”方案。当目标获得焦点时，临时开启 `AlwaysOnTop` 并执行平台强制 Raise 動作，成功后**立即撤销**全局置顶状态。
- **Rationale**: 持续性的置顶会导致窗口管理器在处理其他应用覆盖时出现竞态条件（需要点击两次才生效）。通过瞬时提升，我们将助手送到微信之上的 Z-Order 顶端，但不赋予其“霸占层级”的永久属性，从而允许用户无感切换到任意其他应用。

### 4. Z-Order 恢复增强

- **Decision**: 引入 `ForceRaiseWindowByTitle` 函数。在 Windows 下封装 `SetWindowPos(HWND_TOPMOST)`，在 Linux 下封装 `XRaiseWindow` 或 `_NET_ACTIVE_WINDOW`。
- **Rationale**: 仅仅依靠显示/不显示在全屏切换或复杂窗口堆栈下无法稳定抢夺 Z-Order。显式的 Raise 调用能强制 WM 重新排列堆栈。

## Risks / Trade-offs

- **[Risk]**：DWM API 在禁用 AERO 特效或极老系统上可能执行失败。
  - **Mitigation**：增加分层获取逻辑，如果 DWM 调用返回错误，立即降级回 `GetWindowRect`，确保程序不崩溃且仍能运行。
- **[Trade-off]**：增加 `GetForegroundWindow` 调用微量增加了 CPU 轮询负担。
  - **Mitigation**：Win32 API 调用极其高效，50ms 的间隔对现代系统几乎无感知。
