package codegen

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-openapi/inflect"
	"golang.org/x/tools/imports"
)

// Renderer renders the dom
type Renderer interface {
	Render(w io.Writer)
}

// Generator generates the source code
type Generator struct{}

// Generate generates the source code
func (g *Generator) Generate(spec *SpecDescriptor) error {
	if err := g.render("contract.go", spec.Types); err != nil {
		return err
	}

	return nil
}

func (g *Generator) render(name string, renderer Renderer) error {
	buffer := g.buffer()

	// render the types
	renderer.Render(buffer)

	// format the source
	if err := g.format(buffer); err != nil {
		return err
	}

	output, err := os.Create(name)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, buffer)
	return err
}

func (g *Generator) buffer() *bytes.Buffer {
	buffer := &bytes.Buffer{}

	fmt.Fprintf(buffer, "package service")
	fmt.Fprintln(buffer)

	return buffer
}

func (g *Generator) format(buffer *bytes.Buffer) error {
	data, err := imports.Process("source.go", buffer.Bytes(), nil)
	if err != nil {
		return err
	}

	buffer.Reset()
	_, err = buffer.Write(data)
	return err
}

func camelize(text string) string {
	field := inflect.Camelize(text)

	switch {
	case field == "Id":
		field = strings.ToUpper(field)
	case strings.HasSuffix(field, "Id"):
		field = fmt.Sprintf("%vID", strings.TrimSuffix(field, "Id"))
	}

	return field
}
