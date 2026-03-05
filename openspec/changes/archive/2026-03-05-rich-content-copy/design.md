## Context

当前 `ScriptItem.vue` 有两条独立的复制路径：

- `copyContent()` → `navigator.clipboard.writeText(text)` — 仅写入纯文本
- `copyImage(imgSrc)` → `navigator.clipboard.write([ClipboardItem({ 'image/png': blob })])` — 仅写入单张图片

两者互相覆盖剪贴板，用户无法通过一次操作同时获得文本和图片。目标是在一次点击内将富文本（文字 + 所有图片）写入剪贴板。

## Goals / Non-Goals

**Goals:**

- 点击话术内容区时，若话术含图片，将文本 + 所有图片一并写入剪贴板（`text/html` + `text/plain`）
- 无图片话术保持现有 `writeText()` 纯文本路径，无行为变更
- 保持「用户手势同步帧」约束：`clipboard.write()` 在点击回调同步帧内调用

**Non-Goals:**

- 不处理视频复制（视频无法通过 Clipboard API 传递，保持右键菜单单独操作）
- 不改变右键菜单中单独复制单张图片的行为
- 不支持视频+文本的组合复制

## Decisions

### 决策 1：使用 `text/html` + `text/plain` 双格式而非 `image/png`

**选择**：将所有图片 base64 编码后嵌入 HTML `<img>` 标签，同时提供 `text/plain` fallback。

**原因**：浏览器 Clipboard API 的单个 `ClipboardItem` 不支持 `text/*` 和 `image/*` 混合。`text/html` 是唯一能把文本和多图打包为一个剪贴板条目的标准格式。微信、飞书、企业微信等主流 IM 的富文本粘贴均识别 `text/html`。

**替代方案考虑**：
- 多个 `ClipboardItem`（先文本后图片）→ 被否决，多次 `clipboard.write()` 每次覆盖，最终只保留最后一次
- 只复制第一张图片 → 被否决，用户明确要求全部图片

### 决策 2：`Promise<Blob>` 模式维持用户手势上下文

**选择**：在`handleContentClick` 的同步帧里构造 `ClipboardItem`，将异步操作（图片加载、canvas 转 base64）封装在 `Promise<Blob>` 的 value 里传入，不等待图片加载完成再调用 `clipboard.write()`。

**原因**：Wails WebView 对 Clipboard API 的权限检查依赖用户手势帧，`clipboard.write()` 必须在同步帧内调用。该模式已被现有 `copyImage()` 验证可行。

### 决策 3：多图并行处理

**选择**：使用 `Promise.all` 并行加载和转换所有图片，再组装 HTML 字符串。

**原因**：串行加载在图片多时延迟明显，并行不增加实现复杂度，且 canvas 操作无全局状态冲突。

## Risks / Trade-offs

- **微信 PC 兼容性** → 微信 PC 对 `text/html` 内嵌 `data:image/png;base64` 的处理行为未有官方文档记载，实现后须在微信 PC 端实测粘贴效果。若不兼容，可考虑降级为先粘文本、提示用户单独复制图片。
- **大图 base64 体积** → 单张高分辨率图片 base64 后可能达数 MB，多张累积可能给剪贴板和对端 IM 带来压力。缓解：可在 canvas 绘制时限制最大宽度（如 1920px），但此优化可留作后续迭代，本次不做。
- **canvas 跨域限制** → 图片通过本地 MediaServer（`localhost:port/filename`）提供，无跨域问题，canvas `toBlob` 可正常执行。

## Migration Plan

仅涉及 `ScriptItem.vue` 的函数替换，无数据 schema 变更，无需数据库迁移。Wails bindings 无变更，直接替换函数实现即可上线。回滚：还原 `ScriptItem.vue` 的对应函数。

## Open Questions

- 微信 PC 端实测结果如何？（需手动验证）
- 是否需要对单张图片大于某阈值时给出警告提示？（保留为后续优化项）
