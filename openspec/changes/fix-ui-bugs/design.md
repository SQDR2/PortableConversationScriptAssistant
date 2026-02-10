## Context

Sidekick v2.0.3 在 Windows 和 Linux 平台上存在三个 UI 缺陷。这些问题均来源于前端状态管理和后端窗口跟随逻辑的实现细节遗漏，不涉及架构变更。

当前状态：
- **窗口宽度**: `window_service.go` 的 `checkTarget()` 在首次轮询时若 `sidekickHWID` 尚未解析，会使用硬编码的 `sw := 350` 作为物理像素宽度传给 `SetWindowPos`，但 Wails 的 350 是逻辑像素（在 150% DPI 下应为 525 物理像素）
- **图片上传**: `ScriptsPage.vue` 的 `uploadDummy` ref 仅在组件初始化时设为 `null`，后续上传完成后未重置，导致 Quasar `q-file` 组件不触发 `@update:model-value`
- **滚动范围**: `q-page` 使用 `style="height: 100%"` 但未利用 Quasar 的 `:style-fn` 机制精确计算可用高度，导致页面内容溢出后整个窗口出现滚动条

## Goals / Non-Goals

**Goals:**
- 修复 Windows 平台窗口宽度在高 DPI 下首次关联时被缩窄的问题
- 修复图片上传组件连续多次上传失败的问题
- 将滚动范围限制在话术列表区域内，搜索栏和 Tab 保持固定

**Non-Goals:**
- 不重构窗口跟随架构（仅修复首次对齐的边界条件）
- 不改变图片/视频上传的 API 签名
- 不引入新的 UI 组件或依赖

## Decisions

### D1: sidekickHWID 未就绪时跳过 resize

**决策**: 当 `sidekickHWID` 为空时，跳过当轮的 `ForceMoveResizeWindow` 调用，等待下一个 50ms 轮询周期解析句柄后再对齐。

**替代方案**:
- A) 在 `SetTarget` 时同步解析 `sidekickHWID` → 可能因窗口尚未创建而失败
- B) 使用 DPI-aware 的默认值替代 350 → 需要额外获取 DPI 比例，增加复杂度
- C) 跳过 resize（选中）→ 最安全，最多延迟 50ms，无副作用

### D2: handleFileUpload 完成后重置 uploadDummy

**决策**: 在 `handleFileUpload` 函数末尾（无论成功或失败）和 `openEditor()` 函数中均重置 `uploadDummy.value = null`。

**理由**: 双重保险策略。`handleFileUpload` 末尾重置确保当次上传后组件状态立即清空；`openEditor` 重置确保每次打开编辑器都是干净状态。

### D3: 使用 Quasar `:style-fn` 约束 q-page 高度

**决策**: 使用已存在但未启用的 `tweakPageHeight` 函数，通过 `:style-fn` 属性让 `q-page` 精确占满 header/footer 之间的空间，并在 `q-page` 上设置 `overflow: hidden` 阻止页面级滚动，仅在话术列表容器上启用 `overflow-y: auto`。

**替代方案**:
- A) 使用 CSS `calc(100vh - header - footer)` → 不可靠，header/footer 高度可变
- B) 使用 Quasar `:style-fn`（选中）→ 框架原生方案，自动计算 offset

## Risks / Trade-offs

- [D1] 首次关联时可能有 1-2 个 50ms 周期的延迟才会对齐 → 用户几乎不可察觉
- [D3] `overflow: hidden` 在 `q-page` 上可能影响 dialog 弹出的定位 → Quasar dialog 使用 `position: fixed`，不受影响
