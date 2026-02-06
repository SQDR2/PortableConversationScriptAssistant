## MODIFIED Requirements

### Requirement: 话术数据结构扩展

话术模型必须包含一个图片字段，支持存储多个图片引用。

#### Scenario: 查询话术列表

- **WHEN** 前端请求话术列表
- **THEN** 后端返回的每个 `Script` 对象应包含 `images` 数组（即图片的本地路径列表）

### Requirement: 编辑器图片支持

话术编辑器对话框必须提供界面让用户管理图片。

#### Scenario: 编辑现有话术

- **WHEN** 用户编辑一个已有话术并添加新图片
- **THEN** 保存后，新的图片路径应持久化到该话术记录中
