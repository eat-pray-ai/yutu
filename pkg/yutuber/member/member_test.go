package member

import (
	"reflect"
	"testing"
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
			name: "TestNewMember",
			args: args{
				opts: []Option{
					WithMemberChannelId("memberChannelId"),
					WithHasAccessToLevel("hasAccessToLevel"),
					WithMaxResults(5),
					WithMode("all_current"),
				},
			},
			want: &member{
				memberChannelId:  "memberChannelId",
				hasAccessToLevel: "hasAccessToLevel",
				maxResults:       5,
				mode:             "all_current",
			},
		},
		{
			name: "TestNewMember",
			args: args{
				opts: []Option{
					WithMemberChannelId("memberChannelId"),
					WithMaxResults(5),
					WithMode("updates"),
				},
			},
			want: &member{
				memberChannelId: "memberChannelId",
				maxResults:      5,
				mode:            "updates",
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
