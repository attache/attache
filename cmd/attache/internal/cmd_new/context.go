package cmd_new

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/alecthomas/template"
	"github.com/iancoleman/strcase"
)

type Context struct {
	flags *flag.FlagSet
	Dir   string
	Name  string
}

func (c *Context) Execute(args []string) error { return c.do(args) }

func (c *Context) do(args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	c.flags = flag.NewFlagSet("new", flag.ContinueOnError)
	name := c.flags.String("n", "", "application name")
	if err := c.flags.Parse(args); err != nil {
		return err
	}

	if *name == "" {
		return errors.New("must provide name with -n")
	}

	c.Name = strcase.ToCamel(*name)
	c.Dir = strcase.ToSnake(*name)

	info, err := os.Stat(filepath.Join(cwd, c.Dir))
	if err == nil {
		if info.IsDir() {
			return fmt.Errorf("directory %s already exists", c.Dir)
		}

		return fmt.Errorf("%s already exists and is not a directory", c.Dir)
	}

	if !os.IsNotExist(err) {
		return err
	}

	buildErr := Dir{
		Name: c.Dir,
		Files: []File{
			{Name: ".gitignore", Body: []byte("secret")},
			{Name: "attache.json", BodyFunc: c.FileTemplate("attache.json.tpl")},
			{Name: "main.go", BodyFunc: c.FileTemplate("main.go.tpl")},
			{Name: "Taskfile.yml", BodyFunc: c.FileTemplate("Taskfile.yml.tpl")},
			{Name: "go.mod", BodyFunc: c.FileTemplate("go.mod.tpl")},
		},

		Dirs: []Dir{
			{Name: "views", Files: []File{
				{Name: "layout.tpl", BodyFunc: c.FileTemplate("layout.tpl.tpl")},
				{Name: "index.tpl", BodyFunc: c.FileTemplate("index.tpl.tpl")},
			}},
			{Name: "models"},
			{Name: "web", Dirs: []Dir{
				{Name: "dist", Dirs: []Dir{
					{Name: "js"},
					{Name: "css"},
					{Name: "img"},
				}},
				{Name: "src", Dirs: []Dir{
					{Name: "script"},
					{Name: "styles"},
				}},
			}},
			{Name: "secret", Files: []File{
				{Name: "schema.sql", Body: []byte("")},
				{Name: "dev.env", Body: []byte("DB_DRIVER=sqlite3\nDB_DSN=:memory:")},
				{Name: "run.sh", Body: []byte("go run *.go")},
			}},
		},
	}.Build(cwd)

	if buildErr != nil {
		return buildErr
	}

	return nil
}

func (c *Context) FileTemplate(name string) func() ([]byte, error) {
	fullPath := path.Join("templates", name)
	return func() ([]byte, error) {
		tpl := template.New("")
		if strings.HasSuffix(name, ".tpl.tpl") {
			tpl.Delims("[[", "]]")
		}

		tpl = template.Must(
			tpl.Parse(
				string(MustAsset(fullPath)),
			),
		)

		buf := new(bytes.Buffer)
		if err := tpl.Execute(buf, c); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	}
}
