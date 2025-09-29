# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

Yutu is a CLI tool and MCP (Model Context Protocol) server for YouTube API operations, built with:

- **Command Structure**: Uses Cobra framework with commands organized in `cmd/` directory. Each YouTube resource (video, playlist, channel, etc.) has its own subdirectory with CRUD operations.

- **Main Entry Points**:
  - `main.go`: Imports all command modules and calls Execute()
  - `cmd/root.go`: Defines the root command and initializes viper for config
  - `cmd/mcp.go`: Implements the MCP server functionality

- **Authentication Flow**:
  - OAuth2 credentials stored in `client_secret.json`
  - Token cached in `youtube.token.json`
  - Environment variables: `YUTU_CREDENTIAL` and `YUTU_CACHE_TOKEN`

- **Resource Commands**: Each YouTube resource type (video, playlist, channel, comment, etc.) is implemented as a subcommand with standard operations (list, insert, update, delete).

- **MCP Server**: Can run in stdio or HTTP mode to serve as an MCP server for AI assistants, exposing YouTube operations as tools.

## Key Dependencies

- `spf13/cobra`: Command-line interface framework
- `google.golang.org/api`: Google APIs client library for YouTube
- `modelcontextprotocol/go-sdk`: MCP server implementation
- `golang.org/x/oauth2`: OAuth2 authentication

## Development Notes

- Follow gitmoji convention for commit messages
- The codebase supports both standard Go toolchain and Bazel build system
- GitHub Actions are configured for testing, code quality checks, and publishing