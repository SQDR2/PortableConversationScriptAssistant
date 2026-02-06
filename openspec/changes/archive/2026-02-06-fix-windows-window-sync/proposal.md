## Why

系统在 Windows 平台上存在窗口同步精度和响应问题，具体表现为：

1. **高度不一致**：助手窗口高度包含物理阴影边框，导致视觉上比关联应用高出一截。
2. **层级不跟随**：关联应用被其他应用遮挡或从最小化恢复后，助手窗口无法自动跟随展示到最前端。

## What Changes

- **对齐精度优化**：在 Windows 平台改用 `DwmGetWindowAttribute` 获取窗口的 Extended Frame Bounds，以获取视觉上的真实矩形，消除透明阴影边框带来的高度误差。
- **层级同步增强**：在同步循环中增加前台窗口（Foreground Window）检查。当检测到目标应用获得焦点但在 Z-Order 上未被助手窗口跟随覆盖时，强制提升助手窗口层级。
- **稳定性改进**：修正 Windows 平台下 `IsIconic` (最小化) 状态检测后恢复时的布局抖动。

## Capabilities

### Modified Capabilities

- `window-follow`: 提升在 Windows 平台下的窗口对齐精度标准和层位（Z-Order）跟随响应要求。

## Impact

- `backend/utils/window_windows.go`: 引入 `dwmapi.dll` 及其 `DwmGetWindowAttribute` 方法。
- `backend/services/window_service.go`: 更新 `checkTarget` 逻辑，集成层级检查和视觉边界修正。
