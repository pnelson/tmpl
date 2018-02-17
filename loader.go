package tmpl

import (
	"io/ioutil"
	"path/filepath"
)

// Loader represents a template loader.
type Loader interface {
	// Load returns the contents of the template named by name.
	Load(name string) ([]byte, error)
}

// FileSystemLoader is a file system loader implementation.
type FileSystemLoader struct {
	root      string
	extension string
}

// defaultLoader represents the default template loader.
var defaultLoader = NewFileSystemLoader("./templates", ".html")

// NewFileSystemLoader returns a new file system loader.
func NewFileSystemLoader(root, extension string) *FileSystemLoader {
	return &FileSystemLoader{root: root, extension: extension}
}

// Load returns the contents of the template named by name.
// Names are expected to be file names relative to the root
// template directory and omit the extension.
//
// Load implements the Loader interface.
func (l *FileSystemLoader) Load(name string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(l.root, name+l.extension))
}
