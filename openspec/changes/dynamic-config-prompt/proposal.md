## Why

目前应用依赖静态的 `config.json` 文件进行翻译服务的鉴权。如果该文件内容为空（首次使用或用户误删），应用启动后翻译服务将不可用且无明确提示。为了提升用户体验，需在配置缺失时主动引导用户输入凭证，并支持动态保存配置，避免用户手动修改配置文件的繁琐操作。

## What Changes

- **配置管理升级**：
  - `config.json` 管理从“静态只读”升级为“动态读写”。
  - 后端增加配置保存接口，并支持热重载 TMT 客户端。
- **用户交互增强**：
  - 启动时自动检查翻译配置状态。
  - **UI新增**：在未配置时弹出模态对话框，强制或引导用户输入 `SecretId` 和 `SecretKey`。

## Capabilities

### New Capabilities

- `config-management`: 提供配置状态查询、验证、持久化以及服务热重载的能力。

### Modified Capabilities

- `translation-service`: 增加对动态配置变更的响应能力（依赖 `config-management` 保存配置后刷新 Client）。
- `ui-translation-bar`: 增加启动时的鉴权状态检查逻辑及鉴权输入弹窗交互。

## Impact

- **API**: 新增后端方法 `GetConfigStatus` 和 `UpdateCredentials`。
- **Frontend**: 直接影响 `MainLayout.vue` 的启动流程。
- **Security**: 涉及 `SecretKey` 的传输和本地存储（仍为明文存储，与现状保持一致）。
