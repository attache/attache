package main

import (
	"flag"
	"log"
	"strings"
)

type Context struct {
	Model *Model
	Views []View
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

	noViews := flag.Bool("noviews", false, "disables generation of views")
	name := flag.String("n", "", "name of the resource")
	defs := &fieldDefs{}
	flag.Var(defs, "f", "field definitions [-f NAME:TYPE:FLAG,FLAG...]+")
	flag.Parse()

	if *name == "" {
		panic("name cannot be empty")
	}

	if len(*defs) == 0 {
		panic("must specift at least one field")
	}

	c.Model = buildModel(*name, "", *defs)
	if !*noViews {
		c.Views = viewsFor(c.Model)
	}
}
