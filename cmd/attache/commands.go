package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"golang.org/x/tools/imports"
)

//go:generate go run github.com/kevinburke/go-bindata/go-bindata -pkg main templates/new templates/gen

// CommandNew powers the `attache new` command.
type CommandNew struct {
	Name         string `arg:"positional,required"`
	API          bool   `arg:"-j,--api" default:"false"`
	Dir          string `arg:"-d,--dir" default:""`
	LocalAttache string `arg:"-l,--local-attache" default:"" help:"use a local version of attache"`
	Version      string `arg:"-"`
}

// ProjectStructure calculates the project structure to generate.
func (c *CommandNew) ProjectStructure() (dir Dir) {
	dir.Name = c.Dir
	dir.Dirs = append(
		dir.Dirs,
		Dir{Name: "models"},
		Dir{Name: "secret", Files: []File{
			{Name: "dev.env", BodyFunc: c.FileTemplate("dev.env.tpl")},
		}},
	)
	if !c.API {
		dir.Dirs = append(
			dir.Dirs,
			Dir{Name: "views", Files: []File{
				{Name: "layout.go.html", BodyFunc: c.FileTemplate("layout.go.html.tpl")},
				{Name: "index.go.html", BodyFunc: c.FileTemplate("index.go.html.tpl")},
			}},
			Dir{Name: "web", Dirs: []Dir{
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
		)
	}
	dir.Files = append(
		dir.Files,
		File{Name: ".gitignore", BodyFunc: c.FileTemplate("gitignore.tpl")},
		File{Name: ".prettierignore", BodyFunc: c.FileTemplate("prettierignore.tpl")},
		File{Name: "attache.json", BodyFunc: c.FileTemplate("attache.json.tpl")},
		File{Name: "main.go", BodyFunc: c.FileTemplate("main.go.tpl")},
		File{Name: "Taskfile.yml", BodyFunc: c.FileTemplate("Taskfile.yml.tpl")},
		File{Name: "go.mod", BodyFunc: c.FileTemplate("go.mod.tpl")},
	)
	return dir
}

// Execute runs the command.
func (c *CommandNew) Execute() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	c.Name = strcase.ToCamel(c.Name)
	if c.Dir == "" {
		c.Dir = strcase.ToSnake(c.Name)
	}

	projectDir := path.Join(cwd, c.Dir)

	if err := c.ProjectStructure().Build(cwd); err != nil {
		if remErr := os.RemoveAll(projectDir); remErr != nil {
			log.Println("failed to clean up", path.Join(cwd, c.Dir), ":", remErr)
		}
		return fmt.Errorf("failed to create project: %w", err)
	}

	modCmd := exec.Command("go", "mod", "tidy")
	modCmd.Dir = projectDir
	modCmd.Stderr = os.Stderr
	modCmd.Stdout = os.Stdout
	if err := modCmd.Run(); err != nil {
		return fmt.Errorf("failed to run `go mod tidy`: %w", err)
	}

	fmtCmd := exec.Command("goimports", "-w", ".")
	fmtCmd.Dir = projectDir
	fmtCmd.Stderr = os.Stderr
	fmtCmd.Stdout = os.Stdout
	if err := fmtCmd.Run(); err != nil {
		return fmt.Errorf("failed to run `goimports`: %w", err)
	}

	fmt.Println("all done! run `cd " + c.Dir + "`")
	return nil
}

// FileTemplate renders the named embedded template to a byte slice.
func (c *CommandNew) FileTemplate(name string) func() ([]byte, error) {
	return func() ([]byte, error) {
		var buf bytes.Buffer
		asset, err := AssetString(path.Join("templates", "new", name))
		if err != nil {
			return nil, err
		}
		tpl := template.New("")
		if strings.HasSuffix(name, ".go.html.tpl") {
			tpl.Delims("[[", "]]")
		}
		if tpl, err = tpl.Parse(asset); err != nil {
			return nil, err
		}
		if err := tpl.Execute(&buf, c); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
}

// CommandGen powers the `attache gen` commands
type CommandGen struct {
	Name        string   `arg:"positional,required" help:"name of the model"`
	Table       string   `arg:"-t,--table" help:"name of the table"`
	GenViews    bool     `arg:"-v,--views" help:"generate views" default:"false"`
	GenRoutes   bool     `arg:"-r,--routes" help:"generate routes" default:"false"`
	GenModel    bool     `arg:"-m,--model" help:"generate model" default:"false"`
	JSONRoutes  bool     `arg:"-j,--json" help:"generate json routes" default:"false"`
	Replace     bool     `arg:"--replace" help:"replace existing files" default:"false"`
	Fields      []string `arg:"-f,--field,separate" help:"field specifications of the form NAME[:TYPE[:FLAG[,FLAG...]]]"`
	ContextType string   `arg:"--context" help:"name of the context type to generate routes for"`
	Scope       string   `arg:"-s,--scope" help:"scope under which to generate routes and views"`

	Model      *Model      `arg:"-"`
	Views      [][2]string `arg:"-"`
	NameSnake  string      `arg:"-"`
	TableCamel string      `arg:"-"`
	ScopeCamel string      `arg:"-"`
	ScopeSnake string      `arg:"-"`
	ScopePath  string      `arg:"-"`
}

func (c *CommandGen) Execute() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	conf, err := GetConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	c.JSONRoutes = c.JSONRoutes || conf.GetBool("APIOnly")
	if c.ContextType == "" {
		c.ContextType = conf.GetString("ContextType")
	}

	// Not specifying any specific generators implies all generators
	if !(c.GenModel || c.GenViews || c.GenRoutes) {
		c.GenModel = true
		c.GenRoutes = true
		c.GenViews = !c.JSONRoutes
	}

	if c.GenViews && c.JSONRoutes {
		if conf.GetBool("APIOnly") {
			return fmt.Errorf("failed to generate: cannot generate views for an api-only application")
		}
		return fmt.Errorf("failed to generate: cannot specify both --json and --views")
	}

	if c.Scope != "" {
		c.ScopeCamel = strcase.ToCamel(c.Scope)
		c.ScopeSnake = strcase.ToSnake(c.Scope)
		c.ScopePath = path.Clean("/" + path.Join(strings.Split(c.ScopeSnake, "_")...))
	}

	c.Name = strcase.ToCamel(c.Name)
	if c.Table == "" {
		c.Table = strcase.ToSnake(c.Name)
	} else {
		c.Table = strcase.ToSnake(c.Table)
	}

	c.TableCamel = strcase.ToCamel(c.Table)
	c.NameSnake = strcase.ToSnake(c.Name)

	c.Model, err = buildModel(c.Fields)
	if err != nil {
		return fmt.Errorf("invalid model definition: %w", err)
	}

	// prevent overwrite of existing model file
	modelFile := filepath.Join("models", c.NameSnake+".go")
	if c.GenModel {
		// make sure file doesn;t already exist
		if _, err := os.Stat(modelFile); err == nil {
			if !c.Replace {
				return fmt.Errorf("%s already exists", modelFile)
			}
		} else if !os.IsNotExist(err) {
			return err
		}

		// ensure containing directory exists
		if err := (Dir{Name: "models"}).Build(cwd); err != nil {
			return err
		}
	}

	if c.GenViews {
		createForm := template.Must(template.New("").Delims("[[", "]]").
			Parse(string(MustAsset("templates/gen/view_create.go.html.tpl"))))
		updateForm := template.Must(template.New("").Delims("[[", "]]").
			Parse(string(MustAsset("templates/gen/view_update.go.html.tpl"))))
		listView := template.Must(template.New("").Delims("[[", "]]").
			Parse(string(MustAsset("templates/gen/view_list.go.html.tpl"))))

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

		c.Views = [][2]string{
			{
				filepath.Join("views", c.ScopeSnake, c.NameSnake, "create.go.html"),
				create.String(),
			},
			{
				filepath.Join("views", c.ScopeSnake, c.NameSnake, "update.go.html"),
				update.String(),
			},
			{
				filepath.Join("views", c.ScopeSnake, c.NameSnake, "list.go.html"),
				list.String(),
			},
		}

		// prevent overwrite of view files
		for _, v := range c.Views {
			if _, err := os.Stat(v[0]); err == nil {
				if !c.Replace {
					return fmt.Errorf("%s already exists", v[0])
				}
			} else if !os.IsNotExist(err) {
				return err
			}
		}
	}

	var routeFile string
	if c.ScopeSnake != "" {
		routeFile = fmt.Sprintf("%s_%s_routes.go", c.ScopeSnake, c.NameSnake)
	} else {
		routeFile = fmt.Sprintf("%s_routes.go", c.NameSnake)
	}

	if c.GenRoutes {
		// prevent overwrite of route file
		if _, err := os.Stat(routeFile); err == nil {
			if !c.Replace {
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

	if c.GenModel {
		buf.Reset()

		// generate model file
		tpl, err := template.New("").Parse(string(MustAsset("templates/gen/model.go.tpl")))
		if err != nil {
			return err
		}

		if err = tpl.Execute(&buf, c); err != nil {
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
	if c.GenRoutes {
		buf.Reset()

		var (
			tpl *template.Template
			err error
		)

		if c.JSONRoutes {
			tpl, err = template.New("").Parse(string(MustAsset("templates/gen/routes.json.go.tpl")))
		} else {
			tpl, err = template.New("").Parse(string(MustAsset("templates/gen/routes.go.tpl")))
		}

		if err != nil {
			return err
		}

		if err = tpl.Execute(&buf, c); err != nil {
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

	if c.GenViews && c.Views != nil {
		for _, v := range c.Views {
			err := os.MkdirAll(filepath.Dir(v[0]), os.ModePerm)
			if err != nil {
				return err
			}

			file, err := os.Create(v[0])
			if err != nil {
				return err
			}

			if _, err = file.WriteString(v[1]); err != nil {
				return err
			}

			created = append(created, v[0])
		}
	}

	fmt.Printf("done. files created:\n\t%s\n", strings.Join(created, "\n\t"))

	return nil
}
