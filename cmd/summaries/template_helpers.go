package summaries

import (
	"html/template"
	"io"

	"github.com/pkg/errors"
)

func renderTemplateToWriter(w io.Writer, tpl string, ctx interface{}) error {
	customFuncMap := template.FuncMap{}

	t := template.Must(template.New("t1").Funcs(customFuncMap).Parse(tpl))
	if err := t.Execute(w, ctx); err != nil {
		return errors.Wrap(err, "executing template")
	}
	return nil
}
