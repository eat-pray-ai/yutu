# OVERVIEW

CLI command definitions and MCP tool bindings.

# STRUCTURE

- `root.go`: Root command entry point.
- `<resource>/`: Subcommands (e.g., `video/`, `channel/`).

# WIRING PATTERN

- `main.go` -> `cmd.Execute()` -> `root.go`.
- `<resource>.go`: Registers subcommand + MCP tool.
- `Run`: Binds flags -> Calls `pkg/<resource>` method.

# CONVENTIONS

- `resetFlags` in `PersistentPreRun`: Ensures clean state for MCP.
- MCP tools registered via `mcp.AddTool`.
- Flags bound to package-level variables.
