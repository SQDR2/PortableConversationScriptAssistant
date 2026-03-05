## 1. 核心函数重构（ScriptItem.vue）

- [x] 1.1 新增辅助函数 `imageUrlToBase64(url: string): Promise<string>`，使用 canvas 将图片 URL 转换为 `data:image/png;base64,...` 字符串
- [x] 1.2 新增 `copyRichContent()` 函数：检测 `scriptImages.value` 是否非空；非空时使用 `Promise.all` 并行调用 `imageUrlToBase64`
- [x] 1.3 在 `copyRichContent()` 中组装 HTML 片段（`<p>文本</p><img src="..."/><img src="..."/>`），并构造 `ClipboardItem({ 'text/plain': ..., 'text/html': ... })`，调用 `clipboard.write()`
- [x] 1.4 在 `copyRichContent()` 中处理无图片 fallback：当 `scriptImages` 为空时，退回到 `writeText()` 纯文本路径
- [x] 1.5 将 `handleContentClick()` 中的 `copyContent()` 调用替换为 `copyRichContent()`
- [x] 1.6 将"复制全文"按钮的 `@click="copyContent"` 替换为 `@click="copyRichContent"`

## 2. 用户反馈优化

- [x] 2.1 图文复制成功时 Toast 改为"已复制图文内容"，纯文本时保持"已复制文本内容"

## 3. 验证与测试

- [ ] 3.1 在飞书/企业微信粘贴：验证文字和所有图片同时出现
- [ ] 3.2 在系统记事本/纯文本编辑器粘贴：验证只显示文字内容（text/plain fallback 有效）
- [ ] 3.3 微信 PC 端粘贴测试（风险项，若不支持 base64 嵌入需记录兼容性问题）
