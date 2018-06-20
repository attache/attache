package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

type Dir struct {
	Name  string
	Dirs  []Dir
	Files []File
}

func (d Dir) Build(root string) error {
	self := filepath.Join(root, d.Name)
	if err := os.MkdirAll(self, os.ModePerm); err != nil {
		return err
	}

	for _, file := range d.Files {
		if err := file.Build(self); err != nil {
			return err
		}
	}

	for _, dir := range d.Dirs {
		if err := dir.Build(self); err != nil {
			return err
		}
	}

	return nil
}

type File struct {
	Name     string
	Body     []byte
	BodyFunc func() ([]byte, error)
}

func (f File) Build(root string) error {
	var body []byte
	if f.BodyFunc != nil {
		data, err := f.BodyFunc()
		if err != nil {
			return err
		}
		body = data
	} else {
		body = f.Body
	}

	if body != nil {
		if err := ioutil.WriteFile(filepath.Join(root, f.Name), body, os.ModePerm); err != nil {
			return err
		}
	} else {
		log.Printf("skipping %s (empty file)", filepath.Join(root, f.Name))
	}

	return nil
}

func FileTemplate(name string) func() ([]byte, error) {
	return func() ([]byte, error) {
		tpl := template.New("")
		if strings.HasSuffix(name, ".tpl.tpl") {
			tpl.Delims("[[", "]]")
		}

		tpl = template.Must(
			tpl.Parse(
				MustAssetString(path.Join("templates", name)),
			),
		)

		buf := new(bytes.Buffer)
		if err := tpl.Execute(buf, context); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	}
}
