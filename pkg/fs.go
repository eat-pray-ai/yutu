package pkg

import (
	"io/fs"
	"os"
)

var Fsys fs.FS

func init() {
	root, err := os.OpenRoot("/")
	if err != nil {
		panic(err)
	}
	Fsys = root.FS()
}
