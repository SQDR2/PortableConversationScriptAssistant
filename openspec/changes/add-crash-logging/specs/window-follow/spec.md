## MODIFIED Requirements

### Requirement: 遮挡状态下自动恢复可见

即使目标应用位置未发生变化，只要目标应用处于激活状态，系统 SHALL 定期检查并确保助手窗口处于顶层可见状态，防止助手窗口被其他操作（如点击任务栏打开新窗口）遮挡。

系统 MUST 在 `checkTarget()` 方法中安装 panic 恢复机制，确保窗口跟随轮询不会因为任何异常而终止。当 panic 发生时，系统 MUST 将错误信息写入日志文件并继续运行。

#### Scenario: 助手窗口被新打开的窗口遮挡

- **WHEN** 目标窗口处于激活状态，且有新窗口弹出并遮挡了助手窗口
- **THEN** 在下一个同步周期，系统将助手窗口重新带回最前端显示

#### Scenario: checkTarget 中发生 panic 后轮询继续

- **WHEN** `checkTarget()` 轮询过程中发生 panic
- **THEN** 系统捕获 panic 并写入日志文件，下一个轮询周期正常继续，窗口跟随功能不受影响
