# QuickTodo

QuickTodo 是一个使用 Go 编写的终端 Todo CLI。它支持在命令行中快速添加、查看、修改和删除待办事项，并将数据持久化到本地 JSON 文件。

## 功能状态

已实现：

- `add`：添加一个或多个待办事项。
- `list`：列出待办事项，并支持按完成状态、内容和 ID 范围过滤。
- `mod`：修改待办内容、完成状态和优先级。
- `done`：按 ID 将一个或多个待办事项标记为完成。
- `rm`：删除指定 ID 的待办事项。
- `clear`：清空全部待办事项。

## 构建

构建当前平台可执行文件：

```sh
go build -o qtd .
```

构建全部包：

```sh
go build ./...
```

直接运行：

```sh
go run .
```

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

## 开发

格式化：

```sh
gofmt -w .
```

运行全部测试：

```sh
go test ./...
```

运行单个包测试：

```sh
go test ./internal/data
go test ./internal/processor
```

运行单个测试：

```sh
go test ./internal/data -run TestLocalRepository
go test ./internal/processor -run TestAdd
```

静态检查和构建：

```sh
go vet ./...
go build ./...
```

## Release

GitHub Actions 会在符合 `v1.2.3` 格式的 tag 上创建 Release，并上传以下平台和架构的二进制产物：

- `qtd_darwin_amd64`
- `qtd_darwin_arm64`
- `qtd_linux_amd64`
- `qtd_linux_arm64`
- `qtd_windows_amd64.exe`
- `qtd_windows_arm64.exe`

示例：

```sh
git tag v0.1.0
git push origin v0.1.0
```

Release 构建会从 tag 注入应用版本号：`v0.1.0` 对应程序版本 `0.1.0`。本地调试构建默认版本为 `dev`，也可以手动注入：

```sh
go build -ldflags="-X todo_list/internal/app.APP_VERSION=0.1.0" -o qtd .
```
