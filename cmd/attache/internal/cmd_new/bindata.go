// Code generated by go-bindata. DO NOT EDIT.
// sources:
// templates/.at-conf.json.tpl
// templates/index.tpl.tpl
// templates/layout.tpl.tpl
// templates/main.go.tpl

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
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
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

var _templatesAtConfJsonTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xe6\x52\x50\x50\x50\x50\x72\xce\xcf\x2b\x49\xad\x28\x09\xa9\x2c\x48\x55\xb2\x52\x50\xaa\xae\xd6\xf3\x4b\xcc\x4d\xad\xad\x55\xe2\xaa\x05\x04\x00\x00\xff\xff\x3e\xe1\xfc\xf4\x22\x00\x00\x00")

func templatesAtConfJsonTplBytes() ([]byte, error) {
	return bindataRead(
		_templatesAtConfJsonTpl,
		"templates/.at-conf.json.tpl",
	)
}

func templatesAtConfJsonTpl() (*asset, error) {
	bytes, err := templatesAtConfJsonTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/.at-conf.json.tpl", size: 34, mode: os.FileMode(420), modTime: time.Unix(1530473476, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x25, 0x25, 0xde, 0xdd, 0x98, 0x95, 0x98, 0x11, 0x1b, 0x78, 0x73, 0xa, 0x6, 0x13, 0x27, 0x33, 0x5d, 0x51, 0x5f, 0x5, 0x9d, 0x1, 0x76, 0xfc, 0xdc, 0xb8, 0x27, 0x9c, 0x9c, 0xc, 0x79, 0xd2}}
	return a, nil
}

var _templatesIndexTplTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x64\xcc\x41\x0a\xc2\x40\x0c\x85\xe1\xbd\xa7\x18\x72\x00\xe7\x02\xd6\x6b\xb8\xef\xcc\x2b\x0d\xc6\x14\x9a\x80\x96\x90\xbb\x0b\x83\xe0\xa2\xdb\xef\x3d\xfe\x88\x8e\x85\x15\x85\x9c\x5d\x40\x99\x0f\x48\xdb\x5e\x88\x80\xf6\xcc\xcb\xff\x60\x7e\x08\x8c\x32\x6f\xc2\xfa\x2c\x3b\x64\xfa\xd9\x0a\x38\x95\x75\xc7\x32\x51\x7d\x63\xae\xcd\xac\xb2\x76\x7c\xae\xcd\x8c\xea\xfd\x1c\x9b\xb7\x7e\xd0\x90\xb1\x7c\x03\x00\x00\xff\xff\xa6\x6a\xf6\x2b\x87\x00\x00\x00")

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

	info := bindataFileInfo{name: "templates/index.tpl.tpl", size: 135, mode: os.FileMode(420), modTime: time.Unix(1530470276, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xe1, 0x45, 0xc1, 0xcb, 0x9e, 0x77, 0xf6, 0x60, 0x29, 0xc9, 0x27, 0xae, 0x44, 0xbe, 0xcc, 0x7b, 0xe6, 0xb1, 0x89, 0x2c, 0x65, 0xd6, 0xb4, 0x4a, 0xa3, 0x3c, 0x75, 0x69, 0x9c, 0x17, 0x5, 0x8}}
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

	info := bindataFileInfo{name: "templates/layout.tpl.tpl", size: 251, mode: os.FileMode(420), modTime: time.Unix(1530470276, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x66, 0xb7, 0xff, 0xed, 0xfe, 0xac, 0x1e, 0x7f, 0xdc, 0xb9, 0xe0, 0xea, 0x59, 0xaf, 0xd3, 0xf2, 0x4a, 0x4, 0xc2, 0x4f, 0x1f, 0x17, 0xaa, 0xcf, 0x3e, 0x90, 0xd, 0xa, 0x55, 0x13, 0xe0, 0x26}}
	return a, nil
}

var _templatesMainGoTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x92\xcd\x6e\xdb\x30\x10\x84\xef\x7c\x8a\xad\x0e\x05\x25\x08\xe2\x3d\x40\x2e\xa9\x93\x34\x40\xdb\x00\x8e\xd0\x1e\x0b\x9a\x5e\xdb\xdb\x50\x24\x43\xae\xea\xa4\x86\xde\xbd\xd0\x8f\x95\x1f\xc4\x87\xe8\x24\xed\xce\xf0\xa3\x06\x13\xb4\xb9\xd7\x5b\x84\x46\x93\x13\x82\x9a\xe0\x23\x83\x14\x00\x00\x99\x43\x56\x3b\xe6\x90\x8d\x9f\xd6\x6f\x33\x31\xbe\x6e\x89\x77\xed\xaa\x32\xbe\x51\x8d\x31\xde\xda\x3f\x51\x69\x66\x6d\x76\x98\x89\x5c\x08\x7e\x0a\x08\x87\x43\xf5\x43\x37\xd8\x75\x90\x38\xb6\x86\xe1\x30\x98\x55\x01\x11\x1f\x5a\x8a\xb8\x86\x42\x0d\xa3\xc9\x5a\x5d\xe8\x84\x5f\xbc\x63\x7c\x64\x10\x47\xf1\x1a\x37\xba\xb5\x0c\x46\x07\xbd\x22\x4b\xfc\x04\xd4\x04\x8b\x0d\x3a\xd6\x4c\xde\xa5\xb7\xc7\x2c\x46\xc7\x15\x59\xbc\xc3\xf8\x17\xe3\x7b\xdb\x9f\x84\xfb\x34\x32\xd4\xdb\xdd\x12\x1f\x5a\x4c\xbc\xc4\x14\xbc\x4b\x78\x42\xb5\xb8\x38\xb1\xa8\xfd\x3d\x3a\x38\xb1\xbc\xc3\x94\xc8\x3b\xd1\x09\xb1\x69\x9d\x01\x69\xa0\x98\x93\xca\xe1\xc6\x11\xcb\x3d\xf4\xb1\x57\x47\xfe\xaf\x48\x8c\xb1\x84\x08\xc5\x34\x1f\xae\x97\x3f\x07\x5a\xdf\x2e\x6e\xcf\x80\x1c\x31\x69\x4b\xff\x10\xcc\x94\x62\xa1\x7a\x90\x52\x70\x7d\x59\x83\x7a\x97\x78\x7d\x59\xff\xfe\x20\xd1\x54\xbd\xe9\xc6\xad\xf1\x51\xee\x4b\x88\xf9\x4b\x08\xf5\xe3\x93\xa4\xc9\xf4\x21\xdc\x31\xc1\x25\xba\x35\xc6\xaf\xf5\xf7\x6f\xd2\x94\x90\x0d\xa0\xac\x84\x7d\x09\x8e\x6c\x3e\x27\xda\x57\x59\xce\xe1\x28\x58\x79\xcf\x89\xa3\x0e\xa0\x43\xb0\x64\x86\xd2\xc0\xc6\xc7\x39\xa5\xd7\x75\x1d\x99\x21\x94\x80\x31\xc2\xd9\xf9\x73\x3d\x8f\x07\xc9\xcf\xb3\xf8\xd0\xe5\x83\x9e\x36\x83\xfa\xd3\x79\x7f\x97\x89\xdd\x3f\xd6\x6f\xab\x2b\xcd\xda\x5a\x27\x31\xc6\x51\xdc\x8d\xe5\x7e\xb9\xd3\x21\x54\xcb\xd6\xc9\xbc\xff\x8f\xff\x01\x00\x00\xff\xff\xec\x27\xaf\x24\x95\x03\x00\x00")

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

	info := bindataFileInfo{name: "templates/main.go.tpl", size: 917, mode: os.FileMode(420), modTime: time.Unix(1531444452, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x5d, 0x9a, 0x99, 0xcd, 0xd2, 0xd2, 0x1a, 0xb6, 0x80, 0x18, 0xbf, 0x98, 0x5c, 0x2a, 0x21, 0x68, 0xa, 0x98, 0x83, 0x62, 0xe9, 0x48, 0x5b, 0xce, 0x81, 0x6b, 0xe6, 0x72, 0xe6, 0x4a, 0x10, 0xa5}}
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
	"templates/.at-conf.json.tpl": templatesAtConfJsonTpl,

	"templates/index.tpl.tpl": templatesIndexTplTpl,

	"templates/layout.tpl.tpl": templatesLayoutTplTpl,

	"templates/main.go.tpl": templatesMainGoTpl,
}

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
	"templates": &bintree{nil, map[string]*bintree{
		".at-conf.json.tpl": &bintree{templatesAtConfJsonTpl, map[string]*bintree{}},
		"index.tpl.tpl":     &bintree{templatesIndexTplTpl, map[string]*bintree{}},
		"layout.tpl.tpl":    &bintree{templatesLayoutTplTpl, map[string]*bintree{}},
		"main.go.tpl":       &bintree{templatesMainGoTpl, map[string]*bintree{}},
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
