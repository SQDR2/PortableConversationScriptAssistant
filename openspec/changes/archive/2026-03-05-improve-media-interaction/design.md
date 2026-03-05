## Context

当前媒体文件的服务分为两条路径：
- **图片**：通过 Wails AssetServer 的自定义 handler 在 `/images/` 路径下提供，`<img>` src 为相对路径
- **视频**：通过独立的 `MediaServer`（本地 127.0.0.1 HTTP server）提供，以支持 Range 请求和正常的 `<video>` 流式播放

在生产构建中，Wails WebView 使用自定义 URI scheme（非真实 HTTP），默认关闭浏览器原生右键菜单，导致图片无法通过原生 "Copy Image" 复制。视频从未有过"分享"能力，无论是 dev 还是 build 模式均不支持。

## Goals / Non-Goals

**Goals:**
- 在生产构建中，用户右键单击图片能看到自定义菜单，并可复制图片到剪贴板
- 用户可以通过菜单/按钮在系统文件管理器中定位任意媒体文件（图片或视频）
- 实现在 Linux / macOS / Windows 三平台均可用

**Non-Goals:**
- 不实现与第三方 IM 客户端的直接集成
- 不实现从应用窗口到系统桌面的拖拽（Wails WebView drag-out 支持有限）
- 不支持 JPEG 图片的剪贴板复制（ClipboardItem 仅标准支持 PNG；视实现情况可降级处理）

## Decisions

### 决策一：自定义 Vue 右键菜单，而非 `EnableDefaultContextMenu: true`

**选择**：在 `<img>` 上监听 `@contextmenu.prevent`，用 Quasar `QMenu` 渲染菜单项。

**理由**：`EnableDefaultContextMenu` 会暴露浏览器全部原生菜单（包括"审查元素"等不应对终端用户展示的选项），且在 Linux WebKitGTK 上"Copy Image"是否能处理 Wails 自定义协议不可预期。自定义菜单行为完全可控，UI 风格与 Quasar 一致。

**菜单项**（图片）：
- **复制图片**（Copy to Clipboard）
- **在文件夹中显示**（Reveal in File Manager）

### 决策二：使用 `fetch` + `ClipboardItem` 复制图片

**选择**：通过 `fetch(imgSrc)` 获取图片 Blob → `new ClipboardItem({'image/png': blob})` → `navigator.clipboard.write()`。

**理由**：此路径完全在 WebView JS 层面完成，不依赖浏览器原生"Copy Image"，绕过 Wails 自定义 URI scheme 对原生菜单的限制。Wails WebView 内置安全上下文（trusted context）支持 Clipboard API。

**降级**：若图片非 PNG（如 JPEG），直接将 `src` URL 写入文本剪贴板并提示用户。实际上多数 IM 客户端可识别"复制图片"后的 PNG。

### 决策三：Go 后端实现 `RevealInFileManager`，挂载到现有 Service

**选择**：在 `backend/services/script_service.go` 中新增 `RevealInFileManager(relativePath string) error` 方法（也可新建独立 service，但为减少改动量挂载到 ScriptService）。

**平台实现**：
```
Linux:   exec("xdg-open", filepath.Dir(absPath))   # 打开所在目录
macOS:   exec("open", "-R", absPath)               # 高亮选中文件
Windows: exec("explorer", "/select,"+absPath)      # 高亮选中文件
```

**理由**：纯 Go 标准库 `os/exec` 实现，无新增依赖。路径从相对路径（如 `images/uuid.mp4`）转换为绝对路径由 `os.Getwd()` 完成，与现有 MediaServer 逻辑一致。

### 决策四：视频区域使用图标按钮而非右键菜单

**选择**：在视频元素旁新增一个「在文件夹中显示」图标按钮（与现有「全屏」按钮风格一致），不额外添加右键菜单。

**理由**：视频没有"复制"操作，菜单项只有一个时不如直接暴露按钮；与全屏按钮保持视觉一致性。

## Risks / Trade-offs

- **ClipboardItem PNG 限制** → 降级为提示用户"图片已存储在本地，可通过文件管理器发送"
- **Linux xdg-open 打开目录而非选中文件**（部分文件管理器如 Nautilus 不支持选中指定文件）→ 可接受，用户至少能打开所在目录
- **Wails Secure Context**：Clipboard API 要求 secure context。Wails WebView 将自身 origin 视为可信，应无此问题，但需在构建后验证
- **Windows 路径含空格**：`explorer /select,` 路径参数需要正确引用 → 使用 `exec.Command` 数组参数形式避免 shell 注入
