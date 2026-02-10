## 1. 后端：新增 `SaveScriptMedia` 接口

- [x] 1.1 在 `backend/services/script_service.go` 中新增 `SaveScriptMedia(filePath string) (string, error)` 方法
  - 接收本地文件绝对路径
  - 校验路径安全性（禁止 `..`，检查文件是否存在）
  - 校验文件扩展名（仅允许 `.mp4`, `.webm` 以及现有图片格式）
  - 校验文件大小（视频 ≤ 50MB）
  - 使用 UUID v7 重命名并拷贝到 `images/` 目录
  - 返回相对路径 `images/<uuid>.<ext>`
- [x] 1.2 重新生成 Wails 前端绑定（`wails generate module`）

## 2. 前端：上传逻辑扩展

- [x] 2.1 在 `frontend/src/pages/ScriptsPage.vue` 中添加 `isVideo(path: string): boolean` 工具函数，通过扩展名判断媒体类型
- [x] 2.2 新增视频上传按钮，调用 Wails `runtime.OpenFileDialog` 获取文件路径，然后调用 `SaveScriptMedia` 保存
- [x] 2.3 修改编辑器已上传附件预览区域：图片保持 `<q-img>` 70×70，视频显示 `play_circle` 图标 + 文件名

## 3. 前端：列表缩略图扩展

- [x] 3.1 修改 `frontend/src/components/ScriptItem.vue` 列表缩略图渲染逻辑：视频用 `<q-icon name="play_circle">` 替代 `<q-img>`
- [x] 3.2 badge 数量保持显示剩余附件总数（图片+视频合计）

## 4. 前端：预览弹窗扩展

- [x] 4.1 修改 `ScriptItem.vue` 预览弹窗：视频使用 `<video controls preload="metadata">` 标签渲染，图片保持现有 `<img>` 标签
- [x] 4.2 视频播放器样式：全宽、圆角、阴影，与图片预览保持视觉一致

## 5. 测试与验证

- [x] 5.1 验证 `.mp4` 和 `.webm` 视频上传、存储、展示完整流程
- [x] 5.2 验证视频大小超限（> 50MB）时错误提示正常
- [x] 5.3 验证不支持的视频格式被拒绝
- [x] 5.4 验证删除含视频话术时磁盘文件被清理
- [x] 5.5 验证现有纯图片话术功能不受影响（向后兼容）
