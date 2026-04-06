# Yutu

Go CLI + MCP server + Agent for YouTube.

## Quick Reference

- **Build**: `go build -o yutu .` or `bazel build //...`
- **Test**: `go test ./...` or `bazel test //...`
- **Smoke Tests**: `./scripts/command-test.sh`
- **Entry Point**: `main.go`

## Directory Index

| Directory | Description |
|-----------|-------------|
| [cmd/](cmd/AGENTS.md) | CLI command definitions and MCP tool bindings |
| [pkg/](pkg/AGENTS.md) | Core domain logic and shared infrastructure |
| [internal/](internal/AGENTS.md) | Internal tools (docgen, skillgen) |
| [scripts/](scripts/AGENTS.md) | Utility scripts and smoke tests |
| [docs/](docs/) | Project documentation |

## Documentation

- [docs/FEATURES.md](docs/FEATURES.md) — Feature overview
- [docs/HOW_TO_TEST.md](docs/HOW_TO_TEST.md) — Testing guide
- [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) — Contribution guidelines
- [docs/CODE_OF_CONDUCT.md](docs/CODE_OF_CONDUCT.md) — Code of conduct

## Conventions

- **Secrets**: `client_secret.json` and `youtube.token.json` in root (standard for this project).
- **Build System**: Bazel is primary, standard Go tools also work.
- **BUILD.bazel files are auto-generated** — do NOT create or edit them manually. Run `bazel run //:gazelle` to regenerate.
- After changing dependencies: `bazel run @rules_go//go -- mod tidy -v && bazel mod tidy`.
- See [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) for the full list of useful build/test/release commands.
