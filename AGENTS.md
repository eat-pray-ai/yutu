# Agent Guidelines

Focus mostly on given tasks, avoid unnecessary complexity. Since developers may update code manually, codebase can be different from generated code. Ensure compatibility with existing structures and conventions.

Follow developers instruction first, don't update existing codes unless developer asked explicitly.

## Commands

- **Build**: `go build ./...` or `bazel build //...`
- **Test (All)**: `go test ./...` or `bazel test //...`
- **Test (Single)**: `go test -v ./pkg/path -run TestName`
- **Lint**: `go vet ./...`
- **Update Bazel**: `bazel run //:gazelle` (Run this after adding/removing files or imports)

## Code Style

- **Formatting**: Ensure code is formatted with `go fmt`.
- **Naming**: Use camelCase for multi-word packages, directories, and filenames (e.g., `channelBanner`, `i18nRegion`).
- **Headers**: Run `addlicense -c "eat-pray-ai & OpenWaygate" -f LICENSE -s=only -ignore "**/*.yml" -ignore "**/*.yaml" -ignore "**/*.bazel" .` to add headers.
- **Imports**: Group standard library first, then third-party, then local imports.
- **Errors**: explicit error handling `if err != nil`.
- **Commits**: Follow [gitmoji](https://gitmoji.dev) convention.

## Structure

- **CLI Commands**: Located in `cmd/<resource>/`.
- **Core Logic**: Located in `pkg/<resource>/`.
- **Tests**: Co-located with code in `pkg/` as `_test.go`.
