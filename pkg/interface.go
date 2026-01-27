// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"io"
)

type Getter[T any] interface {
	Get() ([]*T, error)
}

type Lister interface {
	List(io.Writer) error
}

type Inserter interface {
	Insert(io.Writer) error
}

type Deleter interface {
	Delete(io.Writer) error
}

type Updater interface {
	Update(io.Writer) error
}

type Downloader interface {
	Download(io.Writer) error
}

type RatingGetter interface {
	GetRating(io.Writer) error
}

type SpamMaker interface {
	MarkAsSpam(io.Writer) error
}

type Rater interface {
	Rate(io.Writer) error
}

type AbuseReporter interface {
	ReportAbuse(io.Writer) error
}

type Setter interface {
	Set(io.Writer) error
}

type Unsetter interface {
	Unset(io.Writer) error
}

type ModerationSetter interface {
	SetModerationStatus(io.Writer) error
}
