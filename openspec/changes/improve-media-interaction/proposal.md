## Why

生产构建（`wails build`）中图片和视频的右键菜单完全缺失，导致用户无法快速复制图片发送给客户；同时视频媒体没有任何"分享"路径——用户无法将视频文件方便地发送给对方。这两个问题严重影响客服人员的使用效率，必须在当前版本修复。

## What Changes

- **自定义图片右键菜单**：在 `<img>` 元素上拦截原生 `contextmenu` 事件，展示 Quasar 风格的上下文菜单，提供"复制图片"和"在文件夹中显示"两个操作。复制图片使用 `fetch` + `ClipboardItem` API，完全绕开 Wails 自定义协议的限制，在 dev 和 build 模式均可正常工作。
- **视频添加「在文件夹中显示」按钮**：在视频卡片上新增操作按钮，点击后调用 Go 后端打开 OS 文件管理器并定位到对应的视频文件，用户可直接拖拽文件到 IM 客户端上传。
- **Go 后端新增 `RevealInFileManager` 方法**：跨平台实现，Linux 使用 `xdg-open`，macOS 使用 `open -R`，Windows 使用 `explorer /select,`。

## Capabilities

### New Capabilities

- `image-context-menu`：图片右键自定义菜单能力，包含"复制图片到剪贴板"和"在文件夹中显示"操作，在生产构建中正常工作。
- `media-reveal-in-folder`：通过 Go 后端调用系统文件管理器，定位并高亮指定的媒体文件（图片或视频），跨平台支持 Linux/macOS/Windows。

### Modified Capabilities

（无已有 spec 的需求变更）

## Impact

- **前端**：`frontend/src/components/ScriptItem.vue` — 加右键菜单逻辑、视频区域新增按钮
- **后端**：`backend/services/script_service.go` 或新文件 — 新增 `RevealInFileManager(path string)` 方法，需在 `main.go` Bind 列表中注册
- **Wails 绑定**：运行 `wails dev` 或 `wails generate module` 重新生成 `wailsjs/` 绑定
- **依赖**：无新增第三方依赖
