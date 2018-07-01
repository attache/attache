package main

import (
	"github.com/mccolljr/attache/cmd/attache/internal/cmd_gen"
	"github.com/mccolljr/attache/cmd/attache/internal/cmd_new"
	"github.com/mccolljr/attache/cmd/attache/internal/shared"
)

var builtins = map[string]shared.Plugin{
	cmd_new.Export.Name(): cmd_new.Export,
	cmd_gen.Export.Name(): cmd_gen.Export,
}
