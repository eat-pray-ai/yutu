# Contributing

yutu is a cli tool built using the [cobra](https://github.com/spf13/cobra). Feel free to contribute to the project under these conventions:

- Commit messages should follow the [gitmoji](https://gitmoji.dev) convention.
- Follow the existing naming and project structure.

Here are some commands which may useful.

```shell
# upgrade all dependencies
❯ go get -u ./...

# run tests
## unit tests
❯ go test ./...
## or with bazel
❯ bazel test //...
## verify binary commands, detect shorthands conflicts, etc.
❯ ./scripts/command-test.sh

# build the binary
## install goreleaser: https://goreleaser.com/install
❯ GITHUB_REPOSITORY=eat-pray-ai/yutu goreleaser build --clean --auto-snapshot
# or with bazel
❯ bazel run //:gazelle  # (re)generate BUILD files
## update go.mod, go.sum, and use_repo in MODULE.bazel
❯ bazel run @rules_go//go -- mod tidy -v 
❯ bazel build //:yutu   # build the binary for the current platform
❯ bazel build //...     # build all targets
❯ bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //:yutu

# script to install yutu
❯ ./scripts/install.sh
```
