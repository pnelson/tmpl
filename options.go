package tmpl

// Option represents a functional option for configuration.
type Option func(*Template)

// WithLoader sets the loader to retrieve templates from.
func WithLoader(loader Loader) Option {
	return func(t *Template) {
		t.loader = loader
	}
}

// WithPool sets the pool to to get and put byte buffers from.
func WithPool(pool Pool) Option {
	return func(t *Template) {
		t.pool = pool
	}
}

// WithRecompile sets the flag that indicates if templates are
// to be reompiled on demand. This may be useful for development.
// Defaults to false.
func WithRecompile(recompile bool) Option {
	return func(t *Template) {
		t.recompile = recompile
	}
}
