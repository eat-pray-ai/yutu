# cmd/

CLI command definitions and MCP tool bindings.

## Wiring

- `main.go` → `cmd.Execute()` → `root.go`
- Each `<resource>/` dir registers a subcommand + MCP tool
- Each subcommand defines `<subcommand>Opts` with `bindFlags(c *cobra.Command)` and `build(...)`, constructed by `new<Sub>Cmd()`

## Conventions

- No shared package-level flag variables; bind all flags (including `--yes` and `--output`) to per-command `<subcommand>Opts` structs.
- Use `utils.OptIfChanged` inside `build(f *pflag.FlagSet)` for optional pointer/slice flags so omitted flags stay `nil` (and typos fail fast).
- Use `utils.ConfirmPreRun(c, opts.yes, msg)` in `PreRunE` for mutating/destructive commands.
- MCP tools registered via `mcp.AddTool`.
- `BUILD.bazel` files are auto-generated — do NOT create or edit them manually.

## Subcommands

Each subdirectory corresponds to a YouTube API resource (e.g. `video/`, `channel/`, `playlist/`) or tool subcommand (`agent/`). All follow the same pattern — see [Wiring](#wiring) above.
