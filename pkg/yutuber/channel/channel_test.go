package channel

import (
	"reflect"
	"testing"
)

func TestNewChannel(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want Channel
	}{
		{
			name: "TestNewChannel",
			args: args{
				opts: []Option{
					WithCategoryId("15"),
					WithForHandle("handle"),
					WithForUsername("username"),
					WithHl("hl"),
					WithID("id"),
					WithChannelManagedByMe(true, true),
					WithMaxResults(5),
					WithMine(false, true),
					WithMySubscribers(true, true),
					WithOnBehalfOfContentOwner("contentOwner"),
				},
			},
			want: &channel{
				CategoryId:             "15",
				ForHandle:              "handle",
				ForUsername:            "username",
				Hl:                     "hl",
				ID:                     "id",
				ManagedByMe:            &[]bool{true}[0],
				MaxResults:             5,
				Mine:                   &[]bool{false}[0],
				MySubscribers:          &[]bool{true}[0],
				OnBehalfOfContentOwner: "contentOwner",
			},
		},
		{
			name: "TestNewChannel",
			args: args{
				opts: []Option{
					WithCategoryId("20"),
					WithForHandle("handle"),
					WithCountry("US"),
					WithCustomUrl("customUrl"),
					WithDefaultLanguage("en"),
					WithDescription("description"),
					WithTitle("title"),
				},
			},
			want: &channel{
				CategoryId:      "20",
				ForHandle:       "handle",
				Country:         "US",
				CustomUrl:       "customUrl",
				DefaultLanguage: "en",
				Description:     "description",
				Title:           "title",
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewChannel(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewChannel() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
