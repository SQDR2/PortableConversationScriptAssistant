## ADDED Requirements

### Requirement: 无 cgo 构建下数据库仍可用

系统 MUST 支持在 `CGO_ENABLED=0` 的构建环境下正常初始化并使用本地 sqlite 数据库，以保证脚本管理能力在常见的 Windows/CI 构建方式下可用。

#### Scenario: CGO_ENABLED=0 构建启动后脚本列表可用
- **WHEN** 应用以 `CGO_ENABLED=0` 构建并启动完成
- **THEN** 数据库初始化成功，脚本列表查询可正常返回数据

### Requirement: 数据库不可用时脚本接口不得崩溃

当数据库未就绪或初始化失败时，脚本管理相关接口 MUST 返回可诊断的错误，并且 MUST NOT 因空指针等异常导致应用崩溃。

#### Scenario: 数据库初始化失败后调用 ListScripts
- **WHEN** `db.InitDB` 初始化失败且前端调用 `ListScripts`
- **THEN** 后端返回错误信息并记录日志，同时返回空数组（而非 `null`）作为列表结果
