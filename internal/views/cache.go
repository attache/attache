package views

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"html/template"
)

var cache = map[string]View{
	// starts empty
}

func init() {
	err := loadDir(RootDirectory, "", nil)

	if err != nil && !os.IsNotExist(err) {
		log.Fatalln("loading views:", err)
	}
}

func loadDir(path, prefix string, layouts []string) error {
	stats, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	layout, err := ioutil.ReadFile(filepath.Join(path, LayoutFilename))
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if layout != nil {
		layouts = append(layouts, string(layout))
	}

	var subdirs []string

	for _, file := range stats {
		fpath := filepath.Join(path, file.Name())

		if file.IsDir() {
			subdirs = append(subdirs, fpath)
			continue
		}

		if file.Name() == LayoutFilename {
			continue
		}

		tpl := template.New("")
		for _, l := range layouts {
			if _, err := tpl.Parse(l); err != nil {
				return err
			}
		}

		if _, err := tpl.ParseFiles(fpath); err != nil {
			return err
		}

		var tplName string

		if prefix == "" {
			tplName = strings.TrimSuffix(file.Name(), ".tpl")
		} else {
			tplName = fmt.Sprintf("%s.%s", prefix, strings.TrimSuffix(file.Name(), ".tpl"))
		}

		cache[tplName] = tpl
	}

	for _, dir := range subdirs {
		var newPrefix string
		if prefix == "" {
			newPrefix = filepath.Base(dir)
		} else {
			newPrefix = fmt.Sprintf("%s.%s", prefix, filepath.Base(dir))
		}
		if err := loadDir(dir, newPrefix, layouts); err != nil {
			return err
		}
	}

	return nil
}
