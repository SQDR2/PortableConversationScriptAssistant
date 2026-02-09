## Why

当前的话术管理系统仅支持扁平化列表展示，随着话术数量的增加，用户难以有效组织和查找内容。用户需要一种结构化的方式来管理话术，同时保持界面的简洁性。

核心问题：

1.  话术缺乏分类，查找困难。
2.  话术列表展示信息冗余（不需要标题），且操作不够便捷。
3.  误删除风险高，缺乏安全确认机制。

本次变更旨在引入“目录（Category）”概念，提供更高效的组织方式，并优化话术的展示和操作体验。

## What Changes

本次变更将对后端数据模型和前端交互进行全面升级：

- **引入单层目录结构**：
  - 新增“目录”概念，话术可以归属到一个目录中。
  - 目录支持创建、重命名、删除。
- **双视图切换**：
  - **时间轴视图** (Timeline)：保留现有的按时间倒序展示，方便查看最新。
  - **目录视图** (Directory)：按目录分组展示，结构化管理。
- **话术卡片优化**：
  - 去除“标题”字段，仅展示内容预览。
  - 增加快捷操作按钮：复制、编辑、删除。
- **交互安全优化**：
  - 删除话术或目录时，增加二次确认弹窗。
  - 目录删除支持“级联删除”或“保留话术（移至未分类）”的选择。
- **后端 API 扩展**：
  - 支持目录的 CRUD 操作。
  - 支持话术关联目录。

## Capabilities

### New Capabilities

- `script-categories`: 核心目录管理能力，包括目录的增删改查、以及目录与话术的关联关系管理。

### Modified Capabilities

- `script-management`: 扩充现有话术管理能力，支持目录归属、视图过滤以及优化的删除逻辑。

## Impact

- **Database**:
  - New Table: `Category`
  - Migration: `Script` table adds `category_id` FK.
- **Backend (Go)**:
  - New `CategoryService`.
  - Update `ScriptService` to handle `category_id`.
- **Frontend (Vue)**:
  - Major refactor of `ScriptsPage.vue`.
  - New components: `CategoryList`, `ScriptItem`, `CategoryDialog`.
  - Update `MainLayout` or page header to support view switching.
