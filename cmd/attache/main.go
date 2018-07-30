package main

import (
	"log"
	"os"
	"path/filepath"
	"plugin"

	"github.com/mccolljr/attache/cmd/attache/shared"
)

const helpText = `attache: cli for the attache framework

COMMANDS:
	new: create a new project
	gen: generate models, views, and/or routes
`

func init() {
	log.SetFlags(0)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln(helpText)
	}

	cmd, cmdArgs := os.Args[1], os.Args[2:]

	if plug := builtins[cmd]; plug != nil {
		plugCmd := plug.Command()
		if err := plugCmd.Execute(cmdArgs); err != nil {
			log.Fatalf("%s: %s", plug.Name(), err)
		}

		return // end
	}

	plugDir := filepath.Join(os.Getenv("HOME"), ".attache", "plugins")
	plugFile := filepath.Join(plugDir, cmd+".so")

	loaded, err := plugin.Open(plugFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("cannot find plugin %s", plugFile)
		}

		log.Fatalf("error loading plugin %s", plugFile)
	}

	sym, _ := loaded.Lookup("Export")

	if plug, ok := sym.(shared.Plugin); ok {
		plugCmd := plug.Command()
		if err := plugCmd.Execute(cmdArgs); err != nil {
			log.Fatalf("%s: %s", plug.Name(), err)
		}

		return // end
	}

	log.Fatalf("invalid plugin %s", plugFile)
}
