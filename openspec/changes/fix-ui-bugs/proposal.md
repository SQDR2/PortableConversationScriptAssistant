## Why

v2.0.3 在 Windows 和跨平台使用中暴露了三个 UI Bug：(1) Windows 下窗口宽度在关联目标后异常变窄；(2) 图片上传组件的内部状态未重置导致第二次上传完全无效果；(3) 滚动条作用在整个页面而非话术列表区域内，导致搜索栏和 Tab 切换器跟随滚动。这些问题严重影响日常使用体验，需要尽快修复。

## What Changes

- 修复 Windows 平台下首次 `checkTarget` 时 `sidekickHWID` 未解析导致使用错误默认宽度（350 逻辑像素被当作物理像素），使窗口在高 DPI 下缩窄约 33%
- 修复 `ScriptsPage.vue` 中 `q-file` 组件的 `uploadDummy` ref 在文件选择后未重置，导致后续上传无法触发 `@update:model-value` 事件
- 修复 `ScriptsPage.vue` 中 `q-page` 高度约束失效，将滚动范围限制在话术列表区域内，搜索栏和视图切换按钮保持固定

## Capabilities

### New Capabilities

无。

### Modified Capabilities

- `window-follow`: 新增要求——首次对齐时若尚未获取助手窗口句柄，MUST 跳过 resize 操作等待下一轮轮询，避免使用错误的默认宽度
- `image-storage`: 新增要求——图片上传组件在每次文件选择完成后 MUST 重置内部状态，确保连续多次上传均能正常触发
- `ui-translation-bar`: 新增要求——话术列表区域 MUST 为唯一可滚动区域，搜索栏和视图切换器 MUST 保持固定不随滚动

## Impact

- **后端**: `backend/services/window_service.go` — 调整 `checkTarget` 对齐逻辑，`sidekickHWID` 未就绪时跳过宽度 resize
- **前端**: `frontend/src/pages/ScriptsPage.vue` — 重置 `uploadDummy`、修复 `q-page` 高度约束和滚动容器定位
- **跨平台**: 滚动条修复影响所有平台（Windows + Linux）；窗口宽度修复仅影响 Windows
