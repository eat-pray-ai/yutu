// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package resources_test

import (
	"testing"

	"github.com/eat-pray-ai/yutu/cmd"
	_ "github.com/eat-pray-ai/yutu/cmd/agent"
	_ "github.com/eat-pray-ai/yutu/cmd/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestFlagConsistency(t *testing.T) {
	canonical := make(map[string]string)
	firstSeen := make(map[string]string)

	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		cmdPath := c.CommandPath()
		c.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Name == "help" || f.Hidden {
				return
			}
			if f.Shorthand == "y" && f.Name != "yes" {
				t.Errorf(
					"command %q: flag --%s uses disallowed shorthand -y",
					cmdPath, f.Name,
				)
			}
			if f.Shorthand != "" {
				if prev, ok := canonical[f.Name]; !ok {
					canonical[f.Name] = f.Shorthand
					firstSeen[f.Name] = cmdPath
				} else if prev != f.Shorthand {
					t.Errorf(
						"flag --%s shorthand drift: %q on %q vs %q on %q",
						f.Name, prev, firstSeen[f.Name], f.Shorthand, cmdPath,
					)
				}
			}
		})
		for _, sub := range c.Commands() {
			walk(sub)
		}
	}

	walk(cmd.RootCmd)
}
