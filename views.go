package attache

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"html/template"
)

type View interface {
	Execute(out io.Writer, data interface{}) error
}

type noopView struct{}

func (noopView) Execute(_ io.Writer, _ interface{}) error { return nil }

type ViewCache interface {
	Get(name string) View
	Has(name string) (ok bool)
	Render(name string, data interface{}) ([]byte, error)

	// ViewCache should not be implemented outside of this package
	private()
}
type viewCache map[string]View

func (v viewCache) private() {}

func (v viewCache) Has(name string) (ok bool) { return v != nil && v[name] != nil }

func (v viewCache) Get(name string) View {
	if v.Has(name) {
		return v[name]
	}

	return noopView{}
}

func (v viewCache) Render(name string, data interface{}) ([]byte, error) {
	buf := getbuf()
	defer putbuf(buf)

	view := v.Get(name)
	if err := view.Execute(buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (v viewCache) load(path, prefix string, layouts []string) error {
	stats, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	layout, err := ioutil.ReadFile(filepath.Join(path, d_LAYOUT_FILE))
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

		if file.Name() == d_LAYOUT_FILE {
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

		v[tplName] = tpl
	}

	for _, dir := range subdirs {
		var newPrefix string
		if prefix == "" {
			newPrefix = filepath.Base(dir)
		} else {
			newPrefix = fmt.Sprintf("%s.%s", prefix, filepath.Base(dir))
		}
		if err := v.load(dir, newPrefix, layouts); err != nil {
			return err
		}
	}

	return nil
}
