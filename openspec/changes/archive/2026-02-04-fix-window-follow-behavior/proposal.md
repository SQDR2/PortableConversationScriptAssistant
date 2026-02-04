## Why

当前助手窗口在关联应用后不会立即同步高度，且在目标窗口最小化后恢复时出现位置漂移（跑到中间顶部）。这会破坏“贴右侧伴随”的核心体验，需要尽快修复以保证跟随一致性。

## What Changes

- 关联应用后立即重置助手窗口高度为目标窗口高度，并将位置贴紧目标右侧外边缘。
- 目标窗口最小化后恢复时，助手保持与最小化前一致的相对位置（仍贴右侧），避免恢复后位置漂移。
- 增加恢复过程的稳定处理，确保首次有效定位不会被跳过。

## Capabilities

### New Capabilities

- `window-follow`: 定义助手窗口与目标窗口的跟随与恢复行为规范（贴右侧、同步高度、恢复稳定性）。

### Modified Capabilities

- （无）

## Impact

- backend/services/window_service.go（窗口跟随逻辑）
- Wails runtime 窗口 API（WindowSetPosition/WindowSetSize/WindowShow/WindowHide 等调用路径）
