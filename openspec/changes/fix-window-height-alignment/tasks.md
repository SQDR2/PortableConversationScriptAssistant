## 1. 新增 DWM 偏移量辅助函数

- [x] 1.1 在 `backend/utils/window_windows.go` 中新增 `GetDWMFrameOffsets(hwnd)` 函数，返回 DWM Extended Frame Bounds 与 GetWindowRect 四个方向的偏移量（top, bottom, left, right）
- [x] 1.2 为偏移量定义返回类型 `FrameOffsets struct`（Top, Bottom, Left, Right int），放在 `backend/utils/types.go` 中（仅 Windows 需要，但类型可跨平台定义）

## 2. 重构对齐逻辑

- [x] 2.1 修改 `backend/services/window_service.go` 中 `checkTarget` 的对齐代码块：移除对 `GetWindowDecorationHeightByTitle` 的调用和 `decorationHeight` 相关字段
- [x] 2.2 在对齐逻辑中调用新的 `GetDWMFrameOffsets` 获取 sidekick 自身的阴影偏移量
- [x] 2.3 计算 `SetWindowPos` 所需的物理坐标和尺寸：`y = target_dwm.Top - shadow_top`，`height = target_dwm_height + shadow_top + shadow_bottom`
- [x] 2.4 替换 `runtime.WindowSetSize` + `runtime.WindowSetPosition` 为直接调用 `ForceMoveResizeWindowByTitle`（统一使用 Win32 `SetWindowPos`）
- [x] 2.5 移除 fallback 逻辑（`actualX != newX` 的检查），因为现在始终使用直接 API

## 3. 清理与兼容

- [x] 3.1 移除 `WindowService` 中不再需要的 `decorationHeight` 和 `decorationKnown` 字段
- [x] 3.2 确认 `GetWindowDecorationHeightByTitle` 函数是否还有其他调用方；若无，将其标记为废弃或移除
- [x] 3.3 确保 `window_linux.go` 和 `window_darwin.go` 的对齐路径不受影响（条件编译隔离）

## 4. 验证

- [ ] 4.1 在 100% DPI 的 Windows 环境下验证顶部和底部对齐
- [ ] 4.2 在 150%（或其他非 100%）DPI 环境下验证对齐
- [ ] 4.3 验证目标窗口调整大小后助手底部仍保持对齐
