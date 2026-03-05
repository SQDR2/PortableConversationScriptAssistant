## ADDED Requirements

### Requirement: 在文件管理器中定位媒体文件

系统 SHALL 提供"在文件夹中显示"操作，允许用户通过操作系统的文件管理器定位到指定的媒体文件（图片或视频）。此功能需在 Linux、macOS 和 Windows 三个平台上均可用。

#### Scenario: 用户通过图片右键菜单触发「在文件夹中显示」

- **WHEN** 用户右键单击话术卡片中的图片，并点击「在文件夹中显示」
- **THEN** 系统调用后端 `RevealInFileManager`，使用平台对应的文件管理器打开图片所在目录（macOS 高亮选中文件，Linux 打开目录，Windows 高亮选中文件）

#### Scenario: 用户通过视频操作按钮触发「在文件夹中显示」

- **WHEN** 用户点击话术卡片视频区域中的「在文件夹中显示」图标按钮
- **THEN** 系统调用后端 `RevealInFileManager`，在文件管理器中定位到对应的视频文件

#### Scenario: 文件不存在时给予错误提示

- **WHEN** 媒体文件对应的本地路径不存在或无法访问
- **THEN** 系统弹出错误通知，提示「文件不存在或已被删除」，不打开文件管理器

#### Scenario: 跨平台命令正确执行

- **WHEN** `RevealInFileManager` 被调用，且当前运行环境为 Linux
- **THEN** 后端执行 `xdg-open <目录路径>`，打开文件所在目录

- **WHEN** `RevealInFileManager` 被调用，且当前运行环境为 macOS
- **THEN** 后端执行 `open -R <文件路径>`，高亮该文件

- **WHEN** `RevealInFileManager` 被调用，且当前运行环境为 Windows
- **THEN** 后端执行 `explorer /select,<文件路径>`，高亮该文件
