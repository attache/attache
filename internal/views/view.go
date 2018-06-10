package views

import (
	"io"
)

type View interface {
	Execute(out io.Writer, data interface{}) error
}
