# TODO

本文件记录当前仓库中尚未完成的功能和后续改进方向。

## 功能增强

- [ ] 支持按优先级排序和过滤。关联：`internal/processor/list.go`、`internal/ui/table.go`
- [ ] 支持截止日期或计划时间字段。关联：`internal/data/repository.go`
- [ ] 支持多清单或 profile。关联：`internal/data/local_repository.go`、`internal/config/config.go`
- [ ] 完善 `mod -p` 的优先级范围校验。关联：`internal/processor/mod.go`

## TUI 原型

- [ ] 明确是否正式接入 `gocui` TUI 主流程。关联：`internal/ui/main_window.go`
- [ ] 完善 Details 视图，展示选中 Todo 的内容、状态、优先级和时间。关联：`internal/ui/main_window.go`
- [ ] 清理 TUI 中的坐标调试输出，改为正式状态展示。关联：`internal/ui/main_window.go`

## 配置与持久化

- [ ] 补全配置模块，支持自定义数据文件路径、默认清单名等选项。关联：`internal/config/config.go`
- [ ] 改进仓储初始化错误处理，减少 `panic`，向 CLI 层返回明确错误。关联：`internal/data/local_repository.go`
- [ ] 为损坏 JSON 数据提供更友好的错误提示或恢复流程。关联：`internal/data/local_repository.go`

## 测试与工程

- [ ] 为 `main.go` 的命令注册和根命令行为补集成测试。
- [ ] 为 `internal/ui/table.go` 补充更细粒度的渲染或过滤测试。
- [ ] 扩充 `.gitignore`，忽略构建产物和本地调试数据。
- [ ] 视命令数量增长情况，增加 `Makefile`、`Taskfile` 或 `justfile` 统一常用开发命令。
