# cmd/

CLI command definitions and MCP tool bindings.

## Wiring

- `main.go` → `cmd.Execute()` → `root.go`
- Each `<resource>/` dir registers a subcommand + MCP tool
- `Run` function: binds flags → calls `pkg/<resource>` method

## Conventions

- `resetFlags` in `PersistentPreRun`: ensures clean state for MCP.
- MCP tools registered via `mcp.AddTool`.
- Flags bound to package-level variables.
- `BUILD.bazel` files are auto-generated — do NOT create or edit them manually.

## Subcommands

Each subdirectory is a 1:1 mapping to a YouTube API resource:

`activity/`, `agent/`, `caption/`, `channel/`, `channelBanner/`, `channelSection/`, `comment/`, `commentThread/`, `i18nLanguage/`, `i18nRegion/`, `member/`, `membershipsLevel/`, `playlist/`, `playlistImage/`, `playlistItem/`, `search/`, `subscription/`, `superChatEvent/`, `thumbnail/`, `video/`, `videoAbuseReportReason/`, `videoCategory/`, `watermark/`

All follow the same pattern — see [Wiring](#wiring) above.
