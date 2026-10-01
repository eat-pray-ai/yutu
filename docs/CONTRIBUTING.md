# Contributing

Feel free to contribute to the project under these conventions:

- Commit messages should follow the [gitmoji](https://gitmoji.dev) convention.
- Follow the existing naming and project structure.

## Git Pre-Hooks

Install Git hooks to automate code hygiene, license headers, and test checks:

```shell
❯ ./scripts/pre-hooks.sh install
❯ ./scripts/pre-hooks.sh help
```

Here are some commands which may useful.

## Go Standard Toolchain

```shell
## upgrade all dependencies
❯ go get -u ./...

## run tests
❯ go test ./...
### with coverage report
❯ go test ./... -coverprofile=./coverage.out -coverpkg=./...
❯ go tool cover -html=coverage.out -o=coverage.html

## build binaries with GoReleaser (https://goreleaser.com/install)
❯ goreleaser build --clean --auto-snapshot

## verify binary commands, detect shorthands conflicts, etc.
❯ ./scripts/command-test.sh dist/yutu_darwin_arm64_v8.0/yutu-darwin-arm64

## script to install yutu
❯ ./scripts/install.sh
```

## Bazel Toolchain

```shell
## upgrade all dependencies
❯ bazel run @rules_go//go -- get -u ./...

## run tests
❯ bazel test //...
### with coverage report
❯ bazel coverage //...
❯ genhtml -o genhtml "$(bazel info output_path)/_coverage/_coverage_report.dat"

## build the binary
❯ bazel run //:gazelle  # (re)generate BUILD files
### update go.mod, go.sum, and use_repo in MODULE.bazel
❯ bazel run @rules_go//go -- mod tidy -v
❯ bazel mod tidy
❯ bazel build //:yutu   # build the binary for the current platform
❯ bazel build //...     # build all targets
❯ bazel build --platforms=@rules_go//go/toolchain:linux_amd64 //:yutu
❯ bazel cquery --output=files //:yutu-linux-amd64

## verify binary commands, detect shorthands conflicts, etc.
❯ ./scripts/command-test.sh "$(bazel info bazel-bin)/yutu_/yutu"

## script to install yutu
❯ ./scripts/install.sh
```
