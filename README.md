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
readygo -n example -p github.com/yuk1ty/example
```

## How to use

### Basic usage

`readygo` has the following options:

- `--name` or `-n`: Name for the directory which is created by `readygo`. This option can be omitted.
- `--pkg-name` or `-p`: Package name for `go mod init` command.
- `--style`: Directory layout style. You can choose [Standard Go Project Layout](https://github.com/golang-standards/project-layout) style (`standard`) or empty style (`default`). This option can be omitted. The default value is `default`, creates an empty directory.

### Examples

```shell
readygo --name example --pkg-name github.com/yuk1ty/example
```

generates

```shell
cd example
ls -la
drwxr-xr-x  - a14926 25 4 17:17 .git
.rw-r--r-- 42 a14926 25 4 17:17 go.mod
.rw-r--r-- 72 a14926 25 4 17:17 main.go
```

```shell
readygo -n example -p github.com/yuk1ty/example -s standard
```

generates

```shell
cd example
ls -la
drwxr-xr-x  - a14926 25 4 17:28 .git
drwxr-xr-x  - a14926 25 4 17:28 cmd
.rw-r--r-- 42 a14926 25 4 17:28 go.mod
drwxr-xr-x  - a14926 25 4 17:28 internal
.rw-r--r-- 72 a14926 25 4 17:28 main.go
drwxr-xr-x  - a14926 25 4 17:28 pkg
```
