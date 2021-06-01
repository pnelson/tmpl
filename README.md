# tmpl

Package tmpl implements a minimal view interface over html/template.

Use of this package encourages proper separation of concerns. Templates
are expected to be "dumb" using only the base `text/template` actions
and leverage methods on the view for more complex output.
