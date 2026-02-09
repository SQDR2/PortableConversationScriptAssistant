## Why

在 Windows 环境下，助手窗口绑定目标应用后，顶部能与目标窗口齐平，但底部始终无法对齐——要么超出目标窗口，要么不够。根本原因是窗口尺寸计算中存在多个维度的 API 不一致：DWM Extended Frame Bounds 与 GetWindowRect 混用导致阴影像素偏差，以及 Wails 框架内部的 DPI 缩放与物理像素坐标混用导致高度被错误放大或缩小。

## What Changes

- 修复 `GetWindowDecorationHeightByTitle` 函数，统一使用 DWM API 获取 sidekick 自身的视觉边界，消除阴影像素偏差
- 重构 `checkTarget` 中的对齐逻辑，绕过 Wails 的 `WindowSetSize` / `WindowSetPosition`（其内部有 DPI 缩放和监视器偏移），改为直接使用 Win32 `SetWindowPos` API 进行物理像素级别的精确定位
- 在 sidekick 的窗口尺寸计算中，正确处理 DWM frame bounds 与 GetWindowRect 之间的偏移量（阴影补偿），确保视觉边界精确对齐

## Capabilities

### New Capabilities

（无新增能力）

### Modified Capabilities

- `window-follow`: 修改窗口高度同步的实现要求——明确要求在 Windows 平台上使用 DWM API 获取 sidekick 自身的视觉边界进行对齐计算，并要求绕过框架层的尺寸设置 API 直接使用原生 Win32 API 定位

## Impact

- **受影响文件**:
  - `backend/utils/window_windows.go` — `GetWindowDecorationHeightByTitle` 函数重构及新增 DWM 辅助函数
  - `backend/services/window_service.go` — `checkTarget` 对齐逻辑重写
- **平台影响**: 仅影响 Windows 平台，Linux/macOS 不受影响
- **破坏性**: 无，纯内部实现优化
