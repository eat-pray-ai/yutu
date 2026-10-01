# Pre-Release Checklist

Follow this checklist before creating and pushing a release tag to ensure all assets, dependencies, tests, and documentation are properly synchronized.

## 1. Code & Resource Registration

- [ ] **Verify Resource Registration**:
  Every resource package under `pkg/` and `cmd/` must be blank-imported in:
  - `main.go` (`_ "github.com/eat-pray-ai/yutu/cmd/<resource>"`)
  - `internal/tools/skillgen/main.go`
  - `internal/tools/cmdtestgen/main.go`

- [ ] **Run Code Generators**:
  ```shell
  # Regenerate agent instructions and skill
  go run ./internal/tools/skillgen

  # Regenerate smoke tests
  go run ./internal/tools/cmdtestgen
  ```

- [ ] **Apply License Headers**:
  ```shell
  addlicense \
    -c "eat-pray-ai & OpenWaygate" \
    -f LICENSE \
    -s=only \
    -y $(date +%Y) \
    -ignore "**/*.yml" \
    -ignore "**/*.yaml" \
    -ignore "**/*.bazel" .
  ```

## 2. Dependency & Build System Sync

GitHub Actions CI runs `git diff --exit-code` after syncing dependencies. Ensure Go and Bazel modules are completely in sync:

```shell
go mod tidy
bazel run @rules_go//go -- mod tidy -v
bazel run //:gazelle
bazel mod tidy
```

## 3. Testing & Quality Checks

Run the same test and lint suites that CI will execute:

```shell
# Linter
golangci-lint run

# Go unit tests
go test ./...

# Bazel tests
bazel test //...

# CLI smoke tests (verifies flags, help, and catches shorthand conflicts)
./scripts/command-test.sh
```

## 4. Documentation & Metadata

- [ ] **Documentation Alignment**:
  - `README.md` & `README_zh.md`: CLI flags, command list, environment variables, MCP configs.
  - `docs/FEATURES.md`: Any new resources or command verbs documented.
  - `server.json`: Check MCP registry metadata if new tools or capabilities were added.
- [ ] **Gitmoji Commit Messages**:
  - Ensure all commits follow the [gitmoji](https://gitmoji.dev) convention (e.g. `:sparkles:` for Features, `:bug:` for Fixes, `:memo:` for Docs). GoReleaser uses gitmoji regexes to generate categorized release notes.

## 5. Working Tree & Release Tag

- [ ] **Clean Working Directory**:
  - Check `git status` to ensure no untracked or unstaged files remain.
  - Verify sensitive files (`client_secret.json`, `youtube.token.json`) are not tracked.
- [ ] **Sync with Remote Main**:
  ```shell
  git checkout main
  git pull origin main
  ```
- [ ] **Tag and Push**:
  The tag name **must** start with `v` (e.g., `v0.5.2`) to trigger the release workflow in `.github/workflows/publish.yml`:
  ```shell
  git tag -a vX.Y.Z -m ":bookmark: vX.Y.Z"
  git push origin main
  git push origin vX.Y.Z
  ```
