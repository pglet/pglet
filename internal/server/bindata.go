package server

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

var _client_build_index_html = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x54\x91\xb1\x6e\xf3\x30\x0c\x84\x77\x03\x7e\x87\xfb\xb3\x5b\x44\xd6\x1f\xaa\x97\xb6\x73\x33\x64\xe9\x28\x5b\x84\xac\x56\x96\x04\x99\x4e\x9a\xb7\x2f\x2c\xc7\x43\x17\x11\xe4\x91\xc7\x0f\xa2\xfe\xf7\xf6\xf1\x7a\xfd\xbc\xbc\x63\x92\x39\xf4\x6d\xa3\xb7\x88\x60\xa2\x7b\x39\x71\x3c\xd5\x0a\x1b\xbb\x45\xf1\x12\xb8\xbf\xb8\xc0\x02\xcb\x37\x0e\x29\xcf\x1c\x05\xd9\x38\xd6\xb4\xab\x6d\xa3\xe9\xe8\x1f\x92\x7d\xd4\xf9\x73\x9f\xb7\x21\x4d\xd3\x79\xcb\x73\x7f\x9d\xfc\x02\xbf\xfc\x71\xd1\x63\xb2\xdc\xfb\x68\xf9\x47\x6d\x0c\x9a\x6a\xa1\xba\xab\xb6\x01\x80\x6b\xc2\xc2\xe5\xc6\xc8\x25\xd9\x75\x14\x9f\x22\x4c\xce\xc8\xeb\x10\xfc\x32\xc1\x44\x0b\x9e\x07\xb6\xf0\x82\xbb\x97\xe9\xbf\xa6\x5c\x37\x96\x8d\x6c\xb4\x18\x83\xe7\x28\x6d\xf3\x30\x25\xee\x2f\x86\xd5\x07\x5b\x45\xa5\xda\xc6\x25\x38\x16\x74\x2b\x9c\x97\x69\x1d\xd4\x98\x66\xfa\x12\xe6\xf5\xce\x91\x5c\xea\x06\x1f\xad\x11\x43\x6a\xef\x3e\x72\x74\x09\x3e\x0a\x97\x68\x02\x55\xc8\x42\x4f\x49\xb9\x84\xae\xcb\xdf\x6e\x87\x2f\x4f\x08\xaa\x8b\x77\x1f\x4d\x3b\xa1\xa6\xe3\xcf\xa8\x9e\xe3\x37\x00\x00\xff\xff\xa5\xb7\xcb\xdf\x9e\x01\x00\x00")

func client_build_index_html() ([]byte, error) {
	return bindata_read(
		_client_build_index_html,
		"client/build/index.html",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
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
var _bindata = map[string]func() ([]byte, error){
	"client/build/index.html": client_build_index_html,
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
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"client/build/index.html": &_bintree_t{client_build_index_html, map[string]*_bintree_t{
	}},
}}
