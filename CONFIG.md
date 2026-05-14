# 内置配置项

本文档说明 QuickTodo 当前支持的内置配置项。配置可通过 `qtd conf` 命令查看和修改。

## 命令

查看配置：

```sh
qtd conf
qtd conf list
qtd conf list REPOSITORY_NAME REPOSITORY_LOCAL_TABLE DISPLAY_TABLE_MODE
```

修改配置：

```sh
qtd conf set <key> <value>
```

删除配置：

```sh
qtd conf del <key...>
```

查看配置历史：

```sh
qtd conf list --history
qtd conf list <key> --history
```

## 配置项

| Key | 默认值 | 可选值 | 说明 |
| --- | --- | --- | --- |
| `REPOSITORY_NAME` | `local` | `local` | Todo 数据仓储。当前仅内置本地 JSON 仓储。 |
| `REPOSITORY_LOCAL_TABLE` | `default` | 任意非空字符串 | 本地仓储的数据表名。不同表名会对应不同的 Todo 数据文件。 |
| `DISPLAY_TABLE_MODE` | `table` | `table`、`markdown` | 表格输出模式。`table` 使用终端表格，`markdown` 使用 Markdown 表格。 |

## 删除行为

`conf del` 对内置配置项和自定义配置项有不同处理：

- 内置配置项：恢复为默认值，并将历史记录重置为默认值。
- 自定义配置项：从当前配置和历史记录中移除。

示例：

```sh
qtd conf del DISPLAY_TABLE_MODE
```

执行后，`DISPLAY_TABLE_MODE` 会恢复为 `table`。

## 历史记录

使用 `conf set` 修改配置时，QuickTodo 会记录该 key 设置过的不同值。

- 同一个 key 下重复设置相同 value 时，只保存一次。
- 每个 key 最多保留最近 20 个不同值。
- `conf list --history` 会在表格中显示历史记录。
