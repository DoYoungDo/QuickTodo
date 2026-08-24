# QuickTodo：一个终端里的 Todo CLI

QuickTodo 是一个用 Go 编写的终端 Todo 工具，默认命令名是 `qtd`。它的用途很直接：在命令行里快速添加、查看、修改、完成和删除待办事项。

如果你的日常工作主要在终端里完成，`qtd` 可以把临时任务记录、任务状态更新和待办列表查看都放在当前工作流里，不需要切换到单独的应用。

## 安装

GitHub 仓库地址：<https://github.com/DoYoungDo/QuickTodo>

预编译二进制下载地址：<https://github.com/DoYoungDo/QuickTodo/releases>

如果从源码构建，需要先安装 Go。当前项目的 `go.mod` 声明版本为 Go 1.26。

```sh
git clone https://github.com/DoYoungDo/QuickTodo.git
cd QuickTodo
go build -o qtd .
```

构建完成后，可以直接在当前目录运行：

```sh
./qtd --version
./qtd list
```

如果希望在任意目录使用 `qtd`，可以把二进制放到 PATH 中的目录。例如：

```sh
mv qtd /usr/local/bin/qtd
qtd --version
```

如果通过 GitHub Release 页面下载预编译二进制，可以选择对应系统和架构的文件，改名为 `qtd` 后放入 PATH。

## 使用效果

![BQACAgUAAyEGAASHRsPbAAEZuY9qi7o47psjd_U3H5Z20bu4_OxG2AACkCUAAg-uYVTMOvj45u6ZBT0E.jpg](https://img.remit.ee/api/file/BQACAgUAAyEGAASHRsPbAAEZuY9qi7o47psjd_U3H5Z20bu4_OxG2AACkCUAAg-uYVTMOvj45u6ZBT0E.jpg)

![BQACAgUAAyEGAASHRsPbAAEZuZxqi7rdSEmVSL7Z-5twjlI_JeW6SwACnSUAAg-uYVRa9E4zjHTGbz0E.jpg](https://img.remit.ee/api/file/BQACAgUAAyEGAASHRsPbAAEZuZxqi7rdSEmVSL7Z-5twjlI_JeW6SwACnSUAAg-uYVRa9E4zjHTGbz0E.jpg)

## 基本使用

查看帮助信息：

```sh
qtd -h
qtd --help
```

如果只想查看某个子命令的用法，可以在子命令后使用 `-h`：

```sh
qtd add -h
qtd mod -h
qtd list -h
qtd conf -h
```

帮助信息会列出当前命令支持的参数、选项和子命令。比如 `list -h` 可以查看过滤内容、忽略大小写、ID 范围等选项；`mod -h` 可以查看追加、插入、完成状态和优先级相关选项。

添加一个待办事项：

```sh
qtd add "write README"
```

查看待办列表：

```sh
qtd list
```

根命令没有参数时，默认等价于 `list`：

```sh
qtd
```

直接传入文本时，默认等价于 `add`：

```sh
qtd "review pull request"
```

把某个待办标记为完成：

```sh
qtd done 0
```

修改待办内容：

```sh
qtd mod 0 "update README"
```

删除待办：

```sh
qtd rm 1
```

移动待办顺序：

```sh
qtd mv 2 0
```

清空所有待办：

```sh
qtd clear
```

如果不需要确认提示，可以使用：

```sh
qtd clear -f
```

## 常用操作

添加时直接标记为完成：

```sh
qtd add "ship release" -d
```

只查看已完成事项：

```sh
qtd list -d
```

按内容过滤：

```sh
qtd list -f "readme"
```

忽略大小写过滤：

```sh
qtd list -f "readme" -i
```

按 ID 范围查看：

```sh
qtd list -b 0 -e 3
```

追加内容：

```sh
qtd mod 0 " - urgent" --append
```

在原内容前插入文本：

```sh
qtd mod 0 "NOW: " --insert
```

设置优先级：

```sh
qtd mod 0 -p 2
```

## 配置

QuickTodo 内置本地 JSON 仓储，配置也保存在本地。查看配置：

```sh
qtd conf
qtd conf list
```

修改配置：

```sh
qtd conf set DISPLAY_TABLE_MODE markdown
```

删除配置：

```sh
qtd conf del DISPLAY_TABLE_MODE
```

查看配置历史：

```sh
qtd conf list --history
```

`DISPLAY_TABLE_MODE` 支持两种输出模式：

- `table`：默认终端表格。
- `markdown`：Markdown 表格，适合复制到文档、Issue 或日报中。

## 适合的使用场景

QuickTodo 比较适合这些场景：

- 开发过程中临时记录待办事项。
- 在终端里快速维护个人任务列表。
- 用本地文件保存 Todo，不依赖账号或远程服务。
- 需要把待办列表以 Markdown 表格形式放进文档或工作记录。

它不是完整的项目管理系统，也不处理多人协作、看板、提醒和同步。它更适合承担一个明确的小角色：让 Todo 操作留在命令行里，并尽量减少打断。
