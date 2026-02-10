## ADDED Requirements

### Requirement: Windows 平台窗口跟随长期稳定运行

当启用窗口跟随（已关联目标窗口且轮询运行）时，系统 MUST 确保应用在 Windows 平台上长期稳定运行，不得因 Windows 回调资源耗尽而崩溃退出（例如 `fatal error: too many callback functions`）。

在 Windows 平台上，窗口跟随的轮询热路径（`checkTarget()` 的周期性调用链）MUST NOT 以“每个轮询周期创建新的回调函数/回调桩”的方式工作。

#### Scenario: 关联目标窗口后持续运行不崩溃
- **WHEN** 用户关联目标窗口并保持窗口跟随运行至少 10 分钟
- **THEN** 应用持续运行且窗口跟随功能正常，不发生崩溃退出

#### Scenario: 轮询热路径不产生无限回调资源
- **WHEN** 窗口跟随轮询持续运行
- **THEN** 系统不会随轮询周期增长而持续创建新的 Windows 回调资源（回调数量应保持稳定，不随时间线性增长）
