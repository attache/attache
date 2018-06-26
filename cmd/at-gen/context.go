package main

import (
	"flag"
	"log"
	"strings"
)

type Context struct {
	Model *Model
	Views []View

	DoViews  bool
	DoRoutes bool
	DoModel  bool

	Replace bool
}

type fieldDefs []string

func (f fieldDefs) String() string      { return strings.Join(f, " ") }
func (f *fieldDefs) Set(s string) error { *f = append(*f, s); return nil }

func (c *Context) Init() {
	defer func() {
		val := recover()
		if val != nil {
			log.Fatalln(val)
		}
	}()

	flag.BoolVar(&c.DoModel, "model", false, "generate model")
	flag.BoolVar(&c.DoViews, "views", false, "generate views")
	flag.BoolVar(&c.DoRoutes, "routes", false, "generate routes")
	flag.BoolVar(&c.Replace, "replace", false, "replace existing files")
	name := flag.String("n", "", "name of the resource")
	table := flag.String("t", "", "name of table")
	defs := &fieldDefs{}
	flag.Var(defs, "f", "-f NAME:TYPE:FLAGs [...]")

	flag.Parse()

	if *name == "" {
		panic("name cannot be empty")
	}

	if len(*defs) == 0 {
		panic("must specify at least one field")
	}

	// if none specified, do all
	if !c.DoModel && !c.DoViews && !c.DoRoutes {
		c.DoModel = true
		c.DoViews = true
		c.DoRoutes = true
	}

	// needed for more than just model
	c.Model = buildModel(*name, *table, *defs)

	if c.DoViews {
		c.Views = viewsFor(c.Model)
	}

}
