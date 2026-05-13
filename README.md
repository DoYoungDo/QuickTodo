# QuickTodo

QuickTodo 是一个使用 Go 编写的终端 Todo CLI。它支持在命令行中快速添加、查看、修改和删除待办事项，并将数据持久化到本地 JSON 文件。

## 使用示例

添加待办：

```sh
./qtd add "write README"
./qtd add "ship release" "add tests"
```

添加时标记为完成：

```sh
./qtd add "already done" -d
```

列出待办：

```sh
./qtd list
```

按内容过滤：

```sh
./qtd list -f readme -i
```

按 ID 范围过滤：

```sh
./qtd list -b 0 -e 3
```

修改内容：

```sh
./qtd mod 0 "update README"
```

追加或前插内容：

```sh
./qtd mod 0 " today" --append
./qtd mod 0 "urgent: " --insert
```

标记完成或设置优先级：

```sh
./qtd mod 0 -d
./qtd mod 0 -p 3
./qtd done 0
./qtd done 0 1 2
```

删除待办：

```sh
./qtd rm 0
./qtd rm 0 1 2
```

清空全部待办：

```sh
./qtd clear
./qtd clear -f
```

`clear` 默认会先提示确认；使用 `-f` 或 `--force` 会跳过确认并直接清空。

根命令没有参数时等价于 `list`；直接传入文本时等价于 `add`：

```sh
./qtd
./qtd "quick task"
```

## 数据存储

默认数据文件位于：

```text
os.UserConfigDir()/QuickTodo/todos/default.json
```

该路径由运行系统决定，例如 macOS 通常位于用户配置目录下。手工验证命令会写入真实数据文件，测试代码使用临时目录隔离。
