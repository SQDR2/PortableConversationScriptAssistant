## Why

当前话术管理系统仅支持图片附件。用户在实际业务场景中需要为话术附加视频素材（如产品演示短视频、操作教程等），以提升话术的表达力和实用性。现有的图片管线（上传 → 存储 → 展示 → 删除）已经成熟，视频功能可以在此基础上自然扩展，改动成本低。

## What Changes

- 话术编辑器支持上传视频文件（`.mp4`, `.webm`），与图片共享同一上传入口
- 话术列表和详情预览中，视频以内联播放器（`<video>` 标签）展示
- 列表缩略图区域：视频显示为播放图标，图片保持现有 `<q-img>` 缩略图
- 复用现有 `images` 字段和 `images/` 存储目录，前端通过文件扩展名区分渲染方式
- 新增 Go 后端接口 `SaveScriptMedia`，支持通过文件路径直接拷贝（避免大文件 base64 传输瓶颈）
- 图片+视频总附件数量上限保持 10 个
- 单个视频文件大小限制 ≤ 50MB

## Capabilities

### New Capabilities
- `video-storage`: 视频文件的上传、存储、静态服务和大文件传输策略

### Modified Capabilities
- `image-storage`: 存储目录 `images/` 现在同时承载视频文件；文件保存接口泛化以支持视频
- `script-management`: 话术数据模型的 `images` 字段语义扩展为"媒体附件"，兼容视频路径；编辑器 UI 和列表 UI 需区分图片/视频渲染

## Impact

- **后端代码**: `backend/services/script_service.go`（新增 `SaveScriptMedia` 接口）、`main.go`（AssetServer Handler 无需改动，`http.ServeFile` 已自动处理视频 MIME type）
- **前端代码**: `frontend/src/pages/ScriptsPage.vue`（编辑器上传逻辑扩展）、`frontend/src/components/ScriptItem.vue`（列表缩略图和预览弹窗视频渲染）
- **数据模型**: `backend/models/script.go` 无需改动（复用 `Images` 字段）
- **数据库**: 零迁移（`images` 列的 JSON 内容兼容视频路径）
- **存储**: `images/` 目录存储空间需求增长（视频文件较大）
- **系统依赖**: 无新外部依赖（不引入 ffmpeg 等）
- **Wails 绑定**: 新增 `SaveScriptMedia` 需要重新生成 `wailsjs` 绑定
