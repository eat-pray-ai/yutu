package main

import (
	"github.com/eat-pray-ai/yutu/cmd"
	_ "github.com/eat-pray-ai/yutu/cmd/channel"
	_ "github.com/eat-pray-ai/yutu/cmd/playlist"
	_ "github.com/eat-pray-ai/yutu/cmd/video"
)

func main() {
	cmd.Execute()
}
