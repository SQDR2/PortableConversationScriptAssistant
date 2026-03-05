## 1. Go 后端：RevealInFileManager 实现

- [x] 1.1 在 `backend/services/script_service.go` 中新增 `RevealInFileManager(relativePath string) error` 方法，接受相对路径（如 `images/uuid.mp4`）并转换为绝对路径
- [x] 1.2 实现 Linux 平台逻辑：路径存在检查通过后执行 `xdg-open <目录路径>`（打开文件所在目录）
- [x] 1.3 实现 macOS 平台逻辑：执行 `open -R <文件绝对路径>`（高亮选中文件）
- [x] 1.4 实现 Windows 平台逻辑：执行 `explorer /select,<文件绝对路径>`（高亮选中文件，注意路径引用）
- [x] 1.5 路径不存在时返回有意义的 error（前端显示「文件不存在或已被删除」通知）
- [x] 1.6 运行 `wails dev` 或 `wails generate module` 重新生成 `frontend/wailsjs/` 绑定

## 2. 前端：图片右键上下文菜单

- [x] 2.1 在 `ScriptItem.vue` 的话术卡片内、所有 `<img>` 元素上添加 `@contextmenu.prevent="showImgMenu($event, img)"` 绑定
- [x] 2.2 使用 Quasar `QMenu` + `QList`/`QItem` 实现上下文菜单 UI（跟随鼠标位置，Quasar 风格）
- [x] 2.3 实现 `copyImage(imgSrc: string)` 函数：`fetch(imgSrc)` → `Blob` → `new ClipboardItem({'image/png': blob})` → `navigator.clipboard.write()`
- [x] 2.4 实现复制失败降级逻辑：ClipboardItem 不支持时弹出错误通知，引导用户使用「在文件夹中显示」
- [x] 2.5 菜单项「在文件夹中显示」：调用 Wails 生成绑定中的 `ScriptService.RevealInFileManager(relativePath)`（relativePath 来自 scriptImages 数组原始值）
- [x] 2.6 在详情预览 Dialog 中的 `<img>` 上同步添加相同的右键菜单逻辑

## 3. 前端：视频「在文件夹中显示」按钮

- [x] 3.1 在 `ScriptItem.vue` 视频卡片区域（`.relative-position` 容器）中新增「在文件夹中显示」`q-btn` 图标按钮（使用 `folder_open` 图标，样式与全屏按钮一致，定位到视频控制栏旁边）
- [x] 3.2 按钮点击回调中调用 `RevealInFileManager(vid)`（vid 为视频原始相对路径）
- [x] 3.3 在详情预览 Dialog 的视频区域（`.relative-position` 容器）中同步添加「在文件夹中显示」按钮
- [x] 3.4 处理 `RevealInFileManager` 调用失败时的 Quasar `notify` 错误通知

## 4. 验证与测试

- [x] 4.1 运行 `go build ./...` 确保后端代码无编译错误
- [ ] 4.2 在 `wails dev` 模式下验证图片右键菜单可正常弹出、复制图片可粘贴到 IM
- [ ] 4.3 在 `wails dev` 模式下验证「在文件夹中显示」可打开文件管理器（Linux：确认 Nautilus/Thunar 等打开了正确目录）
- [ ] 4.4 执行 `wails build` 后，在生产构建中验证图片右键菜单功能保持正常
- [ ] 4.5 验证视频「在文件夹中显示」按钮在卡片视图和全屏 Dialog 中均可见且可用
