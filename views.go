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

// View is the interface implemented by a type that can be rendered
type View interface {
	Execute(out io.Writer, data interface{}) error
}

type noopView struct{}

func (noopView) Execute(_ io.Writer, _ interface{}) error { return nil }

// A ViewCache is a read-only view of a set of cached views. It is safe for
// concurrent use between goroutines
type ViewCache interface {
	// Get retrieves the View from the cache mapped to the given name.
	// If there is no view associated with name, an empty, no-op View
	// implementation is returned. The returned View will never be nil
	Get(name string) View

	// Has returns true if there is a View in the cache mapped to
	// the given name, otherwise it returns false
	Has(name string) (ok bool)

	// Render will call Get(name) to retrieve a View, and will
	// then execute that view with data as the View's data argument.
	// If an error is encountered, it is returned. If not, the rendered
	// bytes are returned
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
