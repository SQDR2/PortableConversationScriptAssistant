## 1. Windows 窗口跟随崩溃修复（callback 泄漏）

- [x] 1.1 在 `backend/utils/window_windows.go` 中识别所有 `syscall.NewCallback` 调用点，确认哪些会被轮询热路径触发
- [x] 1.2 为 `EnumWindows` 回调实现“可复用”方案（进程级单例 callback 或彻底移除热路径枚举），确保不会随时间持续创建 callback
- [x] 1.3 新增基于 hwnd 的窗口定位/尺寸设置函数（例如 `ForceMoveResizeWindow`），并在轮询热路径中改用 hwnd 而非按标题查找
- [x] 1.4 优化 sidekick 自身 hwnd 的获取与缓存策略：首次获取、失效检测、低频重取（不得引入 callback 泄漏）
- [ ] 1.5 回归验证：窗口跟随开启后持续运行 10 分钟以上不崩溃；最小化/恢复/对齐逻辑保持与现有一致

## 2. 数据库初始化在 CGO_DISABLED 环境下可用

- [x] 2.1 确认当前 sqlite driver 在 `CGO_ENABLED=0` 下的失败路径与错误信息，并将其纳入启动错误日志
- [x] 2.2 将 sqlite driver 调整为不依赖 cgo 的实现，确保 Windows 下 `CGO_ENABLED=0` 构建仍能 `db.InitDB` 成功
- [x] 2.3 验证 WAL/外键等 PRAGMA 在新 driver 下的兼容性；不支持时提供等价处理或降级策略

## 3. 服务层 DB 不可用时的稳定性与可诊断性

- [x] 3.1 为 `CategoryService`/`ScriptService` 添加 DB 就绪检查：当 `db.DB == nil` 时返回可诊断错误且不触发 nil pointer
- [x] 3.2 确保列表接口在错误场景返回空数组（而非 `null`）并记录日志：`ListScripts`、`ListCategories`
- [x] 3.3 回归验证：在模拟 DB 初始化失败场景下，前端调用脚本/分类接口不会导致进程崩溃，且能看到明确错误信息

## 4. 错误日志补强（db_init_error）

- [x] 4.1 将 `db.InitDB` 失败写入 crash-logging 规范的独立日志文件，错误类型为 `db_init_error`，包含错误信息与堆栈跟踪
- [x] 4.2 验证日志文件名格式与内容字段符合现有 crash-logging 规范
