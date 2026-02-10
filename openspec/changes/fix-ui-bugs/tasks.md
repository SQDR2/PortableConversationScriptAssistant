## 1. 修复 Windows 窗口宽度变窄

- [x] 1.1 在 `backend/services/window_service.go` 的 `checkTarget()` 对齐逻辑中，当 `sidekickHWID` 为空时跳过 `ForceMoveResizeWindow` 调用，直接进入 `update_tracker`
- [x] 1.2 验证 Windows 平台下首次关联目标窗口后助手窗口宽度不再缩窄

## 2. 修复图片上传第二次失败

- [x] 2.1 在 `frontend/src/pages/ScriptsPage.vue` 的 `handleFileUpload` 函数末尾添加 `uploadDummy.value = null` 重置上传组件状态
- [x] 2.2 在 `openEditor()` 函数中添加 `uploadDummy.value = null` 确保每次打开编辑器时上传组件为干净状态
- [x] 2.3 验证连续多次图片上传均能正常触发

## 3. 修复滚动条范围

- [x] 3.1 在 `frontend/src/pages/ScriptsPage.vue` 的 `<q-page>` 上使用 `:style-fn="tweakPageHeight"` 替代手写 `style="height: 100%"`，并添加 `overflow: hidden` 阻止页面级滚动
- [x] 3.2 确保话术列表容器 div 为唯一的 `overflow-y: auto` 滚动区域
- [x] 3.3 验证搜索栏和视图切换按钮在内容滚动时保持固定
- [x] 3.4 验证翻译栏显隐切换时列表区域高度自动调整

## 4. 修复视频全屏播放黑屏

- [x] 4.1 用 Quasar maximized dialog + Wails `WindowFullscreen()`/`WindowUnfullscreen()` 替代原生 Fullscreen API，`controlslist="nofullscreen"` 隐藏原生全屏按钮
