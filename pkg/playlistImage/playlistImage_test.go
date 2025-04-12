package playlistImage

import (
	"reflect"
	"testing"
)

func TestNewPlaylistImage(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want PlaylistImage
	}{
		{
			name: "TestNewPlaylistImage",
			args: args{
				opts: []Option{
					WithID("id"),
					WithKind("kind"),
					WithHeight(100),
					WithPlaylistID("playlistID"),
					WithType("type"),
					WithWidth(200),
					WithFile("file"),
					WithParent("parent"),
					WithMaxResults(10),
					WithOnBehalfOfContentOwner("onBehalfOfContentOwner"),
					WithOnBehalfOfContentOwnerChannel("onBehalfOfContentOwnerChannel"),
				},
			},
			want: &playlistImage{
				ID:                            "id",
				Kind:                          "kind",
				Height:                        100,
				PlaylistID:                    "playlistID",
				Type:                          "type",
				Width:                         200,
				File:                          "file",
				Parent:                        "parent",
				MaxResults:                    10,
				OnBehalfOfContentOwner:        "onBehalfOfContentOwner",
				OnBehalfOfContentOwnerChannel: "onBehalfOfContentOwnerChannel",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPlaylistImage(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPlaylistImage() = %v, want %v", got, tt.want)
			}
		})
	}
}
