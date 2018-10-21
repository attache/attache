package attache

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"html/template"

	"github.com/Masterminds/sprig"
	viewDriver "github.com/attache/attache/drivers/view"
)

func init() {
	driver := attacheViews{}
	viewDriver.DriverRegister("", driver)
	viewDriver.DriverRegister("attache", driver)
}

type attacheViews struct{}

func (a attacheViews) Init(root string) (*viewDriver.Cache, error) {
	cache := new(viewDriver.Cache)
	if err := a.load(cache, root, "", nil); err != nil {
		return nil, err
	}
	return cache, nil
}

type viewLayout struct {
	File, Body string
}

var viewParseErrRx = regexp.MustCompile("^template: (?:[^:]*):([0-9]+): (.*)")

func viewParseError(currentFile string, e error) error {
	result := viewParseErrRx.FindStringSubmatch(e.Error())
	if result != nil {
		return fmt.Errorf("parse %s: line %s: %s", currentFile, result[1], result[2])
	}
	return fmt.Errorf("parse %s: %s", currentFile, e)
}

func (a attacheViews) load(c *viewDriver.Cache, path, prefix string, layouts []viewLayout) error {
	const layoutFileName = "layout.tpl"

	stats, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	layoutPath := filepath.Join(path, layoutFileName)
	layout, err := ioutil.ReadFile(layoutPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if layout != nil {
		layouts = append(layouts, viewLayout{
			File: layoutPath,
			Body: string(layout),
		})
	}

	var subdirs []string

	for _, file := range stats {
		fpath := filepath.Join(path, file.Name())

		if file.IsDir() {
			subdirs = append(subdirs, fpath)
			continue
		}

		if file.Name() == layoutFileName {
			continue
		}

		tpl := template.New("").Funcs(sprig.FuncMap())
		for _, l := range layouts {
			if _, err := tpl.Parse(l.Body); err != nil {
				return viewParseError(l.File, err)
			}
		}

		if _, err := tpl.ParseFiles(fpath); err != nil {
			return viewParseError(fpath, err)
		}

		var tplName string

		if prefix == "" {
			tplName = strings.TrimSuffix(file.Name(), ".tpl")
		} else {
			tplName = fmt.Sprintf("%s.%s", prefix, strings.TrimSuffix(file.Name(), ".tpl"))
		}

		c.Put(tplName, tpl)
	}

	for _, dir := range subdirs {
		var newPrefix string
		if prefix == "" {
			newPrefix = filepath.Base(dir)
		} else {
			newPrefix = fmt.Sprintf("%s.%s", prefix, filepath.Base(dir))
		}
		if err := a.load(c, dir, newPrefix, layouts); err != nil {
			return err
		}
	}

	return nil
}
