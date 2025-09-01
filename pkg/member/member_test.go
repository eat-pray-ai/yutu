package member

import (
	"reflect"
	"testing"

	"google.golang.org/api/youtube/v3"
)

func TestNewMember(t *testing.T) {
	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want Member
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithMemberChannelId("member123"),
					WithHasAccessToLevel("level1"),
					WithMaxResults(100),
					WithMode("all_current"),
					WithService(&youtube.Service{}),
				},
			},
			want: &member{
				MemberChannelId:  "member123",
				HasAccessToLevel: "level1",
				MaxResults:       100,
				Mode:             "all_current",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &member{},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &member{
				MaxResults: 1,
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-15),
				},
			},
			want: &member{
				MaxResults: 1,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithMemberChannelId(""),
					WithHasAccessToLevel(""),
					WithMode(""),
				},
			},
			want: &member{
				MemberChannelId:  "",
				HasAccessToLevel: "",
				Mode:             "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithMemberChannelId("channel456"),
					WithMaxResults(50),
				},
			},
			want: &member{
				MemberChannelId: "channel456",
				MaxResults:      50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewMember(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewMember() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
