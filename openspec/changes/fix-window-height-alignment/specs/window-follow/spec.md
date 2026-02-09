## MODIFIED Requirements

### Requirement: 关联后立即贴右侧并同步高度

助手窗口在用户关联目标窗口后，系统 MUST 立即将助手窗口定位到目标窗口右侧外边缘，并将助手窗口高度同步为目标窗口高度（目标窗口视觉高度）。在 Windows 平台上，MUST 使用 DWM API (Extended Frame Bounds) 获取目标窗口的视觉矩形。

系统 MUST 确保助手窗口的视觉边界（DWM Extended Frame Bounds）与目标窗口的视觉边界在顶部和底部精确对齐。在 Windows 平台上，系统 MUST 计算助手窗口自身的 DWM frame bounds 与 GetWindowRect 之间的偏移量（阴影补偿），并将其纳入 `SetWindowPos` 的坐标和尺寸参数中，以实现像素级精确对齐。

系统 MUST 在 Windows 平台上使用原生 Win32 API (`SetWindowPos`) 直接设置窗口位置和大小，而非依赖框架层的窗口管理 API，以避免 DPI 缩放和监视器偏移带来的误差。

该对齐 MUST 在以下 DPI 缩放比例下均能正确工作：100%、125%、150%、175%、200%。

#### Scenario: 首次关联目标窗口

- **WHEN** 用户选择并关联一个目标窗口
- **THEN** 助手窗口被放置在目标窗口右侧外边缘且高度与目标窗口视觉边界一致（顶部和底部均精确对齐）

#### Scenario: 在 150% DPI 缩放下关联目标窗口

- **WHEN** 系统 DPI 缩放设置为 150%，用户选择并关联一个目标窗口
- **THEN** 助手窗口的视觉顶部和底部与目标窗口精确对齐，不因 DPI 缩放而产生偏差

#### Scenario: 目标窗口调整大小后重新同步

- **WHEN** 目标窗口高度发生变化
- **THEN** 助手窗口的视觉底部与目标窗口的视觉底部保持精确对齐
