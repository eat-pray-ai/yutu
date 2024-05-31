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
					WithId("id"),
					WithChannelManagedByMe("true"),
					WithMaxResults(5),
					WithMine("false"),
					WithMySubscribers("true"),
					WithOnBehalfOfContentOwner("contentOwner"),
				},
			},
			want: &channel{
				categoryId:             "15",
				forHandle:              "handle",
				forUsername:            "username",
				hl:                     "hl",
				id:                     "id",
				managedByMe:            "true",
				maxResults:             5,
				mine:                   "false",
				mySubscribers:          "true",
				onBehalfOfContentOwner: "contentOwner",
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
				categoryId:      "20",
				forHandle:       "handle",
				country:         "US",
				customUrl:       "customUrl",
				defaultLanguage: "en",
				description:     "description",
				title:           "title",
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
