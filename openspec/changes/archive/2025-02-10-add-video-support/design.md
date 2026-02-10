## Context

当前话术管理系统的媒体管线仅支持图片。整体架构为：前端 `<q-file>` 组件选择文件 → `FileReader.readAsDataURL()` 读取为 base64 → Wails RPC 调用 `SaveScriptImage(base64, ext)` → Go 端 base64 解码 → UUID v7 命名写入 `images/` 目录 → 返回相对路径 → 存储到 `scripts.images` JSON 字段。

展示侧通过 Wails AssetServer 的自定义 Handler 拦截 `/images/*` 请求，使用 `http.ServeFile` 提供静态文件服务。

该架构对图片（通常 100KB-5MB）运行良好，但视频文件（通常 10MB-100MB+）的 base64 传输方式会导致内存峰值和传输延迟问题。

## Goals / Non-Goals

**Goals:**
- 在现有图片功能基础上，以最小改动量支持视频附件
- 解决大文件通过 base64 + RPC 传输的性能问题
- 向后兼容现有数据（零 DB 迁移）
- 图片和视频在 UI 中清晰区分渲染

**Non-Goals:**
- 不引入视频转码或压缩（不依赖 ffmpeg）
- 不实现视频缩略图的首帧截取（使用图标代替）
- 不新增独立的 `videos` 数据库字段（复用 `images` 字段）
- 不支持视频在线录制或编辑
- 不实现视频的富文本剪贴板支持

## Decisions

### Decision 1: 复用 `images` 字段，不新增数据库列

**选择**: 保留 `scripts.images` 字段名和 JSON 数组格式，视频路径直接存入该数组。

**替代方案**: 新增 `scripts.videos` 独立字段。

**理由**: 
- 零 DB 迁移，GORM AutoMigrate 不需要加列
- 前端只需维护一个 `editorImages[]` 数组
- 图片和视频的 CRUD 逻辑完全统一
- 前端通过文件扩展名判断渲染方式即可

### Decision 2: 大文件使用 Wails 原生文件对话框 + 文件路径拷贝

**选择**: 新增 `SaveScriptMedia(filePath string)` 后端接口，接收本地文件绝对路径，Go 端直接从磁盘读取并拷贝到 `images/` 目录。前端通过 `runtime.OpenFileDialog` 获取文件路径。

**替代方案 A**: 继续使用 base64 传输，限制文件大小。
**替代方案 B**: 前端分片读取，多次 RPC 调用拼接。

**理由**:
- 方案 A 对 50MB 的视频仍会导致 ~67MB 的 base64 字符串在内存中，不可接受
- 方案 B 实现复杂度高，需要处理断点续传等
- 文件路径拷贝完全绕过 RPC 数据传输，内存开销恒定，实现简单
- 此方案同时也优化了大图片的上传体验

### Decision 3: 存储目录保持 `images/`

**选择**: 视频文件仍存储在 `images/` 目录下。

**替代方案**: 新建 `videos/` 或 `media/` 目录。

**理由**:
- AssetServer Handler 已拦截 `/images/*`，零改动
- `http.ServeFile` 自动根据文件扩展名设置正确的 MIME type（`video/mp4`, `video/webm`）
- 减少路径管理复杂度

### Decision 4: 前端通过扩展名判断媒体类型

**选择**: 定义一个 `isVideo(path)` 工具函数，根据扩展名（`.mp4`, `.webm`）判断是否为视频。

**替代方案**: 在 JSON 中存储结构化对象 `{path, type}`。

**理由**:
- 扩展名判断简单可靠，无需改变 JSON 格式
- 向后兼容：现有纯图片数据无需任何变更
- 未来若需支持更多格式，只需扩展扩展名列表

### Decision 5: 视频缩略图使用 Material Icon

**选择**: 列表缩略图区域，视频用 `play_circle` 图标展示，不截取首帧。

**替代方案**: 使用 `<video>` + `<canvas>` 截取首帧做缩略图。

**理由**:
- 不需要额外的视频加载开销
- 不需要引入 ffmpeg 等外部依赖
- 列表性能不受影响（大量视频缩略图加载会拖慢列表）
- 图标方案视觉清晰，一目了然

### Decision 6: 支持的视频格式

**选择**: 仅支持 `.mp4` 和 `.webm`。

**理由**:
- `.mp4` (H.264): 所有平台 WebView 均支持
- `.webm` (VP8/VP9): Linux WebKitGTK 和 Windows WebView2 支持良好
- 不支持 `.mov`（macOS 专有，Linux/Windows WebView 支持差）
- 不支持 `.avi`, `.mkv`（WebView 不支持）

## Risks / Trade-offs

- **[风险] 文件路径安全** → `SaveScriptMedia` 接收任意文件路径，需要校验路径合法性（禁止 `..`、仅允许指定扩展名），防止路径遍历攻击
- **[风险] Linux WebKitGTK 视频解码** → 部分 Linux 发行版的 GStreamer 插件不全，H.264 可能无法播放 → 建议用户优先使用 `.webm` 格式，或在 README 中说明 GStreamer 依赖
- **[风险] 磁盘空间** → 视频文件远大于图片，长期使用可能占用大量空间 → 可在 UI 中显示存储用量提示（future scope）
- **[权衡] `images` 字段名语义不精确** → 字段名不再准确反映内容，但避免了 DB 迁移和代码翻倍
- **[权衡] 无视频压缩** → 用户上传的视频按原始大小存储，可能浪费空间 → 保持简单，不引入外部依赖
