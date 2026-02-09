## ADDED Requirements

### Requirement: 应用启动时初始化日志目录

系统 MUST 在应用启动时（早于所有其他组件初始化之前）在运行目录下创建 `logs/` 子目录。如果目录已存在，系统 MUST 跳过创建。如果目录创建失败，系统 MUST 将错误输出到 `stderr`，但不得阻止应用启动。

#### Scenario: 首次启动创建日志目录

- **WHEN** 应用首次启动且运行目录下不存在 `logs/` 目录
- **THEN** 系统创建 `logs/` 目录

#### Scenario: 日志目录已存在

- **WHEN** 应用启动且运行目录下已存在 `logs/` 目录
- **THEN** 系统跳过目录创建，正常启动

#### Scenario: 日志目录创建失败

- **WHEN** 应用启动且运行目录无写入权限导致 `logs/` 创建失败
- **THEN** 系统将错误信息输出到 `stderr`，应用继续正常启动

### Requirement: 错误事件生成独立日志文件

系统 MUST 在每次记录错误时，在 `logs/` 目录下创建一个独立的日志文件。文件名格式 MUST 为 `{YYYY-MM-DD_HH-mm-ss}_{错误类型}.log`，其中时间使用本地时间。

日志文件内容 MUST 包含以下信息：
- 精确时间戳
- 错误类型标签
- 错误信息
- 完整的堆栈跟踪（stack trace）

#### Scenario: 记录一个错误事件

- **WHEN** 系统调用日志写入函数记录一个类型为 `window_error` 的错误
- **THEN** 系统在 `logs/` 目录下创建文件如 `2026-02-09_10-30-45_window_error.log`，内容包含时间戳、错误信息和堆栈跟踪

#### Scenario: 同一秒内发生多个错误

- **WHEN** 在同一秒内发生两个不同类型的错误
- **THEN** 系统生成两个独立的日志文件，文件名因错误类型不同而不重复

### Requirement: main 函数全局 panic 恢复

系统 MUST 在 `main()` 函数中安装 `defer recover()` 机制。当发生未捕获的 panic 时，系统 MUST 将 panic 信息和完整堆栈跟踪写入日志文件后，再退出应用。

#### Scenario: main 函数中发生 panic

- **WHEN** 应用运行期间 main goroutine 中发生 panic
- **THEN** 系统捕获 panic，将错误信息和堆栈写入 `logs/` 目录下的日志文件，然后退出

### Requirement: checkTarget goroutine panic 恢复

系统 MUST 在 `checkTarget()` 方法内部安装 `defer recover()` 机制。当 `checkTarget()` 中发生 panic 时，系统 MUST 将 panic 信息写入日志文件，但不得终止轮询 goroutine 的继续运行。

#### Scenario: checkTarget 中发生 panic

- **WHEN** `checkTarget()` 轮询过程中发生 panic（如 Windows API 调用异常）
- **THEN** 系统捕获 panic，将错误信息写入日志文件，当前轮询周期结束，下一个轮询周期正常继续

#### Scenario: checkTarget 连续多次 panic

- **WHEN** `checkTarget()` 在连续多个轮询周期中发生 panic
- **THEN** 系统为每次 panic 生成独立的日志文件，轮询 goroutine 始终保持运行

### Requirement: 启动阶段错误日志记录

系统 MUST 在数据库初始化 (`db.InitDB`) 和各服务 `Startup` 过程中的错误发生时，将错误信息写入日志文件。

#### Scenario: 数据库初始化失败

- **WHEN** `db.InitDB` 返回错误
- **THEN** 系统将错误详情写入日志文件
