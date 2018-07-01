package cmd_new

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/imports"
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
		var (
			err      = error(nil)
			filePath = filepath.Join(root, f.Name)
		)

		if strings.HasSuffix(f.Name, ".go") {
			// apply goimports to go files
			body, err = imports.Process(filePath, body, nil)
			if err != nil {
				return err
			}
		}

		if err = ioutil.WriteFile(filePath, body, os.ModePerm); err != nil {
			return err
		}
	} else {
		log.Printf("skipping %s (empty file)", filepath.Join(root, f.Name))
	}

	return nil
}
