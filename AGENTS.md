# Agent Guidelines

## OVERVIEW

Go CLI + MCP server for YouTube.

## STRUCTURE

```text
.
├── assets/     # Static assets (logos, etc.)
├── cmd/        # CLI command definitions (Cobra)
├── dist/       # Build artifacts
├── internal/   # Internal tools and private packages
├── pkg/        # Core domain logic and infrastructure
└── scripts/    # Utility scripts and smoke tests
```

## WHERE TO LOOK

- **CLI Entry**: `main.go` -> `cmd/`
- **Domain Logic**: `pkg/<resource>/` (e.g., `pkg/video/`)
- **Infrastructure**: `pkg/auth`, `pkg/utils`
- **Tests**: Co-located `*_test.go` files within `pkg/`

## CODE MAP

- `main.go`: Application entry point.
- `cmd/root.go`: Cobra root command and global flag definitions.
- `pkg/video/video.go`: Example of domain logic implementation.

## CONVENTIONS

- **Build System**: Bazel is used for build/test, though standard Go tools (`go build`, `go test`) also work.
- **Testing**: Table-driven tests are preferred for consistency and coverage.
- **Secrets**: `client_secret.json` and `youtube.token.json` are typically stored in the root directory (standard for
  this project).

## COMMANDS

- **Build**: `go build -o yutu .` or `bazel build //...`
- **Test**: `go test ./...` or `bazel test //...`
- **Smoke Tests**: `./scripts/command-test.sh`
