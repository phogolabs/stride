package golang

import (
	"bytes"
	"path/filepath"
)

// MainGenerator builds the main
type MainGenerator struct {
	Path string
}

// Generate generates a file
func (g *MainGenerator) Generate() *File {
	var (
		command  = filepath.Base(g.Path)
		filename = filepath.Join(g.Path, "main.go")
	)

	writer := &TemplateWriter{
		Path: "syntax/golang/main.go.tpl",
		Context: map[string]interface{}{
			"command": command,
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
