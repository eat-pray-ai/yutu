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
		want Caption
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIDs([]string{"caption1", "caption2"}),
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
					WithService(&youtube.Service{}),
				},
			},
			want: &caption{
				IDs:                    []string{"caption1", "caption2"},
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
			want: &caption{},
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
			want: &caption{},
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
			want: &caption{
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
			want: &caption{
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
					WithIDs([]string{"caption1"}),
					WithLanguage("fr"),
					WithVideoId("video456"),
					WithIsCC(&isCCTrue),
				},
			},
			want: &caption{
				IDs:      []string{"caption1"},
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
					t.Errorf("NewCation() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
