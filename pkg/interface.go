package pkg

import (
	"io"
)

type Lister interface {
	List([]string, string, io.Writer) error
}

type Inserter interface {
	Insert(string, io.Writer) error
}

type Deleter interface {
	Delete(writer io.Writer) error
}

type Updater interface {
	Update(string, io.Writer) error
}

type Setter interface {
	Set(writer io.Writer) error
}

type Unsetter interface {
	Unset(writer io.Writer) error
}
