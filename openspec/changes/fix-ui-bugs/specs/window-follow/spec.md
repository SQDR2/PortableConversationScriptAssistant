## ADDED Requirements

### Requirement: 首次对齐前 MUST 解析助手窗口句柄

系统在执行窗口位置和尺寸同步时，MUST 先确保已获取助手窗口的原生句柄 (`sidekickHWID`)。若助手窗口句柄尚未解析，MUST 跳过当轮的 `ForceMoveResizeWindow` 调用，等待下一个轮询周期重新尝试，避免使用错误的默认宽度值。

#### Scenario: sidekickHWID 未就绪时不执行 resize

- **WHEN** 轮询周期内 `sidekickHWID` 为空
- **THEN** 系统跳过当轮的窗口位置/尺寸调整，不调用 `SetWindowPos`

#### Scenario: sidekickHWID 解析成功后正常对齐

- **WHEN** 下一轮轮询成功解析到 `sidekickHWID`
- **THEN** 系统使用 `GetWindowPhysicalWidth` 获取真实宽度并执行精确对齐
