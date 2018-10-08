package cmd_new

import "github.com/attache/attache/cmd/attache/shared"

//go:generate go-bindata -pkg cmd_new templates

var Export = shared.NewPlugin(
	"new",
	func() shared.Command { return &Context{} },
)
