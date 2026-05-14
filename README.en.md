# QuickTodo

[中文](README.md) | English

QuickTodo is a terminal Todo CLI written in Go for quickly managing local todo items.

The default executable name is `qtd`.

## Usage

```text
Usage: QuickTodo [options] [command] [todo...]

QuickTodo is a terminal Todo CLI for quickly managing local todo items. It supports adding, listing, updating, completing, removing, and clearing todos, with local configuration management.

Arguments:
  todo  Todo item

Options:
  -h, --help     display help for command
  -V, --version  output the version number

Commands:
  add [options] <todo...>       Add todo item
  rm <index...>                 Remove todo item
  mod [options] <index> [todo]  Modify todo item
  list [options]                Show todo items
  done <index...>               Complete todo item, equivalent to: mod <index> -d
  clear [options]               Clear todo items
  conf                          Configuration
```

Examples:

```sh
./qtd add "write README"
./qtd list
./qtd mod 0 "update README"
./qtd done 0
```

Running the root command without arguments is equivalent to `list`; passing text directly is equivalent to `add`.

![](README_files/example.jpg)