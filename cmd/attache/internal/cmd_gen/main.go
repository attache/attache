package cmd_gen

import "github.com/mccolljr/attache/cmd/attache/shared"

//go:generate go-bindata -pkg cmd_gen templates

var Export = shared.NewPlugin(
	"gen",
	func() shared.Command { return &Context{} },
)
