# internal/

Internal tools and private packages. Not importable by external code.

## Tools

| Tool | Description |
|------|-------------|
| `tools/docgen/` | Documentation generator |
| `tools/skillgen/` | Skill file generator |

Each tool is a standalone `main.go` with its own `BUILD.bazel` (auto-generated — do NOT edit manually).
