## 1. 基础模型与后端服务

- [x] 1.1 修改 `backend/models/script.go`，在 `Script` 结构体中添加 `Images string `json:"images"`` 字段。
- [x] 1.2 在 `backend/services/script_service.go` 中实现 `SaveScriptImage` 方法，使用 **UUID V7** 生成文件名并将图片保存至本地磁盘。
- [x] 1.3 注册 Wails 的 `AssetsHandler` 或类似的自定义协议处理器，以便前端能通过指定路径加载本地图片文件。
- [x] 1.4 修改 `DeleteScript` 方法，在删除数据库记录前，先解析并物理删除关联的图片文件。

## 2. 前端编辑器增强 (`ScriptsPage.vue`)

- [x] 2.1 在 `showEditor` 对应的对话框中添加图片上传区域（使用 `q-file`）。
- [x] 2.2 实现图片上传后的临时预览逻辑。
- [x] 2.3 修改 `saveScript` 逻辑，在保存话术前先将新图片存储到后端，并更新地址列表。
- [x] 2.4 在编辑器中支持单张图片的移除操作。

## 3. 话术展示与预览 (`ScriptItem.vue`)

- [x] 3.1 增加预览对话框 `PreviewDialog`。
- [x] 3.2 实现点击话术卡片弹出预览功能，展示完整文本和图片列表。
- [x] 3.3 为 `ScriptItem` 增加缩略图预览（若存在图片，在卡片上显示一个小图标或第一张图的微型缩略图）。

## 4. 话术复制功能

- [x] 4.1 话术复制功能回归至纯文本复制，以确保在 Linux 等各环境下的极高可靠性。
- [x] 4.2 前端点击“复制”按钮时，使用 `navigator.clipboard.writeText` 快速完成。
- [x] 4.3 若需获取图文，建议用户直接进入详情页进行手动处理或后续扩展。

## 5. 问题修复与体验优化

- [x] 5.1 修复 `ScriptItem.vue` 中三个按钮（复制、编辑、删除）的垂直居中显示问题。
- [x] 5.2 优化编辑器中的图片上传 `q-file` 和预览区域的 UI 显示。
- [x] 5.3 修复时间轴渲染卡顿问题：重构了 `displayedScripts` 的计算逻辑，避免了原地排序（in-place sort）导致的响应式死循环和性能开销。
- [x] 5.4 修复 `q-file` 上传组件在选中文件后显示文件名的问题，保持纯图标设计。
- [x] 5.5 优化话术保存后的加载逻辑，避免列表出现空白闪烁（优化 `loading` 状态处理）。
- [x] 5.6 修复删除话术时的刷新逻辑，确保使用静默刷新（silent refresh）并使用 `QInnerLoading` 优化视觉体验。
