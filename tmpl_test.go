package tmpl_test

import (
	"bytes"
	"embed"
	"flag"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/pnelson/tmpl"
)

var update = flag.Bool("update", false, "update .golden files")

//go:embed testdata
var embedFS embed.FS

// extFS is a fs.FS implementation that appends a common extension.
type extFS struct {
	fs  fs.FS
	ext string
}

// Open implements the fs.FS interface.
func (f extFS) Open(name string) (fs.File, error) {
	return f.fs.Open(name + f.ext)
}

func TestGolden(t *testing.T) {
	var tt = map[string]tmpl.Viewable{
		"basic.html":  basic{Title: "test"},
		"layout.html": index{layout: layout{Title: "test"}},
		"nested.html": nested{layout: layout{Title: "test"}},
	}
	testdata, err := fs.Sub(embedFS, "testdata")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	templates := tmpl.New(extFS{testdata, ".html"})
	for filename, view := range tt {
		compare(t, templates, filename, view)
	}
}

func compare(t *testing.T, set *tmpl.Template, filename string, view tmpl.Viewable) {
	buf := bytes.NewBuffer(make([]byte, 0))
	err := set.Render(buf, view)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	have := buf.Bytes()
	golden := filepath.Join("testdata", "golden", filename)
	if *update {
		os.WriteFile(golden, have, 0644)
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(have, want) {
		t.Errorf("does not match golden file")
	}
}

type basic struct {
	Title string
}

func (basic) Templates() []string {
	return []string{"basic"}
}

type layout struct {
	Title string
}

type index struct {
	layout
}

func (index) Templates() []string {
	return []string{"layout", "index"}
}

type nested struct {
	layout
}

func (nested) Templates() []string {
	return []string{"layout", "nested", "list_item"}
}
