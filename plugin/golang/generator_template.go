package golang

import (
	"io"
	"io/ioutil"
	"text/template"

	"github.com/phogolabs/parcello"
)

// TemplateWriter executes the template
type TemplateWriter struct {
	Path    string
	Context map[string]interface{}
}

// WriteTo writes the executed template
func (g *TemplateWriter) WriteTo(w io.Writer) (int64, error) {
	file, err := parcello.Open(g.Path)
	if err != nil {
		return 0, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return 0, err
	}

	pattern, err := template.New("source").Parse(string(data))
	if err != nil {
		return 0, err
	}

	err = pattern.Execute(w, g.Context)
	return 0, err
}
