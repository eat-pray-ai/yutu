// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/eat-pray-ai/yutu/cmd"
	_ "github.com/eat-pray-ai/yutu/cmd/resources"
)

//go:generate go tool go-winres make --arch amd64,arm64 --product-version git-tag --file-version git-tag
func main() {
	cmd.Execute()
}
