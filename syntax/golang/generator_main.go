package golang

import (
	"bytes"
	"path/filepath"

	"github.com/phogolabs/stride/contract"
)

// MainGenerator builds the main
type MainGenerator struct {
	Path     string
	Reporter contract.Reporter
}

// Generate generates a file
func (g *MainGenerator) Generate() *File {
	var (
		command  = filepath.Base(g.Path)
		filename = filepath.Join(g.Path, "main.go")
	)

	reporter := g.Reporter.With(contract.SeverityHigh)
	reporter.Notice(" Generating main file: %s...", filename)

	writer := &TemplateWriter{
		Path: "syntax/golang/main.go.tpl",
		Context: map[string]interface{}{
			"command": command,
		},
	}

	buffer := &bytes.Buffer{}
	if _, err := writer.WriteTo(buffer); err != nil {
		reporter.Error(" Generating main file: %s fail: %v", filename, err)
		return nil
	}

	root, err := ReadFile(filename, buffer)
	if err != nil {
		reporter.Error(" Generating main file: %s fail: %v", filename, err)
		return nil
	}

	reporter.Notice(" Generating main file: %s successful", filename)
	return root
}
