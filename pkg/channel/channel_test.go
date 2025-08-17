package channel

import (
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/utils"
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
					WithIDs([]string{"id1", "id2"}),
					WithChannelManagedByMe(utils.BoolPtr("true")),
					WithMaxResults(5),
					WithMine(utils.BoolPtr("false")),
					WithMySubscribers(utils.BoolPtr("true")),
					WithOnBehalfOfContentOwner("contentOwner"),
				},
			},
			want: &channel{
				CategoryId:             "15",
				ForHandle:              "handle",
				ForUsername:            "username",
				Hl:                     "hl",
				IDs:                    []string{"id1", "id2"},
				ManagedByMe:            utils.BoolPtr("true"),
				MaxResults:             5,
				Mine:                   utils.BoolPtr("false"),
				MySubscribers:          utils.BoolPtr("true"),
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
