// Code generated by go-bindata.
// sources:
// tmpl/flat/footer.tmpl
// tmpl/flat/header.tmpl
// tmpl/flat/page1.tmpl
// tmpl/flat/page2and3.tmpl
// tmpl/inheritance/base.tmpl
// tmpl/inheritance/content1.tmpl
// tmpl/inheritance/content2.tmpl
// DO NOT EDIT!

package main

import (
	"bytes"
	"compress/gzip"
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
	bytes []byte
	info  os.FileInfo
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

var _tmplFlatFooterTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xaa\xae\x4e\x49\x4d\xcb\xcc\x4b\x55\x50\x4a\xcb\xcf\x2f\x49\x2d\x52\xaa\xad\xe5\xb2\x81\x30\xed\x6c\x92\xec\x72\x73\x6d\xf4\x93\xec\x92\x8a\xf2\x8b\x15\xd4\x92\xf3\x0b\x2a\xad\x15\x8c\x0c\x0c\xcd\x6c\xf4\xa1\x2a\xb8\x6c\xf4\x93\xf2\x53\x2a\x41\x74\x46\x49\x6e\x8e\x1d\x57\x75\x75\x6a\x5e\x4a\x6d\x2d\x17\x20\x00\x00\xff\xff\xc3\xe7\xc4\x59\x57\x00\x00\x00")

func tmplFlatFooterTmplBytes() ([]byte, error) {
	return bindataRead(
		_tmplFlatFooterTmpl,
		"tmpl/flat/footer.tmpl",
	)
}

func tmplFlatFooterTmpl() (*asset, error) {
	bytes, err := tmplFlatFooterTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "tmpl/flat/footer.tmpl", size: 87, mode: os.FileMode(511), modTime: time.Unix(1480198841, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _tmplFlatHeaderTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xaa\xae\x4e\x49\x4d\xcb\xcc\x4b\x55\x50\xca\x48\x4d\x4c\x49\x2d\x52\xaa\xad\xe5\xb2\xc9\x28\xc9\xcd\xb1\xe3\xb2\x01\x89\xd8\x71\x29\x28\x28\x28\xd8\x94\x64\x96\xe4\xa4\xda\x05\x24\xa6\xa7\x2a\x84\x80\x98\x36\xfa\x10\x11\x2e\x1b\x7d\x88\x2a\x9b\xa4\xfc\x94\x4a\x3b\xae\xea\xea\xd4\xbc\x94\xda\x5a\x40\x00\x00\x00\xff\xff\x43\xd0\xfc\x32\x56\x00\x00\x00")

func tmplFlatHeaderTmplBytes() ([]byte, error) {
	return bindataRead(
		_tmplFlatHeaderTmpl,
		"tmpl/flat/header.tmpl",
	)
}

func tmplFlatHeaderTmpl() (*asset, error) {
	bytes, err := tmplFlatHeaderTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "tmpl/flat/header.tmpl", size: 86, mode: os.FileMode(511), modTime: time.Unix(1433023486, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _tmplFlatPage1Tmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x54\xcc\x31\x0a\xc3\x30\x0c\x46\xe1\xdd\xa7\x10\xde\x4b\xf1\x6e\xb4\xf4\x02\xbd\x82\xc0\x7f\xe2\x80\x2b\x8b\xa0\x4e\xc2\x77\x2f\xa5\x50\xc8\xfa\x78\x7c\x11\x0d\xdb\xa1\xa0\x6c\xb2\xe3\x56\xf2\x5a\x29\xc2\xf1\xb2\x21\x0e\xca\x1d\xd2\x70\x7e\x6b\xed\x85\x9f\xb2\x1f\x2a\x54\xea\xbd\x17\x4e\xb5\x9f\x9c\xaa\xf1\x63\xaa\x43\xdf\x3e\xa9\x61\x0c\x21\xfb\x5f\xc6\x17\x6c\x9b\xd3\x7f\x58\x04\xb4\xad\xf5\x09\x00\x00\xff\xff\x52\x0f\x1b\xe5\x7e\x00\x00\x00")

func tmplFlatPage1TmplBytes() ([]byte, error) {
	return bindataRead(
		_tmplFlatPage1Tmpl,
		"tmpl/flat/page1.tmpl",
	)
}

func tmplFlatPage1Tmpl() (*asset, error) {
	bytes, err := tmplFlatPage1TmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "tmpl/flat/page1.tmpl", size: 126, mode: os.FileMode(511), modTime: time.Unix(1433023486, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _tmplFlatPage2and3Tmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\xcd\x41\x0a\xc2\x30\x10\x85\xe1\x7d\x4e\x11\xb2\x17\x69\xb3\x0d\xd9\x78\x01\xaf\x30\x90\xd7\xa6\x10\x27\x43\x19\x57\x43\xee\x2e\x22\x08\xba\x69\xb7\x8f\xc7\xff\x99\x15\x2c\x1b\xc3\x07\xa1\x15\x97\x39\x8c\xe1\xcc\x14\x0f\x69\xa4\xf0\xa1\x82\x0a\xf6\xf7\x9a\xea\x94\xef\xb4\x6e\x4c\x7e\x4e\xd7\x3a\x65\x97\xea\x9e\x5d\x92\x7c\xeb\xac\xe0\xa7\x76\x5f\xd0\x1a\x79\xf9\xbe\x24\xff\xc4\x96\xde\xf5\x13\x33\x03\x97\x31\x9c\xfb\xf3\xe3\x29\x3f\x9e\xf2\xe3\xb1\xff\x0a\x00\x00\xff\xff\x55\x64\xb0\xee\xfe\x00\x00\x00")

func tmplFlatPage2and3TmplBytes() ([]byte, error) {
	return bindataRead(
		_tmplFlatPage2and3Tmpl,
		"tmpl/flat/page2and3.tmpl",
	)
}

func tmplFlatPage2and3Tmpl() (*asset, error) {
	bytes, err := tmplFlatPage2and3TmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "tmpl/flat/page2and3.tmpl", size: 254, mode: os.FileMode(511), modTime: time.Unix(1433023486, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _tmplInheritanceBaseTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x4c\x90\x41\x4f\xf3\x30\x0c\x86\xef\xf9\x15\xfe\x72\xf8\x4e\x84\x30\x24\x24\x24\x9c\x48\x30\x76\xe0\x04\x87\x22\xc1\x31\x6d\xbd\xd5\x52\xd2\x94\xd6\xeb\x34\x4d\xfb\xef\xa8\xb4\x43\x9c\xa2\x3c\x7e\xfc\xda\x32\xfe\x7b\x7e\x5d\x17\x9f\x6f\x1b\x68\x24\x45\xaf\x70\x7a\x20\x86\x76\xe7\x34\x8b\x9e\x00\x85\xda\x2b\x4c\x24\x01\xaa\x26\xf4\x03\x89\xd3\x7b\xd9\x9a\x7b\x7d\xc1\x8d\x48\x67\xe8\x6b\xcf\xa3\xd3\x1f\xe6\xfd\xd1\xac\x73\xea\x82\x70\x19\x49\x43\x95\x5b\xa1\x56\x9c\x7e\xd9\x38\xaa\x77\xf4\xdb\xd5\x86\x44\x4e\x8f\x4c\x87\x2e\xf7\xf2\x47\x3c\x70\x2d\x8d\xab\x69\xe4\x8a\xcc\xcf\xe7\x0a\xb8\x65\xe1\x10\xcd\x50\x85\x48\x6e\x35\x85\x08\x4b\x24\xff\x14\x06\x82\x82\x52\x17\x83\x10\x14\x13\x43\x3b\x97\x14\xda\x65\xf9\x32\xd7\x47\xaf\x4e\x27\x90\x8b\xa8\x97\x69\x1a\xae\xe1\x7c\x56\xb8\xcd\x59\xa8\xf7\x58\xfa\x94\xd0\x96\xbe\xec\xf3\x00\xff\xab\xdc\x1d\x1f\xe0\xf6\x66\x75\x87\x76\x31\x14\xda\x39\x0d\xed\x7c\xb1\xef\x00\x00\x00\xff\xff\x8e\x7a\xb9\x82\x42\x01\x00\x00")

func tmplInheritanceBaseTmplBytes() ([]byte, error) {
	return bindataRead(
		_tmplInheritanceBaseTmpl,
		"tmpl/inheritance/base.tmpl",
	)
}

func tmplInheritanceBaseTmpl() (*asset, error) {
	bytes, err := tmplInheritanceBaseTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "tmpl/inheritance/base.tmpl", size: 322, mode: os.FileMode(511), modTime: time.Unix(1452196665, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _tmplInheritanceContent1Tmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xaa\xae\x4e\x49\x4d\xcb\xcc\x4b\x55\x50\x4a\xce\xcf\x2b\x49\xcd\x2b\x51\xaa\xad\xe5\x82\x32\x15\x0c\xb9\xaa\xab\x53\xf3\x52\x6a\x6b\x01\x01\x00\x00\xff\xff\x1a\xd3\xbd\xe7\x26\x00\x00\x00")

func tmplInheritanceContent1TmplBytes() ([]byte, error) {
	return bindataRead(
		_tmplInheritanceContent1Tmpl,
		"tmpl/inheritance/content1.tmpl",
	)
}

func tmplInheritanceContent1Tmpl() (*asset, error) {
	bytes, err := tmplInheritanceContent1TmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "tmpl/inheritance/content1.tmpl", size: 38, mode: os.FileMode(511), modTime: time.Unix(1433192246, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _tmplInheritanceContent2Tmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xaa\xae\x4e\x49\x4d\xcb\xcc\x4b\x55\x50\x4a\xce\xcf\x2b\x49\xcd\x2b\x51\xaa\xad\xe5\x82\x32\x15\x8c\xb8\xaa\xab\x53\xf3\x52\x6a\x6b\x01\x01\x00\x00\xff\xff\xdf\xef\x30\xde\x26\x00\x00\x00")

func tmplInheritanceContent2TmplBytes() ([]byte, error) {
	return bindataRead(
		_tmplInheritanceContent2Tmpl,
		"tmpl/inheritance/content2.tmpl",
	)
}

func tmplInheritanceContent2Tmpl() (*asset, error) {
	bytes, err := tmplInheritanceContent2TmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "tmpl/inheritance/content2.tmpl", size: 38, mode: os.FileMode(511), modTime: time.Unix(1433192251, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
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

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
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
	"tmpl/flat/footer.tmpl": tmplFlatFooterTmpl,
	"tmpl/flat/header.tmpl": tmplFlatHeaderTmpl,
	"tmpl/flat/page1.tmpl": tmplFlatPage1Tmpl,
	"tmpl/flat/page2and3.tmpl": tmplFlatPage2and3Tmpl,
	"tmpl/inheritance/base.tmpl": tmplInheritanceBaseTmpl,
	"tmpl/inheritance/content1.tmpl": tmplInheritanceContent1Tmpl,
	"tmpl/inheritance/content2.tmpl": tmplInheritanceContent2Tmpl,
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
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
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
	"tmpl": &bintree{nil, map[string]*bintree{
		"flat": &bintree{nil, map[string]*bintree{
			"footer.tmpl": &bintree{tmplFlatFooterTmpl, map[string]*bintree{}},
			"header.tmpl": &bintree{tmplFlatHeaderTmpl, map[string]*bintree{}},
			"page1.tmpl": &bintree{tmplFlatPage1Tmpl, map[string]*bintree{}},
			"page2and3.tmpl": &bintree{tmplFlatPage2and3Tmpl, map[string]*bintree{}},
		}},
		"inheritance": &bintree{nil, map[string]*bintree{
			"base.tmpl": &bintree{tmplInheritanceBaseTmpl, map[string]*bintree{}},
			"content1.tmpl": &bintree{tmplInheritanceContent1Tmpl, map[string]*bintree{}},
			"content2.tmpl": &bintree{tmplInheritanceContent2Tmpl, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
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
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
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
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
