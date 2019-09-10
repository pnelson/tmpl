// Package tmpl implements a minimal view interface over html/template.
// Use of this package encourages proper separation of concerns.
package tmpl

import (
	"html/template"
	"io"
	"strings"
	"sync"
)

// Template represents a set of HTML templates.
type Template struct {
	mu        sync.Mutex
	loader    Loader
	loaded    map[string]string
	parsed    map[string]*template.Template
	recompile bool
}

// New returns a new template set.
func New(opts ...Option) *Template {
	t := &Template{
		loader:    defaultLoader,
		loaded:    make(map[string]string),
		parsed:    make(map[string]*template.Template),
		recompile: false,
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
	p, err := t.prepare(view)
	if err != nil {
		return err
	}
	return p.Execute(w, view)
}

// prepare returns the parsed templates representing the view.
func (t *Template) prepare(view Viewable) (*template.Template, error) {
	if view == nil {
		return template.New("nil"), nil
	}
	names := view.Templates()
	key := strings.Join(names, ":")
	t.mu.Lock()
	defer t.mu.Unlock()
	p, ok := t.parsed[key]
	if !ok || t.recompile {
		var err error
		p, err = t.parse(names)
		if err != nil {
			return nil, err
		}
		t.parsed[key] = p
	}
	return p, nil
}

// parse returns the parsed template.
func (t *Template) parse(names []string) (*template.Template, error) {
	var rv *template.Template
	for _, name := range names {
		var tmpl *template.Template
		if rv == nil {
			rv = template.New(name)
			tmpl = rv
		} else {
			tmpl = rv.New(name)
		}
		s, err := t.load(name)
		if err != nil {
			return nil, err
		}
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return rv, nil
}

// load returns the template contents retrieved from the loader.
func (t *Template) load(name string) (string, error) {
	s, ok := t.loaded[name]
	if !ok || t.recompile {
		b, err := t.loader.Load(name)
		if err != nil {
			return "", err
		}
		s = string(b)
		t.loaded[name] = s
	}
	return s, nil
}
