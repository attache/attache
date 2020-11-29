package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/imports"
)

// A Dir represents a directory within a directory structure to generate.
type Dir struct {
	Name  string
	Dirs  []Dir
	Files []File
	Perm  os.FileMode
}

// Build creates the directory structure at the appropriate location relative to the root.
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

	if err := os.Chmod(self, getPerm(d.Perm, true)); err != nil {
		return err
	}

	return nil
}

// A File represents a file within a directory structure to generate.
type File struct {
	Name     string
	Body     []byte
	BodyFunc func() ([]byte, error)
	Perm     os.FileMode
}

// Build creates a file at the appropriate location relative to the root.
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

	var (
		err      error
		filePath = filepath.Join(root, f.Name)
	)

	if strings.HasSuffix(f.Name, ".go") {
		// apply goimports to go files
		body, err = imports.Process(filePath, body, nil)
		if err != nil {
			return err
		}
	}

	if err = ioutil.WriteFile(filePath, body, getPerm(f.Perm, false)); err != nil {
		return err
	}

	return nil
}

func getPerm(given os.FileMode, isDirectory bool) os.FileMode {
	permBits := given & os.ModePerm
	if permBits == 0 {
		if isDirectory {
			permBits = 0o755
		} else {
			permBits = 0o644
		}
	}
	return permBits
}
