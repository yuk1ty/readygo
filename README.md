# readygo

The following multi-line command

```shell
mkdir example
cd example
go mod init github.com/yuk1ty/example
git init
touch main.go
# edit Go code with your favourite editors
```

is shortened as below by using `readygo`

```shell
readygo -p github.com/yuk1ty/example
```

## Getting started

Use `go install`.

```
go install github.com/yuk1ty/readygo@latest
```

## How to use

### Basic usage

`readygo` has the following options:

- `--dir-name` or `-n`: Name for the directory which is created by `readygo`. This option can be omitted.
- `--module-path` or `-p`: Module path for `go mod init` command.
- `--layout` or `-l`: Directory layout style. You can choose [Standard Go Project Layout](https://github.com/golang-standards/project-layout) style (`standard`) or empty style (`default`). This option can be omitted. The default value is `default`, creates an empty directory.

### Examples

```
readygo --help
Usage:
  readygo [flags]

Flags:
  -n, --dir-name string      Define the directory name of your project. This can be omitted. If you do so, the name will be extracted from its package name.
  -h, --help                 help for readygo
  -l, --layout default       Define your project layout. You can choose default or `standard`. If you omit this option, the value becomes `default`. (default "default")
  -p, --module-path string   Define your module path. This is used for go mod init [module path].
```

#### Full

```shell
readygo --dir-name newprj --module-path github.com/yuk1ty/example
```

generates

```shell
ls -a --tree --level 1 newprj
newprj
├── .git
├── .gitignore
├── go.mod
└── main.go
```

and its `go.mod` is to be as below:

```shell
module github.com/yuk1ty/example

go 1.18
```

#### With layout style

```shell
readygo -n example -p github.com/yuk1ty/example -l standard
```

generates

```shell
ls -a --tree --level 1 example
example
├── .git
├── .gitignore
├── cmd
├── go.mod
├── internal
├── main.go
└── pkg
```

#### Omit `--dir-name(-n)`

```shell
readygo -p github.com/yuk1ty/example
```

generates

```shell
ls -a --tree --level 1 example
example
├── .git
├── .gitignore
├── go.mod
└── main.go
```

The following one (illustrates the case that is not starting with `github.com/username`) works fine as well.

```shell
readygo -p example
```

generates

```shell
ls -a --tree --level 1 example
example
├── .git
├── .gitignore
├── go.mod
└── main.go
```
