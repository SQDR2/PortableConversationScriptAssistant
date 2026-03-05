## Why

话术列表中，内嵌图片（max-height: 200px）和视频（无高度上限）占据了过多垂直空间，压缩了文字内容的可见性，导致一屏能展示的话术条数偏少。需要收窄媒体预览区域，让话术列表更紧凑。

## What Changes

- 单图预览高度：`max-height: 200px` → `max-height: 120px`
- 多图 grid 单元格高度：`height: 120px` → `height: 80px`
- 视频预览新增高度上限：`max-height: 160px`
- 其余所有样式（布局、圆角、object-fit、交互）保持不变

## Capabilities

### New Capabilities

无新 capability。

### Modified Capabilities

- `script-management`：话术列表项的媒体预览区视觉尺寸缩减，不影响任何功能或数据。

## Impact

- **前端**：`frontend/src/components/ScriptItem.vue` — `<style>` 区块，3 处 CSS 数值修改
- **后端**：无
- **数据**：无
- **其他文件**：无
