package tmpl_test

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/pnelson/tmpl"
)

var update = flag.Bool("update", false, "update .golden files")

func TestGolden(t *testing.T) {
	var tt = map[string]tmpl.Viewable{
		"basic.html":  basic{Title: "test"},
		"layout.html": index{layout: layout{Title: "test"}},
		"nested.html": nested{layout: layout{Title: "test"}},
	}
	opts := tmpl.WithLoader(tmpl.NewFileSystemLoader("testdata", ".html"))
	templates := tmpl.New(opts)
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
		ioutil.WriteFile(golden, have, 0644)
	}
	want, err := ioutil.ReadFile(golden)
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
