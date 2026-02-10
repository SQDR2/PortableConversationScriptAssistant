## ADDED Requirements

### Requirement: 数据库初始化失败写入明确错误类型日志

当数据库初始化（`db.InitDB`）失败时，系统 MUST 按 crash-logging 的规范生成独立错误日志文件，并使用明确的错误类型标签 `db_init_error`，以便用户环境问题可被快速定位。

#### Scenario: 数据库初始化失败生成 db_init_error 日志
- **WHEN** `db.InitDB` 返回错误
- **THEN** 系统在 `logs/` 目录下创建一个错误类型为 `db_init_error` 的日志文件，内容包含错误信息与堆栈跟踪
