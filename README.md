# Sidekick (话术助手)

Sidekick 是一款基于 Wails 构建的桌面辅助工具，旨在为日常办公和即时通讯提供快捷的话术管理与窗口联动功能。它可以“吸附”在任何指定的应用程序窗口旁，方便用户快速检索并复制预设的话术。

## 核心功能

- **窗口联动**：Sidekick 窗口可以自动跟随目标应用窗口移动，实现“贴边”显示。
- **关联提示**：侧边栏展示当前已关联的目标应用，便于确认跟随对象。
- **智能化检索**：内置 SQLite FTS5 (全文搜索)，支持快速搜索成千上万条话术。
- **话术管理**：支持话术的增删改查、标签分类，并支持通过文本文件批量导入。
- **磨砂玻璃特效**：采用现代化的 Glassmorphism 设计，界面美观且不遮挡背景。
- **实时翻译**：集成腾讯云 TMT 翻译引擎，支持多语言实时对照翻译。

## 技术栈

### 后端 (Go)

- **Wails v2**：构建轻量级原生桌面应用。
- **GORM**：强大的 Go ORM 库，管理 SQLite 数据库。
- **SQLite FTS5**：实现高性能的全文本检索。
- **XGB / EWMH (Linux)** & **Win32 API (Windows)**：底层窗口同步与控制。
- **腾讯云 SDK**：集成 TMT 服务，实现高精度的文本翻译。

### 前端 (TypeScript)

- **Vue 3**：渐进式 JavaScript 框架。
- **Quasar Framework**：高性能的移动端和桌面端 UI 组件库。
- **Vite**：下一代前端构建工具。

## 如何使用

### 环境要求

- [Go](https://go.dev/dl/) 1.18+
- [Node.js](https://nodejs.org/) & [pnpm](https://pnpm.io/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)

### 本地开发

1. 克隆项目：
   ```bash
   git clone <project-url>
   cd PortableConversationScriptAssistant
   ```
2. 运行开发模式（支持热重载）：
   ```bash
   wails dev
   ```

### 编译打包

执行以下命令生成对应平台的安装包：

```bash
wails build
```

## 配置说明

### 翻译服务配置

Sidekick 的翻译功能依赖腾讯云文本翻译 (TMT) 服务。

1. **自动配置**：
   - 首次启动应用时，若检测到未配置密钥，会自动弹出配置窗口。
   - 也可以通过侧边栏中的 **“翻译配置”** 选项随时进行修改。

2. **手动配置**：
   - 在应用根目录下的 `config.json` 文件中，填写以下信息：
     ```json
     {
       "tencent_cloud": {
         "secret_id": "你的 SecretId",
         "secret_key": "你的 SecretKey",
         "region": "ap-guangzhou",
         "project_id": 0
       }
     }
     ```
   - 获取密钥：前往 [腾讯云访问管理控制台](https://console.cloud.tencent.com/cam/capi) 获取 SecretId 和 SecretKey。

## 贡献

欢迎提交 Issue 或 Pull Request 来改进 Sidekick！

## 许可

本项目采用 MIT 许可证。
