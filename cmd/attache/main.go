package main

import (
	"log"
	"os"

	"github.com/alexflint/go-arg"
)

func init() { log.SetFlags(0) }

// CLI described the available subcommands for the attache CLI.
type CLI struct {
	New *CommandNew `arg:"subcommand:new" help:"generate a new attache project"`
	Gen *CommandGen `arg:"subcommand:gen" help:"generate files within an attache project"`
}

func main() {
	var args CLI
	p := arg.MustParse(&args)
	switch {
	case args.New != nil:
		args.New.Version = Version
		if err := args.New.Execute(); err != nil {
			log.Fatalln(err)
		}
		return
	case args.Gen != nil:
		if err := args.Gen.Execute(); err != nil {
			log.Fatalln(err)
		}
		return
	default:
		p.WriteHelp(os.Stderr)
		os.Exit(1)
	}
}
