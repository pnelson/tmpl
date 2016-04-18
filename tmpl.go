// Package tmpl implements a minimal view interface over html/template.
// Use of this package encourages proper separation of concerns.
package tmpl

import (
	"html/template"
	"path/filepath"
	"reflect"
	"sync"
)

// Template represents a set of HTML templates.
type Template struct {
	sync.Mutex
	root      string
	extension string
	recompile bool
	pool      Pooler
	templates map[string]*template.Template
}

// New returns a new template set.
func New(root string, opts ...Option) *Template {
	t := &Template{
		root:      root,
		extension: defaultExtension,
		recompile: false,
		pool:      defaultPool,
		templates: make(map[string]*template.Template),
	}
	for _, option := range opts {
		option(t)
	}
	return t
}

// Viewable represents a view.
type Viewable interface {
	// Templates returns a slice of template names to parse.
	// The provided file names are expected to be relative to the
	// root template directory and omit the extension.
	Templates() []string
}

// Render returns the result of applying the templates
// associated with view to the view itself.
func (t *Template) Render(view Viewable) ([]byte, error) {
	p := t.load(view)
	b := t.pool.Get()
	defer t.pool.Put(b)
	err := p.Execute(b, view)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// load returns the parsed templates representing the view.
func (t *Template) load(view Viewable) *template.Template {
	names := view.Templates()
	if t.recompile {
		return t.parse(names)
	}
	t.Lock()
	defer t.Unlock()
	key := reflect.TypeOf(view).String()
	p, ok := t.templates[key]
	if !ok {
		p = t.parse(names)
		t.templates[key] = p
	}
	return p
}

// parse returns the parsed template.
func (t *Template) parse(names []string) *template.Template {
	return template.Must(template.ParseFiles(t.filenames(names)...))
}

// filenames returns the filenames for the template names.
func (t *Template) filenames(names []string) []string {
	rv := make([]string, len(names))
	for i, name := range names {
		rv[i] = t.filename(name)
	}
	return rv
}

// filename returns a fully qualified template filename.
func (t *Template) filename(name string) string {
	return filepath.Join(t.root, name+t.extension)
}
