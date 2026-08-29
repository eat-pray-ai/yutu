// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"cmp"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

type verbEntry struct {
	name  string
	short string
}

type resourceEntry struct {
	name  string
	kebab string
	human string
	short string
	long  string
	verbs []verbEntry
}

func camelToWords(s string) string {
	return strings.Join(camelSplit(s), " ")
}

func camelToKebab(s string) string {
	return strings.Join(camelSplit(s), "-")
}

func camelSplit(s string) []string {
	var words []string
	var cur []byte
	for i := range len(s) {
		ch := s[i]
		if ch >= 'A' && ch <= 'Z' && len(cur) > 0 {
			words = append(words, strings.ToLower(string(cur)))
			cur = cur[:0]
		}
		cur = append(cur, ch)
	}
	if len(cur) > 0 {
		words = append(words, strings.ToLower(string(cur)))
	}
	return words
}

func collectResources(root *cobra.Command) []resourceEntry {
	var resources []resourceEntry
	for _, c := range root.Commands() {
		if !strings.HasPrefix(c.Short, "Manage") {
			continue
		}
		name := c.Name()
		var verbs []verbEntry
		for _, sub := range c.Commands() {
			if sub.Name() == "help" {
				continue
			}
			verbs = append(verbs, verbEntry{name: sub.Name(), short: sub.Short})
		}
		slices.SortFunc(verbs, func(a, b verbEntry) int {
			return cmp.Compare(a.name, b.name)
		})
		resources = append(resources, resourceEntry{
			name:  name,
			kebab: camelToKebab(name),
			human: camelToWords(name),
			short: c.Short,
			long:  c.Long,
			verbs: verbs,
		})
	}
	slices.SortFunc(resources, func(a, b resourceEntry) int {
		return cmp.Compare(a.name, b.name)
	})
	return resources
}