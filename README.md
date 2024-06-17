# gosumfix

This is a command line tool that tries to automatically fix git conflicts in the go.mod and go.sum files.

If this experiment turns out to be useful, I plan to upstream it to the Go project as part of the `go` tool.

There is a tracking proposal for this feature at https://www.github.com/golang/go/issues/32485.

## Installation

```bash
go install github.com/mauri870/gosumfix/cmd/...@latest
```

This will install the binaries in your `$GOPATH/bin` directory.

It includes the following commands:

- `gosumfix`: The main command that tries to fix the conflicts in mod files.
- `gosumdriver`: This is a git merge driver that will automatically run `gosumfix` when you have conflicts in the go.sum file.

## Usage

When you have a conflict in the go.sum file you can run `gosumfix` to try to fix it. If the conflicting lines contain `replace` or `exclude` directives they need to be fixed manually.

```bash
gosumfix
```

If you install the git merge driver you can run `git merge` as usual and the driver will automatically run `gosumfix` when there are conflicts in the go.sum/go.mod files.

```bash
gosumdriver install
gosumdriver uninstall # To uninstall the driver
```
