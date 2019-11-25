package golang

import (
	"bytes"
	"path/filepath"

	"github.com/phogolabs/stride/codegen"
)

// ServerGenerator builds a server
type ServerGenerator struct {
	Path        string
	Controllers codegen.ControllerDescriptorCollection
}

// Generate generates a file
func (g *ServerGenerator) Generate() *File {
	filename := filepath.Join(g.Path, "server.go")

	writer := &TemplateWriter{
		Path: "syntax/golang/server.go.tpl",
		Context: map[string]interface{}{
			"controllers": g.Controllers,
		},
	}

	buffer := &bytes.Buffer{}
	if _, err := writer.WriteTo(buffer); err != nil {
		panic(err)
	}

	root, err := ReadFile(filename, buffer)
	if err != nil {
		panic(err)
	}

	return root
}
