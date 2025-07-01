package pkg

import (
	"io"
)

type Lister interface {
	List([]string, string, string, io.Writer) error
}

type Inserter interface {
	Insert(string, string, io.Writer) error
}

type Deleter interface {
	Delete(writer io.Writer) error
}

type Updater interface {
	Update(string, string, io.Writer) error
}

type Setter interface {
	Set(string, io.Writer) error
}

type Unsetter interface {
	Unset(io.Writer) error
}
