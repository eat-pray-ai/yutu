# Contributing

yutu is a mcp server & cli tool built using the [go-sdk](https://github.com/modelcontextprotocol/go-sdk) and [cobra](https://github.com/spf13/cobra). Feel free to contribute to the project under these conventions:

- Commit messages should follow the [gitmoji](https://gitmoji.dev) convention.
- Follow the existing naming and project structure.

Here are some commands which may useful.

```shell
# go standard toolchain

## upgrade all dependencies
❯ go get -u ./...

## run tests
### unit tests
❯ go test ./... -coverprofile=./cover.out -coverpkg="$(go list || go list -m | head -1)/pkg/..."
❯ go tool cover -html=cover.out -o=cover.html

## build the binary
### install goreleaser: https://goreleaser.com/install
❯ goreleaser build --clean --auto-snapshot

❯ tree dist
dist
├── artifacts.json
├── config.yaml
├── metadata.json
├── yutu_darwin_amd64_v1
│   └── yutu-darwin-amd64
├── yutu_darwin_arm64_v8.0
│   └── yutu-darwin-arm64
├── yutu_linux_amd64_v1
│   └── yutu-linux-amd64
├── yutu_linux_arm64_v8.0
│   └── yutu-linux-arm64
├── yutu_windows_amd64_v1
│   └── yutu-windows-amd64.exe
└── yutu_windows_arm64_v8.0
    └── yutu-windows-arm64.exe

7 directories, 9 files

## verify binary commands, detect shorthands conflicts, etc.
❯ ./scripts/command-test.sh dist/yutu_darwin_arm64_v8.0/yutu-darwin-arm64

## script to install yutu
❯ ./scripts/install.sh
```

```shell
# bazel toolchain

## upgrade all dependencies
❯ bazel run @rules_go//go -- get -u ./...

## run tests
### unit tests
❯ bazel test //...
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

❯ tree -L 1 bazel-bin/yutu_
bazel-bin/yutu_
├── yutu
├── yutu-0.params
├── yutu.repo_mapping
├── yutu.runfiles
└── yutu.runfiles_manifest

2 directories, 4 files

## verify binary commands, detect shorthands conflicts, etc.
❯ ./scripts/command-test.sh bazel-bin/yutu_/yutu

## script to install yutu
❯ ./scripts/install.sh
```
