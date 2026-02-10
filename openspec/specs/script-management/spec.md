## MODIFIED Requirements

### Requirement: 话术数据结构扩展

话术模型必须包含一个媒体附件字段（`images`），支持存储多个图片和视频引用。前端通过文件扩展名区分渲染方式。

#### Scenario: 查询话术列表

- **WHEN** 前端请求话术列表
- **THEN** 后端返回的每个 `Script` 对象应包含 `images` 数组（包含图片和视频的本地路径列表）

#### Scenario: 前端渲染媒体附件

- **WHEN** 前端遍历 `images` 数组渲染附件
- **THEN** 对于扩展名为 `.mp4` 或 `.webm` 的路径使用 `<video>` 标签渲染，其他扩展名使用 `<img>` / `<q-img>` 标签渲染

### Requirement: 编辑器图片支持

话术编辑器对话框必须提供界面让用户管理图片和视频附件。图片和视频共享上传入口和附件数量上限（10 个）。

#### Scenario: 编辑现有话术并添加视频

- **WHEN** 用户编辑一个已有话术并添加新的视频文件
- **THEN** 保存后，新的视频路径应持久化到该话术记录的 `images` 字段中

#### Scenario: 附件数量上限

- **WHEN** 话术已有 10 个媒体附件（图片+视频合计）
- **THEN** 上传入口隐藏，用户无法继续添加

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
