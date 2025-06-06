package caption

import (
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"reflect"
	"testing"
)

func TestNewCation(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want Caption
	}{
		{
			name: "TestNewCaption",
			args: args{
				opts: []Option{
					WithIDs([]string{"id1", "id2"}),
					WithFile("File"),
					WithAudioTrackType("AudioTrackType"),
					WithIsAutoSynced(utils.BoolPtr("true")),
					WithIsCC(utils.BoolPtr("true")),
					WithIsDraft(utils.BoolPtr("true")),
					WithIsEasyReader(utils.BoolPtr("true")),
					WithIsLarge(utils.BoolPtr("true")),
					WithLanguage("Language"),
					WithName("Name"),
				},
			},
			want: &caption{
				IDs:            []string{"id1", "id2"},
				File:           "File",
				AudioTrackType: "AudioTrackType",
				IsAutoSynced:   utils.BoolPtr("true"),
				IsCC:           utils.BoolPtr("true"),
				IsDraft:        utils.BoolPtr("true"),
				IsEasyReader:   utils.BoolPtr("true"),
				IsLarge:        utils.BoolPtr("true"),
				Language:       "Language",
				Name:           "Name",
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
