// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"io"
)

type Getter[T any] interface {
	Get([]string) ([]*T, error)
}

type Lister interface {
	List([]string, string, string, io.Writer) error
}

type Inserter interface {
	Insert(string, string, io.Writer) error
}

type Deleter interface {
	Delete(writer io.Writer) error
}

type Updater interface {
	Update(string, string, io.Writer) error
}
