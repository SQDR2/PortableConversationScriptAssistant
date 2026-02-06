## 1. Windows Utils 增强 (Backend)

- [x] 1.1 在 `backend/utils/window_windows.go` 中通过 `syscall` 载入 `dwmapi.dll` 并定义 `procDwmGetWindowAttribute`。
- [x] 1.2 定义 `DWMWA_EXTENDED_FRAME_BOUNDS = 9` 常量及相关的 `winRect` 转换逻辑。
- [x] 1.3 在 `user32.dll` 封装中增加 `GetForegroundWindow` 获取函数。
- [x] 1.4 实现一个新的获取矩形方法，逻辑为：优先尝试 `DwmGetWindowAttribute` 获取视觉边界；若失败或 DWM 未启用，则降级使用 `GetWindowRect`。
- [x] 1.5 扩展 `WindowProvider` 接口，增加 `GetForegroundHandle() string` 方法，并在 Windows 下返回 HWND 字符串，Linux 下返回空或 DUMMY。

## 2. 同步逻辑优化 (Backend Service)

- [x] 2.1 修改 `backend/services/window_service.go` 中的 `checkTarget` 方法，调用增强后的视觉矩形获取逻辑。
- [x] 2.2 在轮询循环中调用 `GetForegroundHandle()` 获取当前系统焦点窗口。
- [x] 2.3 在 `WindowService` 中实现“瞬时置顶”触发逻辑：当目标获焦，助手闪烁开启/关闭 AlwaysOnTop。
- [x] 2.4 优化 `checkTarget`：移除所有持续性的置顶状态保持逻辑。
- [x] 2.5 确保 `ForceRaiseWindowByTitle` 在执行时包含瞬时的层级状态清理。
- [x] 2.6 在 `restoreSkipCount` 结束后，显式重置 `hasLast` 或 `forceUpdate`，确保重新计算正确的视觉边界而非沿用过时的物理边界。

## 3. 验证与回归测试

- [x] 3.1 验证 Windows 下助手窗口高度与目标应用（如微信、记事本）视觉边缘完全重合，不再高出一截。
- [x] 3.2 验证从任务栏激活被遮挡的目标应用时，助手窗口能瞬间跟随呈现在最前。
- [x] 3.3 验证 Linux (Ubuntu/EWMH) 环境下的自适应对齐功能依然正常，无回归。
