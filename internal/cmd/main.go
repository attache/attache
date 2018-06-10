package main

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"golang.org/x/tools/imports"
)

func main() {
	log.SetFlags(0)

	var ctx Context
	ctx.Init()

	// prevent overwrite of existing model file
	modelFile := filepath.Join("models", ctx.Model.Table+".go")
	if _, err := os.Stat(modelFile); err == nil {
		log.Fatalln(modelFile, "already exists")
	} else if !os.IsNotExist(err) {
		log.Fatalln(err)
	}

	// prevent overwrite of view files
	for _, v := range ctx.Views {
		if _, err := os.Stat(v.File); err == nil {
			log.Fatalln(v.File, "already exists")
		} else if !os.IsNotExist(err) {
			log.Fatalln(err)
		}
	}

	var buf bytes.Buffer

	tpl, err := template.New("").Parse(modelTemplate)
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

	if ctx.Views != nil {
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
