// Package tmpl implements a minimal view interface over html/template.
// Use of this package encourages proper separation of concerns.
package tmpl

import (
	"html/template"
	"io"
	"reflect"
	"sync"
)

// Template represents a set of HTML templates.
type Template struct {
	mu        sync.Mutex
	loader    Loader
	recompile bool
	templates map[string]*template.Template
}

// New returns a new template set.
func New(opts ...Option) *Template {
	t := &Template{
		loader:    defaultLoader,
		recompile: false,
		templates: make(map[string]*template.Template),
	}
	for _, option := range opts {
		option(t)
	}
	return t
}

// Viewable represents a view.
type Viewable interface {
	// Templates returns a slice of template names to load and parse.
	Templates() []string
}

// Render writes the result of applying the templates
// associated with view to the view itself.
func (t *Template) Render(w io.Writer, view Viewable) error {
	p, err := t.load(view)
	if err != nil {
		return err
	}
	return p.Execute(w, view)
}

// load returns the parsed templates representing the view.
func (t *Template) load(view Viewable) (*template.Template, error) {
	if view == nil {
		return template.New("nil"), nil
	}
	names := view.Templates()
	if t.recompile {
		return t.parse(names)
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	key := reflect.TypeOf(view).String()
	p, ok := t.templates[key]
	if !ok {
		var err error
		p, err = t.parse(names)
		if err != nil {
			return nil, err
		}
		t.templates[key] = p
	}
	return p, nil
}

// parse returns the parsed template.
func (t *Template) parse(names []string) (*template.Template, error) {
	var rv *template.Template
	for _, name := range names {
		b, err := t.loader.Load(name)
		if err != nil {
			return nil, err
		}
		var tmpl *template.Template
		if rv == nil {
			rv = template.New(name)
			tmpl = rv
		} else {
			tmpl = rv.New(name)
		}
		_, err = tmpl.Parse(string(b))
		if err != nil {
			return nil, err
		}
	}
	return rv, nil
}
