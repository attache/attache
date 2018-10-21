package view

import "io"

type View interface {
	Execute(out io.Writer, data interface{}) error
}

// None is a no-op View
var None View = noView{}

type noView struct{}

// Execute implements View for noView
func (noView) Execute(_ io.Writer, _ interface{}) error { return nil }
