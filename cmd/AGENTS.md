# cmd/

CLI command definitions and MCP tool bindings.

## Wiring

- `main.go` → `cmd.Execute()` → `root.go`
- Each `<resource>/` dir registers a subcommand + MCP tool
- `Run` function: binds flags → calls `pkg/<resource>` method

## Conventions

- `utils.ResetFlags` in `PersistentPreRun`: resets omitted optional pointer flags (`!flag.Changed`) to `nil`.
- MCP tools registered via `mcp.AddTool`.
- Flags bound to package-level variables.
- `BUILD.bazel` files are auto-generated — do NOT create or edit them manually.

## Subcommands

Each subdirectory corresponds to a YouTube API resource (e.g. `video/`, `channel/`, `playlist/`) or tool subcommand (`agent/`). All follow the same pattern — see [Wiring](#wiring) above.
