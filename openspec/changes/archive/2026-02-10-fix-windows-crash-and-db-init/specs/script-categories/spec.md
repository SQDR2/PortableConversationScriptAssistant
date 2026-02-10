## ADDED Requirements

### Requirement: 数据库不可用时目录接口不得崩溃

当数据库未就绪或初始化失败时，目录/分类管理相关接口 MUST 返回可诊断的错误，并且 MUST NOT 因空指针等异常导致应用崩溃。

#### Scenario: 数据库初始化失败后调用 ListCategories
- **WHEN** `db.InitDB` 初始化失败且前端调用 `ListCategories`
- **THEN** 后端返回错误信息并记录日志，同时返回空数组（而非 `null`）作为列表结果
