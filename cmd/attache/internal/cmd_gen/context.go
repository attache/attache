package cmd_gen

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"

	"github.com/attache/attache/cmd/attache/shared"
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

	ScopeCamel, ScopeSnake, ScopePath string

	Replace bool

	JSONRoutes bool
}

type fieldDefs []string

func (f fieldDefs) String() string      { return strings.Join(f, " ") }
func (f *fieldDefs) Set(s string) error { *f = append(*f, s); return nil }

func (c *Context) recover(err *error) {
	val := recover()
	if val != nil {
		*err = fmt.Errorf("%v", val)
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
	c.flags.BoolVar(&c.JSONRoutes, "json", false, "create JSON endpoints")
	name := c.flags.String("n", "", "name of the resource")
	table := c.flags.String("t", "", "name of table")
	scope := c.flags.String("s", "", "scope to generate routes and views under")
	defs := &fieldDefs{}
	c.flags.Var(defs, "f", "-f NAME:TYPE:FLAGs [...]")

	if err := c.flags.Parse(args); err != nil {
		return err
	}

	if c.ContextType == "" {
		return errors.New("ctx cannot be \"\"")
	}

	if *scope != "" {
		c.ScopeCamel = strcase.ToCamel(*scope)
		c.ScopeSnake = strcase.ToSnake(*scope)
		c.ScopePath = path.Clean("/" + strings.Replace(c.ScopeSnake, "_", "/", -1))
	}

	if *name == "" {
		return errors.New("name cannot be \"\"")
	}

	if len(*defs) == 0 {
		return errors.New("must specify at least one field")
	}

	noSpec := false
	// if none specified, do all
	if !c.DoModel && !c.DoViews && !c.DoRoutes {
		c.DoModel = true
		c.DoViews = true
		c.DoRoutes = true
		noSpec = true
	}

	if c.JSONRoutes {
		// implied
		c.DoRoutes = true

		if c.DoViews {
			if noSpec {
				// if the user didn't explicitly request views,
				// we'll assume they don't need them for JSON
				// endpoints, and so will skip their generation
				c.DoViews = false
			} else {
				// if the user DID specifically request views,
				// we'll still generate them but will give a
				// warning that they're not used by default
				// when generating JSON endpoints
				log.Println("warning: views are unused when generating JSON routes")
			}
		}
	}

	// needed for more than just model
	c.Model = buildModel(*name, *table, *defs)

	if c.DoViews {
		createForm := template.Must(template.New("").Delims("[[", "]]").
			Parse(MustAssetString("templates/view_create.tpl")))
		updateForm := template.Must(template.New("").Delims("[[", "]]").
			Parse(MustAssetString("templates/view_update.tpl")))
		listView := template.Must(template.New("").Delims("[[", "]]").
			Parse(MustAssetString("templates/view_list.tpl")))

		var (
			create = &strings.Builder{}
			update = &strings.Builder{}
			list   = &strings.Builder{}
		)

		if err := createForm.Execute(create, c); err != nil {
			return err
		}

		if err := updateForm.Execute(update, c); err != nil {
			return err
		}

		if err := listView.Execute(list, c); err != nil {
			return err
		}

		c.Views = []View{
			View{
				File: filepath.Join("views", c.ScopeSnake, c.Model.Table, "create.tpl"),
				Body: create.String(),
			},

			View{
				File: filepath.Join("views", c.ScopeSnake, c.Model.Table, "update.tpl"),
				Body: update.String(),
			},

			View{
				File: filepath.Join("views", c.ScopeSnake, c.Model.Table, "list.tpl"),
				Body: list.String(),
			},
		}
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

	var routeFile string
	if ctx.ScopeSnake != "" {
		routeFile = fmt.Sprintf("%s_%s_routes.go", ctx.ScopeSnake, ctx.Model.Table)
	} else {
		routeFile = fmt.Sprintf("%s_routes.go", ctx.Model.Table)
	}

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

	var (
		buf     bytes.Buffer
		created []string
	)

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
			return err
		}

		file, err := os.Create(modelFile)
		if err != nil {
			return err
		}

		if _, err = file.Write(formattedModel); err != nil {
			return err
		}

		created = append(created, modelFile)
	}

	// generate route file
	if ctx.DoRoutes {
		buf.Reset()

		var (
			tpl *template.Template
			err error
		)

		if ctx.JSONRoutes {
			tpl, err = template.New("").Parse(MustAssetString("templates/routes.json.tpl"))
		} else {
			tpl, err = template.New("").Parse(MustAssetString("templates/routes.tpl"))
		}

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

		created = append(created, routeFile)
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

			created = append(created, v.File)
		}
	}

	fmt.Printf("done. files created:\n\t%s\n", strings.Join(created, "\n\t"))

	return nil
}
