## Context

当前窗口跟随逻辑在 backend/services/window_service.go 中以轮询方式更新位置。现有实现未能在关联后立即贴右侧外边缘，也未同步高度；且在目标窗口从最小化恢复时，位置更新可能被“短路判断”跳过，从而导致助手窗口停留在系统默认位置（例如中间顶部）。

约束：

- 使用 Wails runtime 窗口 API（WindowSetPosition/WindowSetSize/WindowShow/WindowHide），必要时在 Linux 上使用 X11 作为兜底定位手段；Wayland 下需强制 GDK_BACKEND=x11 以确保定位生效。
- 不能修改 .proto 及其生成文件。
- 以当前轮询机制为基础，尽量小改动提升稳定性。

## Goals / Non-Goals

**Goals:**

- 关联应用后，助手窗口立即同步高度并贴右侧外边缘。
- 目标窗口最小化后恢复时，助手保持与最小化前一致的相对位置（贴右侧），不发生位置漂移。
- 恢复流程稳定：首次有效定位不会被缓存短路跳过。

**Non-Goals:**

- 不引入新的跟随模式（如贴左侧、吸附角落）。
- 不做跨平台窗口 API 改造或新增依赖。
- 不调整前端 UI 布局或样式。

## Decisions

1. **在首次关联时立即同步高度**

- 决策：在 SetTarget() 触发或首次有效定位时，读取 rect.Bottom-rect.Top 并调用 WindowSetSize（保持当前宽度）。
- 理由：确保“关联即贴合”的体验，避免依赖后续移动才同步。
- 备选方案：仅在位置变化时同步高度；缺点是首次关联不一定同步。

2. **恢复后强制刷新定位一次**

- 决策：当检测到 lastIconic=true 且 isIconic=false 时，设置“待强制定位”标志并清空 hasLast/lastRect，确保恢复后首次有效定位不被短路判断拦截。
- 理由：现有短路判断会导致恢复后不重新定位，从而固定在系统默认位置。
- 备选方案：在短路判断前增加“恢复后强制刷新”分支；需要额外标记，复杂度略高但可行。

3. **最小化/恢复跟随而非隐藏**

- 决策：目标窗口最小化时调用 WindowMinimise；恢复时 WindowUnminimise 并立即显示助手窗口。
- 理由：与用户预期一致，保持“跟随最小化”的一致体验。
- 备选方案：最小化时隐藏；风险是用户误以为助手消失。

4. **Linux 下定位失败的 X11 兜底**

- 决策：当 WindowSetPosition 后的实际坐标与目标不一致时，使用 X11 根据窗口标题强制 Move/Resize。
- 理由：某些 Linux WM/会话对程序化定位支持有限，X11 兜底可确保“贴右侧”行为生效。
- 备选方案：完全依赖 Wails runtime；风险是定位可能被忽略。

5. **Wayland 会话强制 X11 后端**

- 决策：检测 Wayland 会话时设置 GDK_BACKEND=x11，确保窗口定位能力可用。
- 理由：Wayland 对客户端窗口定位限制较多，导致 WindowSetPosition 不生效。
- 备选方案：提示用户切换 X11 会话；风险是体验不一致。

## Risks / Trade-offs

- [恢复延迟可感知] → Mitigation：restoreSkipCount 设为极短窗口（例如 50~100ms）并尽量保持可见。
- [不同 WM 行为差异] → Mitigation：保持最小改动，仅使用现有 rect/isIconic 数据，避免引入新依赖。
- [短路逻辑误判] → Mitigation：恢复触发时强制刷新一次定位与尺寸。

## Migration Plan

- 无需数据迁移。
- 仅代码逻辑调整，发布后即时生效；如出现问题可回滚到上一版本。

## Open Questions

- restoreSkipCount 的最佳时长是否需按平台区分？（可在实现后通过日志观察再调参）
