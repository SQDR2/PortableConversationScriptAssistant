## Why

用户点击话术时，当前仅复制纯文本内容；若话术附有图片，必须再手动右键单独复制图片，操作割裂效率低下。需要在**一次点击**内将文本与所有图片同步写入剪贴板，覆盖已有 spec 中的核心场景。

## What Changes

- 将 `copyContent()` 升级为 `copyRichContent()`：当话术含有图片时，自动构造 `text/html`（内嵌所有图片 base64）+ `text/plain` 双格式写入剪贴板
- 无图片话术保持现有 `writeText()` 纯文本复制行为，不改变任何交互
- 复制成功的 Toast 提示文案根据是否有图片体现差异（"已复制文本内容" vs "已复制图文内容"）

## Capabilities

### New Capabilities

无新 capability，本 change 是对已有 spec 的补全实现。

### Modified Capabilities

- `rich-content-copy`：将已规划的"组合内容复制"需求落地到 `ScriptItem.vue`，实现文本 + 多图一次性复制。现有 spec 已包含准确的场景描述，本 change 按其要求实现，无需新增需求条目。

## Impact

- **前端**: `frontend/src/components/ScriptItem.vue` — `copyContent()` / `handleContentClick()` 逻辑替换
- **其他文件**: 无后端改动，无 API 变更，无数据库 schema 变更
- **依赖**: 仅依赖浏览器原生 Clipboard API（已在 `copyImage()` 中验证可用）
- **风险**: 微信 PC 端对 `text/html` 内嵌 base64 图片的支持行为需实测确认
