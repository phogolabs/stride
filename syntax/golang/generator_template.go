package golang

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/phogolabs/parcello"
	"github.com/phogolabs/stride/inflect"
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

	m := template.FuncMap{
		"camelize":  inflect.Camelize,
		"dasherize": inflect.Dasherize,
		"uppercase": inflect.UpperCase,
		"titleize":  inflect.Titleize,
		"comment":   g.comment,
		"key":       g.key,
	}

	pattern, err := template.New("source").Funcs(m).Parse(string(data))
	if err != nil {
		return 0, err
	}

	err = pattern.Execute(w, g.Context)
	return 0, err
}

func (g *TemplateWriter) comment(parts ...string) string {
	text := strings.Join(parts, " ")
	text = strings.TrimSpace(text)

	if text == "" {
		return ""
	}

	return fmt.Sprintf("\n// %s", text)
}

func (g *TemplateWriter) key(parts ...string) string {
	buffer := &bytes.Buffer{}

	for _, part := range parts {
		if buffer.Len() > 0 {
			fmt.Fprint(buffer, ":")
		}

		fmt.Fprint(buffer, inflect.Dasherize(part))
	}

	return buffer.String()
}
