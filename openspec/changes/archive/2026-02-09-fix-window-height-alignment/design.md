## Context

当前助手窗口跟随目标窗口的高度同步逻辑位于 `backend/services/window_service.go` 的 `checkTarget` 方法中。其工作流程为：

1. 通过 `DwmGetWindowAttribute(EXTENDED_FRAME_BOUNDS)` 获取目标窗口的视觉矩形（不含阴影）
2. 通过 `GetWindowDecorationHeightByTitle("sidekick")` 计算 sidekick 自身的装饰高度（标题栏+边框）
3. 用 `目标视觉高度 - 装饰高度` 得到期望的 sidekick 窗口高度
4. 调用 Wails 的 `WindowSetSize` 和 `WindowSetPosition` 设置尺寸和位置

该流程在 Windows 上存在三个关键问题导致底部不齐平：
- 装饰高度使用 `GetWindowRect - GetClientRect` 计算，包含了 Win10/11 的窗口阴影像素
- Wails 的 `SetSize` 内部会做 DPI 缩放（`scaleWithWindowDPI`），而传入的是物理像素
- Wails 的 `SetPos` 会加上工作区偏移（`workRect.Left + x`），而传入的是屏幕绝对坐标

## Goals / Non-Goals

**Goals:**
- 在 Windows 上实现 sidekick 窗口与目标窗口的顶部和底部精确视觉对齐
- 在不同 DPI 缩放比例（100%/125%/150%/175%/200%）下都能正确对齐
- 保持 Linux 和 macOS 平台的现有行为不变

**Non-Goals:**
- 不改变 Wails 框架本身的行为
- 不重构窗口跟随的整体架构（轮询机制、Z-Order 管理等）
- 不处理多显示器切换场景中的对齐问题（后续改进）

## Decisions

### 决策 1: 在 Windows 上完全绕过 Wails 窗口管理 API 进行定位

**选择**: 对齐逻辑中统一使用 `ForceMoveResizeWindowByTitle`（直接调用 Win32 `SetWindowPos`），不再使用 `runtime.WindowSetSize` / `runtime.WindowSetPosition`。

**理由**: Wails 的窗口 API 设计用于普通应用的窗口管理（逻辑像素、相对坐标），不适合像素级的精确定位需求。其内部的 DPI 缩放和监视器偏移处理会引入不可控的误差。直接使用 Win32 API 是最可靠的方式。

**备选方案**:
- 反向计算 Wails 的 DPI 缩放因子来补偿 → 脆弱、依赖 Wails 内部实现细节
- 修改 Wails 源码 → 维护成本太高

### 决策 2: 使用 DWM 偏移量计算替代装饰高度减法

**选择**: 不再使用 `decorationHeight` 减法模式。改为获取 sidekick 自身的 DWM frame bounds 和 GetWindowRect 的差值，计算出精确的阴影偏移，然后：

```
sidekick_dwm = DwmGetWindowAttribute(sidekick, EXTENDED_FRAME_BOUNDS)
sidekick_wr  = GetWindowRect(sidekick)

shadow_top    = sidekick_dwm.Top - sidekick_wr.Top      // 通常 ≈ 0
shadow_bottom = sidekick_wr.Bottom - sidekick_dwm.Bottom // 通常 ≈ 7-8px
shadow_left   = sidekick_dwm.Left - sidekick_wr.Left
shadow_right  = sidekick_wr.Right - sidekick_dwm.Right

// SetWindowPos 使用的尺寸 = 目标视觉高度 + 上下阴影补偿
setWindowPos_height = target_dwm_height + shadow_top + shadow_bottom
setWindowPos_y      = target_dwm.Top - shadow_top
```

**理由**: 这种方式精确补偿了阴影的存在，使 sidekick 的 DWM 视觉边界与目标窗口的 DWM 视觉边界完全重合。

**备选方案**:
- 只改 `GetWindowDecorationHeightByTitle` 用 DWM → 只修一半问题，还需要处理 Wails API 的 DPI 缩放
- 设 `DWM_WINDOW_CORNER_PREFERENCE` 去掉圆角阴影 → 不可行，会影响整体视觉

### 决策 3: 新增 `GetDWMFrameOffsets` 辅助函数

**选择**: 在 `window_windows.go` 中新增一个函数，给定窗口句柄返回 DWM frame bounds 和 GetWindowRect 之间四个方向的偏移量。

**理由**: 封装通用逻辑，`checkTarget` 中调用更清晰。也方便未来其他需要 DWM 偏移的场景复用。

## Risks / Trade-offs

- **[风险] DWM API 在某些旧版 Windows 上不可用** → 已有 fallback 到 GetWindowRect 的逻辑，保持兼容
- **[风险] 绕过 Wails API 可能导致框架内部状态不一致** → 影响极小，因为我们只是设置物理位置和大小，框架的内部状态（min/max size 约束等）不依赖外部位置变化
- **[权衡] 完全绕过 Wails 意味着未来 Wails 升级不会自动修复该部分** → 可接受，窗口精确定位是我们的核心需求，必须自己掌控
