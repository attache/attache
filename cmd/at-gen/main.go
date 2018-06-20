package main

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"golang.org/x/tools/imports"
)

//go:generate go-bindata templates

func main() {
	log.SetFlags(0)

	var ctx Context
	ctx.Init()

	// prevent overwrite of existing model file
	modelFile := filepath.Join("models", ctx.Model.Table+".go")
	if ctx.DoModel {
		// make sure file doesn;t already exist
		if _, err := os.Stat(modelFile); err == nil {
			if !ctx.Replace {
				log.Fatalln(modelFile, "already exists")
			}
		} else if !os.IsNotExist(err) {
			log.Fatalln(err)
		}

		// ensure containing directory exists
		if err := os.MkdirAll("models", os.ModePerm); err != nil {
			log.Fatalln(err)
		}
	}

	if ctx.DoViews {
		// prevent overwrite of view files
		for _, v := range ctx.Views {
			if _, err := os.Stat(v.File); err == nil {
				if !ctx.Replace {
					log.Fatalln(v.File, "already exists")
				}
			} else if !os.IsNotExist(err) {
				log.Fatalln(err)
			}
		}
	}

	routeFile := ctx.Model.Table + "_routes.go"
	if ctx.DoRoutes {
		// prevent overwrite of route file
		if _, err := os.Stat(routeFile); err == nil {
			if !ctx.Replace {
				log.Fatalln(routeFile, "already exists")
			}
		} else if !os.IsNotExist(err) {
			log.Fatalln(err)
		}
	}

	var buf bytes.Buffer

	if ctx.DoModel {
		buf.Reset()

		// generate model file
		tpl, err := template.New("").Parse(MustAssetString("templates/model.tpl"))
		if err != nil {
			log.Fatalln(err)
		}

		if err = tpl.Execute(&buf, ctx.Model); err != nil {
			log.Fatalln(err)
		}

		formattedModel, err := imports.Process(modelFile, buf.Bytes(), nil)
		if err != nil {
			log.Fatalln(err)
		}

		file, err := os.Create(modelFile)
		if err != nil {
			log.Fatalln(err)
		}

		if _, err = file.Write(formattedModel); err != nil {
			log.Fatalln(err)
		}
	}

	// generate route file
	if ctx.DoRoutes {
		buf.Reset()

		tpl, err := template.New("").Parse(MustAssetString("templates/routes.tpl"))
		if err != nil {
			log.Fatalln(err)
		}

		if err = tpl.Execute(&buf, ctx.Model); err != nil {
			log.Fatalln(err)
		}

		formattedRoute, err := imports.Process(routeFile, buf.Bytes(), nil)
		if err != nil {
			log.Fatalln(err)
		}

		file, err := os.Create(routeFile)
		if err != nil {
			log.Fatalln(err)
		}

		if _, err = file.Write(formattedRoute); err != nil {
			log.Fatalln(err)
		}
	}

	if ctx.DoViews && ctx.Views != nil {
		for _, v := range ctx.Views {
			err := os.MkdirAll(filepath.Dir(v.File), os.ModePerm)
			if err != nil {
				log.Fatalln(err)
			}

			file, err := os.Create(v.File)
			if err != nil {
				log.Fatalln(err)
			}

			if _, err = file.WriteString(v.Body); err != nil {
				log.Fatalln(err)
			}
		}
	}
}
