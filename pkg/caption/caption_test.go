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
					WithID("ID"),
					WithFile("File"),
					WithAudioTrackType("AudioTrackType"),
					WithIsAutoSynced(true, true),
					WithIsCC(true, true),
					WithIsDraft(true, true),
					WithIsEasyReader(true, true),
					WithIsLarge(true, true),
					WithLanguage("Language"),
					WithName("Name"),
				},
			},
			want: &caption{
				ID:             "ID",
				File:           "File",
				AudioTrackType: "AudioTrackType",
				IsAutoSynced:   &[]bool{true}[0],
				IsCC:           &[]bool{true}[0],
				IsDraft:        &[]bool{true}[0],
				IsEasyReader:   &[]bool{true}[0],
				IsLarge:        &[]bool{true}[0],
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
