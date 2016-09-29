// Code generated by go-bindata.
// sources:
// events.js
// DO NOT EDIT!

package events

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

var _eventsJs = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xb4\x18\xcd\x6e\xdb\x38\xf3\xee\xa7\x60\x7c\xa8\x25\xc4\x55\xe3\xb6\x27\xab\xc6\x77\xe8\xd7\x02\x05\x8a\xed\x1e\x7a\x33\x8c\x40\x96\x28\x9b\xa9\x42\x19\x12\x95\x26\x48\xfc\xee\x3b\x43\x8a\x12\x49\x51\x8e\x37\x9b\xea\x60\x8b\xc3\xf9\xff\xe3\x50\x93\x69\x53\x53\x52\x8b\x8a\xa5\x62\x1a\x4f\xee\x92\x8a\xb0\xec\x73\xd9\x70\x41\x2b\xb2\x22\x57\xf1\x24\x6f\x78\x2a\x58\xc9\xc9\x8e\x8a\x6f\xff\x0f\x42\xf2\x38\x21\xf0\x54\x54\x34\x15\x27\xd3\x29\xb9\x24\xc1\xe5\x65\x47\x14\xc6\x93\x63\x4f\x93\x26\x45\xf1\x15\x16\x41\xce\xe7\x24\xa9\x76\xb5\x26\x67\x39\x09\x70\x4d\x56\xab\x15\xb9\x2b\x59\x46\xae\x60\x8b\x28\x10\x59\x6f\x62\x72\x94\x78\xa8\x50\x01\x90\x9c\x47\x05\xe5\x3b\xb1\x9f\x13\x06\xcb\xb7\x0b\x60\xb7\x80\x17\x24\x58\x5f\x6d\x60\xf5\x5e\xaf\x16\xb8\xfa\xa0\x57\xef\x71\xf5\x51\xaf\x3e\x6c\x62\xc9\xb6\xfe\xcd\x44\xba\x57\x2a\xb4\x8c\xb5\x66\xf8\xa4\x09\x38\xe5\x6a\xd9\xad\xf1\xf9\xbd\x67\x05\x95\xa6\x92\x4f\xa4\x08\xad\x3d\x7c\x72\xbe\x66\x9b\x68\x9f\xf0\xac\xa0\x55\x84\x86\x07\x0a\x94\x8a\xfb\x30\xb6\xd0\x95\xef\x62\x5b\xdc\xe2\xf5\xc4\xa1\x6f\xce\x11\xf9\xfe\x75\x45\x62\x10\xce\x11\xfb\xe1\xf5\xc5\x62\xc4\xcf\x11\xfd\xf1\xcf\x88\xc6\x14\x7b\x46\x7c\x46\xf3\xa4\x29\xc4\x7f\x91\x9f\x1c\x0e\xc5\x83\xa5\x00\x16\xd4\xb8\xd8\x23\x54\x22\xbd\x3f\x94\x95\xa8\x23\x5d\x88\x50\x08\xfa\xd5\x28\x6d\x56\x7f\x6d\x5f\x83\xc4\x29\x70\xf1\x70\xa0\x65\x4e\x12\x59\xa8\x33\x4d\x31\x8b\x0d\xde\x3d\x35\x70\xef\x17\x16\xff\x2f\x77\x94\x8b\x2f\xb7\x4c\x40\x8b\x18\xc8\x48\x08\xe3\xb5\x48\x78\x8a\x92\x4c\x4c\xf2\xf4\x44\x02\x53\xb9\xa8\xe4\x21\x79\xf3\x86\xd8\xb0\x3c\x1f\x02\xa1\xa3\xed\x76\xd0\x8e\x42\x5b\x53\x8b\xfb\xca\x51\x4c\xf5\x3f\x07\x25\xe8\xac\xe8\x7a\x5f\x07\xb1\xac\xd2\xbb\xaa\x71\xfd\xd8\xde\xd0\x54\x44\x10\x76\xc6\xe9\xdf\x55\x79\xa0\x95\x78\x08\x4c\x82\xe8\x50\x95\xa2\x44\xff\xce\xc9\xb4\x60\xb5\xa0\x9c\x56\xf5\x74\x6e\xf4\x21\x68\xb9\x4b\x32\x94\x6f\x07\x9b\x88\x3d\xab\xa3\xeb\x8e\x43\x9f\x10\xc7\x79\xf7\x4a\x79\x73\x4b\xab\x64\x5b\xd0\x25\x11\x55\x43\xfb\x9d\xb4\xe4\x39\xdb\x35\xc6\x9e\xb2\xa2\x4d\x2c\xbf\xc6\x91\x8c\x75\xaf\x19\x45\xac\x39\xc1\x2e\x2f\x33\xb3\x84\x58\x9a\xda\x62\xbb\x47\x98\xd3\xee\x15\x88\xe4\x49\x51\x53\xdd\xf1\xf1\xc1\x30\x48\x96\x78\x18\x04\x8e\x7d\x32\x29\x5c\xd8\x8a\x3c\x1e\xc3\x70\x2d\x89\x36\x3e\x0c\xbd\x85\x87\x8b\x51\x33\x4a\x4a\x74\x68\xea\x7d\x60\x3b\x97\x27\xb7\xe0\x0f\x65\x98\xb5\x81\x4a\x2f\xe5\xaf\x0d\x6f\xcb\x74\x89\x5e\xb0\x36\xc0\x23\x4b\xfc\x41\xb5\x50\xab\x3e\x3e\x86\x22\x46\x2c\xdb\xea\x7d\xc6\xff\xca\x73\x23\x11\x30\x7d\x6f\x66\x09\x94\x86\x1b\x2a\x8c\x78\x78\x96\xc8\x3c\x1f\x4a\xfc\x0b\xbc\x84\xcc\x4c\x81\xc3\xd8\x78\x22\xf8\x78\x8c\xad\xec\xe8\xb8\x41\x8a\x10\xde\x14\x85\x9b\xec\xbe\x88\x1b\xa9\xde\x47\x14\x72\x49\x72\xf4\x26\x00\x4a\xd8\xb8\xac\xad\x74\x1b\x27\xb3\x3b\x2d\x8a\xc8\xf9\x98\xb6\x1e\x8d\x0d\x4e\x6a\xc4\xb1\x08\x8e\xd6\x4a\x1a\x31\xe4\x98\x97\x15\x09\xe4\x8c\x26\x67\x33\x82\xc7\x46\x9b\xc1\x6a\x88\x01\xd0\xe5\xa5\x4f\x17\xcb\xcc\x6b\x9c\x9c\x14\x1d\x9c\x25\xb1\x17\xb9\x8b\x49\x6d\x1c\x41\x68\xad\x1d\xeb\xf3\x6d\x8e\xea\x43\xc1\x52\x1a\xb0\x39\x71\x87\x13\xbf\x0f\x86\x90\xa3\x13\xee\xd3\xf9\xda\x9e\x00\xfe\x9c\x35\x6d\x40\xbf\xf4\x83\xe7\x64\xe0\xec\x6b\xf4\xf6\x22\xc6\xff\x4f\x88\x08\x9d\xd4\xf4\xf8\xb5\xc7\xe5\x72\xd8\x04\xfc\xb7\x64\xb1\x51\xb3\xa7\x22\x02\xd8\xc6\x97\xb4\xef\xde\x15\x54\x8c\xb7\xbc\xe7\x3a\x9e\x74\xf0\xd3\xd3\x89\x94\x5f\x41\xd7\x33\xa4\x45\xd0\xf5\xd3\x44\x0c\x08\x66\x30\x1b\xcc\x80\x93\x89\xfd\xaf\xcb\xf9\x64\xf7\x36\xab\x00\xe8\x40\xd0\x69\x5d\x5a\x24\xbb\x59\x58\x11\xcf\xe8\xb6\xd9\x7d\x06\xe4\x6d\x92\xfe\xb2\x67\xa8\x71\x3c\x25\x0d\x24\xc3\xd5\xa7\x49\x45\x59\x45\x5c\x76\x32\x09\x56\xaf\x46\x7f\xc3\x70\xb6\x00\x73\xe0\xea\x0c\x05\x8c\x39\x81\x7c\xe8\xca\xaa\xbf\xac\xf0\x8c\xde\xdb\x14\x38\x7f\x8d\xe6\x9a\x2c\x6c\x20\x63\x2e\x2b\x55\xec\x8c\x79\x2b\x5c\x62\x8e\xd7\x74\xd2\x5e\x7f\x86\x0d\x4c\x12\x48\x7b\xb1\xb4\xa5\xc7\x7d\xd5\x8d\x0c\xcc\x52\x6e\x23\xe6\x99\x3f\xf1\xe9\x2e\x7d\xed\x99\x0b\xde\x09\x5f\xd0\xec\xa4\x9f\xd4\xc9\x2c\xf9\x9c\xe4\xd1\x1b\xd3\x8d\x19\xf2\x60\xf3\xf0\x95\x21\x19\x6b\xf4\xd2\x17\x9b\x48\xe2\xfc\xc8\xbd\x82\xf1\x39\x45\xab\xfb\x1c\xb2\x18\xf6\x3a\xb7\x87\x69\xe5\x95\xb5\xed\x55\xd4\x73\xf0\xd1\x7b\x9a\x36\x82\x7e\x6f\x45\x2a\xfc\xc1\x15\xc0\x37\x45\x9c\xea\x91\x14\x20\xe0\x8a\xd3\x1d\xf4\x34\x0b\x57\x33\xab\xdf\xe2\x9b\x7d\xf1\xc7\x47\xd6\x8d\xa4\x2a\x2b\xeb\x42\x62\x3a\x64\x50\xa0\x8e\x1c\x3d\xeb\x0f\x6a\xa1\x67\x7c\x2e\x0b\xef\x18\xd1\xb2\x31\x4d\x38\x67\x4a\x52\x39\xf1\xb3\xb4\xdc\x50\x6e\x6f\xe6\xe4\x0f\xcc\xc8\x4a\x18\xe3\xbb\x9f\x25\xf4\x8c\x6c\x4e\x6e\x29\x74\x8a\x0e\xc5\xd8\x76\xf2\x5d\xc1\xdc\x63\x45\x63\x3e\x9a\x83\x29\xe8\xb1\x22\x60\x40\x6b\xd9\xb7\x4c\x92\x59\x80\x95\xfe\x36\x14\x7a\x65\xaf\x59\xb6\x51\x3c\xfa\x6d\x54\x14\x61\x68\xdc\xff\xc8\x0c\xff\x67\x64\x89\x2f\xb3\x1e\x09\x28\xd6\x88\xb8\x31\x87\x56\x54\xf8\xe5\x73\xb3\x0e\xcf\x0f\x77\x7e\xf6\x84\x68\x6c\x88\xd6\x3c\xfc\x61\x3d\x7b\x9e\xae\x45\x79\xf8\xae\xbd\x34\xaa\x4b\xaa\x0f\x36\xa7\x7e\x9e\x89\xad\x5d\x4b\x17\xc6\x4e\x38\x76\x85\xb4\x4f\xa9\x8a\xde\x96\x77\xe8\xa2\x0b\x75\xb8\xc0\xfd\xfa\x42\xeb\xe2\x30\xd7\x60\xc4\x69\x3f\x17\xb4\x07\x12\x9e\x29\xa5\xbc\x08\xcf\x6c\xb1\x1d\xc9\xca\x11\x2d\x8b\x60\x7b\x63\x63\x07\x83\xf4\x5c\x9b\xf9\x37\x48\xae\xfe\x34\xcd\xa0\xdb\x9b\xae\x72\x9b\x05\x90\x01\xb1\x93\xab\xb1\x8b\x82\x57\x9f\xc0\x89\xc7\x20\x11\xb5\xf6\xad\xe3\xa0\x48\x2e\xda\x6f\x00\xbf\xe8\x43\x6d\x14\x0c\x9c\x17\xa1\xb7\xd1\xe3\x93\x51\x18\x04\xe9\x30\x9e\xb6\x5e\xc7\x97\xa6\x7f\x46\xa1\x17\x96\x0f\x56\xb6\x0d\xae\x6e\x56\x62\x06\x86\x89\xea\x0a\x09\xbe\xb0\xf3\xbb\xd5\xc1\xfe\x92\x72\x94\xdd\x40\x7f\x78\x71\xbe\xa9\xd8\xa8\xff\x04\x00\x00\xff\xff\x5b\x78\x9e\x59\x84\x16\x00\x00")

func eventsJsBytes() ([]byte, error) {
	return bindataRead(
		_eventsJs,
		"events.js",
	)
}

func eventsJs() (*asset, error) {
	bytes, err := eventsJsBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "events.js", size: 5764, mode: os.FileMode(420), modTime: time.Unix(1475067918, 0)}
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
	"events.js": eventsJs,
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
	"events.js": &bintree{eventsJs, map[string]*bintree{}},
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

