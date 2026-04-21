# internal/

Internal tools and private packages. Not importable by external code.

## Tools

| Tool | Description |
|------|-------------|
| `tools/skillgen/` | Skill file generator |
| `tools/cmdtestgen/` | Command test script generator |

Each tool is a standalone `main.go` with its own `BUILD.bazel` (auto-generated — do NOT edit manually).
