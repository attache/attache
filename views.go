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

var global_caches = map[string]ViewCache{}

func viewsForRoot(root string) (ViewCache, error) {
	root = filepath.Clean(root)

	if cached, ok := global_caches[root]; ok {
		return cached, nil
	}

	v := ViewCache{cache: map[string]View{}}
	if err := v.load(root, "", nil); err != nil {
		return v, err
	}

	global_caches[root] = v
	return v, nil
}

type View interface {
	Execute(out io.Writer, data interface{}) error
}

type noopView struct{}

func (noopView) Execute(_ io.Writer, _ interface{}) error { return nil }

type ViewCache struct{ cache map[string]View }

func (v ViewCache) Lookup(name string) (view View, ok bool) {
	if v.cache == nil {
		return noopView{}, false
	}

	if cached := v.cache[name]; cached != nil {
		return cached, true
	}

	return noopView{}, false
}

func (v ViewCache) Get(name string) (view View) {
	view, _ = v.Lookup(name)
	return
}

func (v ViewCache) Render(name string, data interface{}) ([]byte, error) {
	buf := getbuf()
	defer putbuf(buf)

	view := v.Get(name)
	if err := view.Execute(buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (v ViewCache) load(path, prefix string, layouts []string) error {
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

		v.cache[tplName] = tpl
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
