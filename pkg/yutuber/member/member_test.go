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
				MemberChannelId:  "memberChannelId",
				HasAccessToLevel: "hasAccessToLevel",
				MaxResults:       5,
				Mode:             "all_current",
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
				MemberChannelId: "memberChannelId",
				MaxResults:      5,
				Mode:            "updates",
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
