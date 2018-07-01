package cmd_gen

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mccolljr/attache/cmd/attache/shared"
	"golang.org/x/tools/imports"
)

type Context struct {
	flags *flag.FlagSet

	ContextType string

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

func (c *Context) recover(err *error) {
	val := recover()
	if val != nil {
		*err = fmt.Errorf("%s", val)
	}
}

func (c *Context) Execute(args []string) (err error) {
	defer c.recover(&err)

	if err = c.init(args); err != nil {
		return err
	}

	if err = c.do(); err != nil {
		return err
	}

	return nil
}

func (c *Context) init(args []string) error {
	c.flags = flag.NewFlagSet("gen", flag.ContinueOnError)

	conf, err := shared.GetConfig()
	if err != nil {
		return err
	}

	c.flags.StringVar(&c.ContextType, "ctx", conf.GetString("ContextType"), "name of app's Context type")
	c.flags.BoolVar(&c.DoModel, "model", false, "generate model")
	c.flags.BoolVar(&c.DoViews, "views", false, "generate views")
	c.flags.BoolVar(&c.DoRoutes, "routes", false, "generate routes")
	c.flags.BoolVar(&c.Replace, "replace", false, "replace existing files")
	name := c.flags.String("n", "", "name of the resource")
	table := c.flags.String("t", "", "name of table")
	defs := &fieldDefs{}
	c.flags.Var(defs, "f", "-f NAME:TYPE:FLAGs [...]")

	if err := c.flags.Parse(args); err != nil {
		return err
	}

	if c.ContextType == "" {
		return errors.New("ctx cannot be \"\"")
	}

	if *name == "" {
		return errors.New("name cannot be \"\"")
	}

	if len(*defs) == 0 {
		return errors.New("must specify at least one field")
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

	return nil
}

func (ctx *Context) do() error {
	// prevent overwrite of existing model file
	modelFile := filepath.Join("models", ctx.Model.Table+".go")
	if ctx.DoModel {
		// make sure file doesn;t already exist
		if _, err := os.Stat(modelFile); err == nil {
			if !ctx.Replace {
				return fmt.Errorf("%s already exists", modelFile)
			}
		} else if !os.IsNotExist(err) {
			return err
		}

		// ensure containing directory exists
		if err := os.MkdirAll("models", os.ModePerm); err != nil {
			return err
		}
	}

	if ctx.DoViews {
		// prevent overwrite of view files
		for _, v := range ctx.Views {
			if _, err := os.Stat(v.File); err == nil {
				if !ctx.Replace {
					return fmt.Errorf("%s already exists", v.File)
				}
			} else if !os.IsNotExist(err) {
				return err
			}
		}
	}

	routeFile := ctx.Model.Table + "_routes.go"
	if ctx.DoRoutes {
		// prevent overwrite of route file
		if _, err := os.Stat(routeFile); err == nil {
			if !ctx.Replace {
				return fmt.Errorf("%s already exists", routeFile)
			}
		} else if !os.IsNotExist(err) {
			return err
		}
	}

	var buf bytes.Buffer

	if ctx.DoModel {
		buf.Reset()

		// generate model file
		tpl, err := template.New("").Parse(MustAssetString("templates/model.tpl"))
		if err != nil {
			return err
		}

		if err = tpl.Execute(&buf, ctx.Model); err != nil {
			return err
		}

		formattedModel, err := imports.Process(modelFile, buf.Bytes(), nil)
		if err != nil {
			log.Println("goimports failed:\n\n\n", buf.String(), "\n\n")
			return err
		}

		file, err := os.Create(modelFile)
		if err != nil {
			return err
		}

		if _, err = file.Write(formattedModel); err != nil {
			return err
		}
	}

	// generate route file
	if ctx.DoRoutes {
		buf.Reset()

		tpl, err := template.New("").Parse(MustAssetString("templates/routes.tpl"))
		if err != nil {
			return err
		}

		if err = tpl.Execute(&buf, ctx); err != nil {
			return err
		}

		formattedRoute, err := imports.Process(routeFile, buf.Bytes(), nil)
		if err != nil {
			return err
		}

		file, err := os.Create(routeFile)
		if err != nil {
			return err
		}

		if _, err = file.Write(formattedRoute); err != nil {
			return err
		}
	}

	if ctx.DoViews && ctx.Views != nil {
		for _, v := range ctx.Views {
			err := os.MkdirAll(filepath.Dir(v.File), os.ModePerm)
			if err != nil {
				return err
			}

			file, err := os.Create(v.File)
			if err != nil {
				return err
			}

			if _, err = file.WriteString(v.Body); err != nil {
				return err
			}
		}
	}

	return nil
}
