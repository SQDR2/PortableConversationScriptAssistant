## Context

目前应用依赖 `backend/services/translation_service.go` 在启动时（Startup）读取 `config.json`。配置逻辑是单向、静态的。如果配置缺失，服务仅记录日志并初始化失败，导致前端无法使用翻译功能且没有任何用户反馈。我们需要建立一个前端到后端的配置写入通道，并支持热更新。

## Goals / Non-Goals

**Goals:**

- 实现后端配置写入能力 (`utils.SaveConfig`)。
- 实现 `TranslationService` 的状态查询 (`IsConfigured`) 和热更新 (`UpdateCredentials`)。
- 实现前端启动时的鉴权状态检查。
- 提供友好的配置输入界面（Dialog）。

**Non-Goals:**

- 复杂的权限管理或多用户配置。
- 配置文件加密存储（保持当前明文 JSON 格式）。
- 非腾讯云 TMT 服务的支持。

## Decisions

### 1. 配置持久化策略

- **Decision**: 在 `utils` 包中增加 `SaveConfig` 函数，使用 `encoding/json` 将内存中的配置结构体序列化并覆盖写入 `config.json`。
- **Rationale**: 简单直接，与现有的 `LoadConfig` 保持对称。
- **Alternative**: 使用 viper 等配置库，但对于目前的单文件配置来说过于重型。

### 2. 服务状态管理

- **Decision**: 在 `TranslationService` 中维护最新的 `AppConfig` 指针。`UpdateCredentials` 方法直接修改该指针指向的配置，并触发 `SaveConfig` 和 `tmt.NewClient`。
- **Rationale**: 避免重启应用即可生效，提升体验。

### 3. 前端交互入口

- **Decision**: 在 `MainLayout.vue` 的 `onMounted` 生命周期钩子中进行检查。
- **Rationale**: `MainLayout` 是应用的主容器，确保无论进入哪个页面（如果有路由的话），只要加载了主布局就会检查。这比在 `TranslationBar` 组件内部检查更稳健，防止因组件条件渲染（`v-if`）导致检查被跳过。

## Risks / Trade-offs

- **Risk**: 并发写入 `config.json`。
  - _Mitigation_: 目前应用是单实例且用户操作是串行的（Dialog 模态），风险极低。
- **Risk**: 用户输入错误的 Key/Id。
  - _Mitigation_: 后端会在 `UpdateCredentials` 中尝试初始化 Client，虽然仅仅是初始化可能不会立即报错（除非格式极度错误）。真正的校验发生在第一次翻译调用时。为了简化，可以在 Update 时尝试做一次简单的 API 调用（如查询支持语言）来验证，但本设计暂定仅保存并初始化 Client，API 错误由翻译调用时的错误处理机制覆盖。
