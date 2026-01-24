// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package caption

import (
	"reflect"
	"testing"

	"google.golang.org/api/youtube/v3"
)

func TestNewCation(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}
	isAutoSyncedTrue := true
	isAutoSyncedFalse := false
	isCCTrue := true
	isCCFalse := false
	isDraftTrue := true
	isDraftFalse := false
	isEasyReaderTrue := true
	isEasyReaderFalse := false
	isLargeTrue := true
	isLargeFalse := false

	tests := []struct {
		name string
		args args
		want ICaption[youtube.Caption]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIds([]string{"caption1", "caption2"}),
					WithFile("/path/to/file.srt"),
					WithAudioTrackType("primary"),
					WithIsAutoSynced(&isAutoSyncedTrue),
					WithIsCC(&isCCTrue),
					WithIsDraft(&isDraftTrue),
					WithIsEasyReader(&isEasyReaderTrue),
					WithIsLarge(&isLargeTrue),
					WithLanguage("en"),
					WithName("English Captions"),
					WithTrackKind("standard"),
					WithOnBehalfOf("channel123"),
					WithOnBehalfOfContentOwner("owner123"),
					WithVideoId("video123"),
					WithTfmt("srt"),
					WithTlang("es"),
					WithService(svc),
				},
			},
			want: &Caption{
				service:                svc,
				Ids:                    []string{"caption1", "caption2"},
				File:                   "/path/to/file.srt",
				AudioTrackType:         "primary",
				IsAutoSynced:           &isAutoSyncedTrue,
				IsCC:                   &isCCTrue,
				IsDraft:                &isDraftTrue,
				IsEasyReader:           &isEasyReaderTrue,
				IsLarge:                &isLargeTrue,
				Language:               "en",
				Name:                   "English Captions",
				TrackKind:              "standard",
				OnBehalfOf:             "channel123",
				OnBehalfOfContentOwner: "owner123",
				VideoId:                "video123",
				Tfmt:                   "srt",
				Tlang:                  "es",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &Caption{},
		},
		{
			name: "with nil boolean options",
			args: args{
				opts: []Option{
					WithIsAutoSynced(nil),
					WithIsCC(nil),
					WithIsDraft(nil),
					WithIsEasyReader(nil),
					WithIsLarge(nil),
				},
			},
			want: &Caption{},
		},
		{
			name: "with false boolean options",
			args: args{
				opts: []Option{
					WithIsAutoSynced(&isAutoSyncedFalse),
					WithIsCC(&isCCFalse),
					WithIsDraft(&isDraftFalse),
					WithIsEasyReader(&isEasyReaderFalse),
					WithIsLarge(&isLargeFalse),
				},
			},
			want: &Caption{
				IsAutoSynced: &isAutoSyncedFalse,
				IsCC:         &isCCFalse,
				IsDraft:      &isDraftFalse,
				IsEasyReader: &isEasyReaderFalse,
				IsLarge:      &isLargeFalse,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithFile(""),
					WithAudioTrackType(""),
					WithLanguage(""),
					WithName(""),
					WithTrackKind(""),
					WithOnBehalfOf(""),
					WithOnBehalfOfContentOwner(""),
					WithVideoId(""),
					WithTfmt(""),
					WithTlang(""),
				},
			},
			want: &Caption{
				File:                   "",
				AudioTrackType:         "",
				Language:               "",
				Name:                   "",
				TrackKind:              "",
				OnBehalfOf:             "",
				OnBehalfOfContentOwner: "",
				VideoId:                "",
				Tfmt:                   "",
				Tlang:                  "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithIds([]string{"caption1"}),
					WithLanguage("fr"),
					WithVideoId("video456"),
					WithIsCC(&isCCTrue),
				},
			},
			want: &Caption{
				Ids:      []string{"caption1"},
				Language: "fr",
				VideoId:  "video456",
				IsCC:     &isCCTrue,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewCation(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("%s\nNewCation() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}
