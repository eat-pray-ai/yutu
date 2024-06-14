package caption

import (
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
					WithId("id"),
					WithFile("file"),
					WithAudioTrackType("audioTrackType"),
					WithIsAutoSynced(true, true),
					WithIsCC(true, true),
					WithIsDraft(true, true),
					WithIsEasyReader(true, true),
					WithIsLarge(true, true),
					WithLanguage("language"),
					WithName("name"),
				},
			},
			want: &caption{
				id:             "id",
				file:           "file",
				audioTrackType: "audioTrackType",
				isAutoSynced:   &[]bool{true}[0],
				isCC:           &[]bool{true}[0],
				isDraft:        &[]bool{true}[0],
				isEasyReader:   &[]bool{true}[0],
				isLarge:        &[]bool{true}[0],
				language:       "language",
				name:           "name",
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
