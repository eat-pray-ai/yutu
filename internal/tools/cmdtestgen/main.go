// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"

	_ "github.com/eat-pray-ai/yutu/cmd/activity"
	_ "github.com/eat-pray-ai/yutu/cmd/agent"
	_ "github.com/eat-pray-ai/yutu/cmd/caption"
	_ "github.com/eat-pray-ai/yutu/cmd/channel"
	_ "github.com/eat-pray-ai/yutu/cmd/channelBanner"
	_ "github.com/eat-pray-ai/yutu/cmd/channelSection"
	_ "github.com/eat-pray-ai/yutu/cmd/comment"
	_ "github.com/eat-pray-ai/yutu/cmd/commentThread"
	_ "github.com/eat-pray-ai/yutu/cmd/i18nLanguage"
	_ "github.com/eat-pray-ai/yutu/cmd/i18nRegion"
	_ "github.com/eat-pray-ai/yutu/cmd/member"
	_ "github.com/eat-pray-ai/yutu/cmd/membershipsLevel"
	_ "github.com/eat-pray-ai/yutu/cmd/playlist"
	_ "github.com/eat-pray-ai/yutu/cmd/playlistImage"
	_ "github.com/eat-pray-ai/yutu/cmd/playlistItem"
	_ "github.com/eat-pray-ai/yutu/cmd/search"
	_ "github.com/eat-pray-ai/yutu/cmd/subscription"
	_ "github.com/eat-pray-ai/yutu/cmd/superChatEvent"
	_ "github.com/eat-pray-ai/yutu/cmd/thumbnail"
	_ "github.com/eat-pray-ai/yutu/cmd/video"
	_ "github.com/eat-pray-ai/yutu/cmd/videoAbuseReportReason"
	_ "github.com/eat-pray-ai/yutu/cmd/videoCategory"
	_ "github.com/eat-pray-ai/yutu/cmd/watermark"
)

const header = `#!/usr/bin/env bash
# Copyright 2025 eat-pray-ai & OpenWaygate
# SPDX-License-Identifier: Apache-2.0

set -euo pipefail

# Accept yutu executable path as first parameter
YUTU_PATH="${1:-}"

# If no path provided, build yutu
if [[ -z "$YUTU_PATH" ]]; then
    echo "No yutu path provided, building yutu..."

    MOD="github.com/eat-pray-ai/yutu/cmd"
    Version="${MOD}.Version=$(git describe --tags --always --dirty)"
    Commit="${MOD}.Commit=$(git rev-parse --short HEAD)"
    CommitDate="${MOD}.CommitDate=$(git log -1 --date='format:%Y-%m-%dT%H:%M:%SZ' --pretty=%cd)"
    Os="${MOD}.Os=$(go env GOOS)"
    Arch="${MOD}.Arch=$(go env GOARCH)"
    Builder="${MOD}.Builder=${GITHUB_ACTOR:-$USER}"
    ldflags="-s -X ${Version} -X ${Commit} -X ${CommitDate} -X ${Os} -X ${Arch} -X ${Builder}"

    go mod download
    go build -ldflags "${ldflags}" -o yutu .
    trap 'rm yutu' EXIT
    chmod +x yutu

    YUTU_PATH="./yutu"
else
    # Verify the provided path exists and is executable
    if [[ ! -f "$YUTU_PATH" ]]; then
        echo "Error: yutu executable not found at: $YUTU_PATH"
        exit 1
    fi
    if [[ ! -x "$YUTU_PATH" ]]; then
        echo "Error: $YUTU_PATH is not executable"
        exit 1
    fi
    echo "Using yutu at: $YUTU_PATH"
fi

# Run tests with the yutu executable
"$YUTU_PATH" help
"$YUTU_PATH" completion
"$YUTU_PATH" version
`

func isResource(c *cobra.Command) bool {
	return strings.HasPrefix(c.Short, "Manage")
}

func generate() string {
	root := cmd.RootCmd
	root.InitDefaultHelpCmd()

	var b strings.Builder
	b.WriteString(header)

	var coreCommands, resourceCommands []*cobra.Command
	for _, c := range root.Commands() {
		name := c.Name()
		if name == "help" || name == "completion" || name == "version" {
			continue
		}
		if isResource(c) {
			resourceCommands = append(resourceCommands, c)
		} else {
			coreCommands = append(coreCommands, c)
		}
	}

	sort.Slice(
		coreCommands, func(i, j int) bool {
			return coreCommands[i].Name() < coreCommands[j].Name()
		},
	)
	sort.Slice(
		resourceCommands, func(i, j int) bool {
			return resourceCommands[i].Name() < resourceCommands[j].Name()
		},
	)

	b.WriteString("\n# yutu\n")
	for _, c := range coreCommands {
		name := c.Name()
		fmt.Fprintf(&b, "echo \"======= %s =======\"\n", name)
		fmt.Fprintf(&b, "\"$YUTU_PATH\" %s --help\n\n", name)
	}

	b.WriteString("# youtube api\n")
	for _, c := range resourceCommands {
		name := c.Name()
		fmt.Fprintf(&b, "echo \"======= %s =======\"\n", name)
		fmt.Fprintf(&b, "\"$YUTU_PATH\" %s --help\n", name)

		var subs []*cobra.Command
		for _, sub := range c.Commands() {
			if sub.Name() == "help" {
				continue
			}
			subs = append(subs, sub)
		}
		sort.Slice(
			subs, func(i, j int) bool {
				return subs[i].Name() < subs[j].Name()
			},
		)

		for _, sub := range subs {
			fmt.Fprintf(&b, "echo \"------- %s -------\"\n", sub.Name())
			fmt.Fprintf(&b, "\"$YUTU_PATH\" %s %s --help\n", name, sub.Name())
		}
		b.WriteString("\n")
	}

	return strings.TrimRight(b.String(), "\n") + "\n"
}

func main() {
	out := flag.String("out", "", "output file (default: stdout)")
	flag.Parse()

	content := generate()

	if *out == "" {
		fmt.Print(content)
		return
	}

	if err := os.WriteFile(*out, []byte(content), 0o755); err != nil {
		log.Fatalf("write %s: %v", *out, err)
	}
	fmt.Printf("Generated %s\n", *out)
}
