package tmpl

// Option represents a functional option for configuration.
type Option func(*Template)

// defaultExtension is the default template file name extension.
const defaultExtension = ".html"

// Extension sets the extension to use for finding templates.
func Extension(extension string) Option {
	return func(t *Template) {
		t.extension = extension
	}
}

// Pool sets the pool to to get and put byte buffers from.
func Pool(pool Pooler) Option {
	return func(t *Template) {
		t.pool = pool
	}
}

// Recompile sets the flag that indicates if templates are to be
// reompiled on demand. This may be useful for development.
// Defaults to false.
func Recompile(recompile bool) Option {
	return func(t *Template) {
		t.recompile = recompile
	}
}
