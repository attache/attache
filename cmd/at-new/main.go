package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/iancoleman/strcase"
)

//go:generate go-bindata templates

var context struct {
	Dir  string
	Name string
}

func main() {
	log.SetFlags(0)

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("fatal:", err)
	}

	name := flag.String("n", "App", "application name")
	flag.Parse()

	context.Name = strcase.ToCamel(*name)
	context.Dir = strcase.ToSnake(*name)

	info, err := os.Stat(filepath.Join(cwd, context.Dir))
	if err == nil {
		if info.IsDir() {
			log.Fatalf("error: %s already exists and is not a directory", context.Dir)
		}

		log.Fatalf("error: directory %s already exists", context.Dir)
	}

	if !os.IsNotExist(err) {
		log.Fatalln("fatal:", err)
	}

	buildErr := Dir{
		Name: context.Dir,
		Files: []File{
			{Name: "main.go", BodyFunc: FileTemplate("main.go.tpl")},
		},

		Dirs: []Dir{
			{Name: "views", Files: []File{
				{Name: "layout.tpl", BodyFunc: FileTemplate("layout.tpl.tpl")},
				{Name: "index.tpl", BodyFunc: FileTemplate("index.tpl.tpl")},
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
			}},
		},
	}.Build(cwd)

	if buildErr != nil {
		log.Fatalln(buildErr)
	}
}
