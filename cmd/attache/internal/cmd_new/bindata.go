// Code generated by go-bindata. DO NOT EDIT.
// sources:
// templates/attache.json.tpl (34B)
// templates/index.tpl.tpl (150B)
// templates/layout.tpl.tpl (251B)
// templates/main.go.tpl (1.038kB)

package cmd_new

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _templatesAttacheJsonTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xe6\x52\x50\x50\x50\x50\x72\xce\xcf\x2b\x49\xad\x28\x09\xa9\x2c\x48\x55\xb2\x52\x50\xaa\xae\xd6\xf3\x4b\xcc\x4d\xad\xad\x55\xe2\xaa\x05\x04\x00\x00\xff\xff\x3e\xe1\xfc\xf4\x22\x00\x00\x00")

func templatesAttacheJsonTplBytes() ([]byte, error) {
	return bindataRead(
		_templatesAttacheJsonTpl,
		"templates/attache.json.tpl",
	)
}

func templatesAttacheJsonTpl() (*asset, error) {
	bytes, err := templatesAttacheJsonTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/attache.json.tpl", size: 34, mode: os.FileMode(0644), modTime: time.Unix(1606360425, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x25, 0x25, 0xde, 0xdd, 0x98, 0x95, 0x98, 0x11, 0x1b, 0x78, 0x73, 0xa, 0x6, 0x13, 0x27, 0x33, 0x5d, 0x51, 0x5f, 0x5, 0x9d, 0x1, 0x76, 0xfc, 0xdc, 0xb8, 0x27, 0x9c, 0x9c, 0xc, 0x79, 0xd2}}
	return a, nil
}

var _templatesIndexTplTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\xcc\x31\x0e\x42\x21\x10\x84\xe1\x9e\x53\xac\xdc\xc4\x53\x58\x23\xac\x71\x13\x64\xc9\x7b\x63\x41\x36\x73\x77\x9b\x97\xd8\x50\x7f\x33\x7f\x44\xd3\x97\x0d\x95\x0c\x43\xd7\x4c\x3e\xb4\x57\xff\x68\x84\x8e\x46\xa6\xff\xe0\xc4\xea\x7a\x66\x72\x43\xf5\xb0\x89\xbd\x3d\xbd\xad\x4c\xa6\xab\x2b\x70\x59\xfe\x3d\xe4\x0e\x94\xfa\x56\x29\x73\x76\xab\x05\xe6\xe3\x96\xae\xfb\x2f\x00\x00\xff\xff\x06\x47\x5f\x5d\x96\x00\x00\x00")

func templatesIndexTplTplBytes() ([]byte, error) {
	return bindataRead(
		_templatesIndexTplTpl,
		"templates/index.tpl.tpl",
	)
}

func templatesIndexTplTpl() (*asset, error) {
	bytes, err := templatesIndexTplTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/index.tpl.tpl", size: 150, mode: os.FileMode(0644), modTime: time.Unix(1606360428, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x99, 0xc8, 0xcd, 0xf4, 0x76, 0x6c, 0xa1, 0x6d, 0x9f, 0x55, 0x8e, 0xc9, 0x7a, 0xc9, 0xcf, 0x21, 0xf6, 0x59, 0x13, 0xa9, 0x59, 0x60, 0xfe, 0x33, 0x9c, 0x3b, 0xf, 0xd3, 0x1f, 0x7, 0x2d, 0xcd}}
	return a, nil
}

var _templatesLayoutTplTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb2\x51\x74\xf1\x77\x0e\x89\x0c\x70\x55\xc8\x28\xc9\xcd\xb1\xe3\xb2\x81\x50\x0a\x0a\x0a\x0a\x36\x19\xa9\x89\x29\x10\x26\x98\x5b\x92\x59\x92\x93\x6a\x57\x5d\x9d\x94\x93\x9f\x9c\xad\xa0\x04\xe6\x2a\x29\xe8\xd5\xd6\x46\x47\xeb\xf9\x25\xe6\xa6\xc6\xc6\x56\x57\xa7\xe6\xa5\xd4\xd6\xda\xe8\x43\x94\xc2\xb5\xc2\xf5\x14\x97\x54\xe6\xa4\x16\x83\x35\x41\xd5\x62\x51\x93\x5c\x94\x59\x50\x82\xa9\xc8\x46\x1f\xe1\x1e\x9b\xa4\xfc\x94\x4a\x2c\xe6\x83\x84\xc1\x1a\x91\xa4\x90\x4d\x80\x68\xb3\xd1\x07\x7b\x12\x10\x00\x00\xff\xff\xd6\xc2\xdf\x26\xfb\x00\x00\x00")

func templatesLayoutTplTplBytes() ([]byte, error) {
	return bindataRead(
		_templatesLayoutTplTpl,
		"templates/layout.tpl.tpl",
	)
}

func templatesLayoutTplTpl() (*asset, error) {
	bytes, err := templatesLayoutTplTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/layout.tpl.tpl", size: 251, mode: os.FileMode(0644), modTime: time.Unix(1606358908, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x66, 0xb7, 0xff, 0xed, 0xfe, 0xac, 0x1e, 0x7f, 0xdc, 0xb9, 0xe0, 0xea, 0x59, 0xaf, 0xd3, 0xf2, 0x4a, 0x4, 0xc2, 0x4f, 0x1f, 0x17, 0xaa, 0xcf, 0x3e, 0x90, 0xd, 0xa, 0x55, 0x13, 0xe0, 0x26}}
	return a, nil
}

var _templatesMainGoTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x53\xcd\x6e\xdb\x30\x0c\xbe\xfb\x29\x38\x1f\x36\x3b\x48\xad\xc3\x6e\x05\x7a\xe9\x92\x76\x05\xda\x75\x4b\x82\xed\x58\xc8\x32\xe3\x10\x93\x25\x45\xa2\xd3\x66\x41\xde\x7d\xf0\x6f\xd3\x22\xf5\xc1\x96\xc8\xef\x47\xa4\x29\x27\xd5\x5f\x59\x22\x54\x92\x4c\x14\x51\xe5\xac\x67\x48\x22\x00\x80\xd8\x20\x8b\x0d\xb3\x8b\xbb\xad\xb6\x65\xb7\x6a\x5f\x42\x40\x21\x59\xe6\x32\x20\x14\x9e\x76\xe8\xc3\x10\x4f\x3c\x56\x76\x87\xa0\xc9\x60\x80\xbd\xad\xa1\xb0\xe6\x0b\x83\x41\x2c\xd2\x16\xf4\x04\x71\x49\xbc\xa9\xf3\x4c\xd9\x4a\x54\x92\xd9\x88\xd2\x5e\x84\xad\x26\xc6\xaf\x71\x23\xb2\xec\xd6\x67\xe0\x1d\xf0\xa2\x33\x15\xd5\x3e\x6c\x75\xcb\x78\xd8\x2f\x7f\xdd\x9f\xc1\x6b\xca\x85\xdb\xb6\x90\x9f\x36\x70\xe9\xb1\xc1\x0d\x87\x95\xcc\x52\x6d\xb0\x2b\xf1\x84\xd5\xc7\x87\x6f\x1c\xa5\x51\xc4\x7b\x87\x70\x38\x64\x3f\x64\x85\xc7\x23\x04\xf6\xb5\x62\x38\x0c\x52\x1e\xb7\x35\x79\x2c\xda\x7d\xcf\xcb\xae\x65\xc0\x6f\xd6\x30\xbe\x30\x8c\xa6\x4a\x3a\x99\x93\x26\x26\x0c\x6f\xd0\x33\x5c\xcb\x5a\xf3\xdc\xec\xc8\x5b\x53\xa1\xe1\x73\xe9\x1b\xd2\xb8\x44\xbf\x43\x7f\x2e\xfb\x9b\xf0\x39\xbc\x2b\x6f\xc8\xcd\xae\x9b\x20\x1a\x99\x6b\x7c\xfd\x7f\xca\x1a\x83\x8a\x69\x47\xbc\xff\x80\xb8\xc4\x10\xc8\x9a\x13\x76\xe8\x23\x81\xad\x97\x25\x46\xc7\x28\x5a\xd7\x46\x41\xa2\x60\x32\xf6\x28\x85\x3b\x43\x9c\x3c\x43\x33\x47\xd9\x02\x83\xb3\x26\xe0\x1f\x4f\x8c\x7e\x0a\x1e\x26\x7d\x7c\x5b\x63\xe0\x74\x68\xe5\x04\x56\x8f\xb3\xc7\x4b\x20\x43\x4c\x52\xd3\xbf\xf6\x88\x6d\x0b\x27\xa2\x31\x12\x02\x6e\xe7\x2b\x10\x67\x1d\x6f\xe7\xab\xa7\x64\xd0\x52\x59\xb3\xbd\x33\x05\xbe\x24\xe9\x29\x95\x9a\xd0\x87\xfc\x9e\xd0\x8b\x0c\xbd\x58\xa0\x29\xd0\x7f\x5f\x3d\xdc\x27\x6a\x0a\x71\x2b\x11\xa7\x63\xe5\xcd\x1d\x1a\x39\x42\x40\x6e\x2d\x07\xf6\xd2\x81\x74\x4e\x93\x92\xdc\xf4\x6b\x6d\xfd\x58\xcd\xdb\x81\xea\xbc\x9c\x9b\x02\x7a\x0f\x97\x57\xaf\x33\x34\x08\x25\x9f\x47\xf0\xe1\xd8\x5d\x25\x5a\xb7\xe8\x4f\x57\x60\x48\xf7\xde\xcd\xa3\x6d\x99\xdd\x48\x96\x5a\x9b\x04\xbd\xef\xc0\xc7\x6e\x02\x4f\x73\xd2\xb9\x6c\x51\x9b\x24\x6d\xea\xf8\x1f\x00\x00\xff\xff\xf2\x94\xf9\x71\x0e\x04\x00\x00")

func templatesMainGoTplBytes() ([]byte, error) {
	return bindataRead(
		_templatesMainGoTpl,
		"templates/main.go.tpl",
	)
}

func templatesMainGoTpl() (*asset, error) {
	bytes, err := templatesMainGoTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/main.go.tpl", size: 1038, mode: os.FileMode(0644), modTime: time.Unix(1606360455, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x34, 0x29, 0x14, 0xb6, 0xf1, 0x36, 0x6c, 0x10, 0x22, 0xb0, 0x5a, 0xb7, 0x19, 0x6c, 0xe, 0x1c, 0xe6, 0xae, 0x57, 0x9a, 0x9a, 0xa9, 0x16, 0xec, 0xda, 0xbf, 0x9c, 0xb2, 0x18, 0x9a, 0xe7, 0xc2}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"templates/attache.json.tpl": templatesAttacheJsonTpl,
	"templates/index.tpl.tpl":    templatesIndexTplTpl,
	"templates/layout.tpl.tpl":   templatesLayoutTplTpl,
	"templates/main.go.tpl":      templatesMainGoTpl,
}

// AssetDebug is true if the assets were built with the debug flag enabled.
const AssetDebug = false

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"templates": {nil, map[string]*bintree{
		"attache.json.tpl": {templatesAttacheJsonTpl, map[string]*bintree{}},
		"index.tpl.tpl": {templatesIndexTplTpl, map[string]*bintree{}},
		"layout.tpl.tpl": {templatesLayoutTplTpl, map[string]*bintree{}},
		"main.go.tpl": {templatesMainGoTpl, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory.
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
