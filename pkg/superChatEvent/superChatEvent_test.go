package superChatEvent

import (
	"reflect"
	"testing"
)

func TestNewSuperChatEvent(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want SuperChatEvent
	}{
		{
			name: "TestNewSuperChatEvent",
			args: args{
				opts: []Option{
					WithHl("hl"),
					WithMaxResults(1),
				},
			},
			want: &superChatEvent{
				Hl:         "hl",
				MaxResults: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewSuperChatEvent(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewSuperChatEvent() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
