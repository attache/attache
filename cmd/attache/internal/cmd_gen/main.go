package cmd_gen

import "github.com/attache/attache/cmd/attache/shared"

//go:generate go run github.com/kevinburke/go-bindata/go-bindata -pkg cmd_gen templates

var Export = shared.NewPlugin(
	"gen",
	func() shared.Command { return &Context{} },
)
