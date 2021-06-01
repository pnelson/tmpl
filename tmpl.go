// Package tmpl implements a minimal view interface over html/template.
package tmpl

import (
	"html/template"
	"io"
	"io/fs"
	"strings"
	"sync"
)

// Template represents a set of HTML templates.
type Template struct {
	mu  sync.Mutex
	fs  fs.FS
	set map[string]*template.Template
}

// Viewable represents a view.
type Viewable interface {
	// Templates returns a slice of template names to parse.
	Templates() []string
}

// New returns a new template set.
func New(fs fs.FS) *Template {
	return &Template{
		fs:  fs,
		set: make(map[string]*template.Template),
	}
}

// Render writes the result of applying the templates
// associated with view to the view itself.
func (t *Template) Render(w io.Writer, view Viewable) error {
	tmpl, err := t.parse(view)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, view)
}

// parse returns the parsed templates representing the view.
func (t *Template) parse(view Viewable) (*template.Template, error) {
	if view == nil {
		return template.New("nil"), nil
	}
	names := view.Templates()
	key := strings.Join(names, ":")
	t.mu.Lock()
	defer t.mu.Unlock()
	tmpl, ok := t.set[key]
	if ok {
		return tmpl, nil
	}
	tmpl, err := template.ParseFS(t.fs, names...)
	if err != nil {
		return nil, err
	}
	t.set[key] = tmpl
	return tmpl, nil
}
