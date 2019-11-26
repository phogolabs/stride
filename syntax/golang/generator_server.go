package golang

import (
	"bytes"
	"path/filepath"

	"github.com/phogolabs/stride/codedom"
	"github.com/phogolabs/stride/contract"
)

// ServerGenerator builds a server
type ServerGenerator struct {
	Path        string
	Controllers codedom.ControllerDescriptorCollection
	Reporter    contract.Reporter
}

// Generate generates a file
func (g *ServerGenerator) Generate() *File {
	filename := filepath.Join(g.Path, "server.go")

	reporter := g.Reporter.With(contract.SeverityHigh)
	reporter.Notice(" Generating server file: %s...", filename)

	writer := &TemplateWriter{
		Path: "syntax/golang/server.go.tpl",
		Context: map[string]interface{}{
			"controllers": g.Controllers,
		},
	}

	buffer := &bytes.Buffer{}
	if _, err := writer.WriteTo(buffer); err != nil {
		reporter.Error(" Generating server file: %s fail: %v", filename, err)
		return nil
	}

	root, err := ReadFile(filename, buffer)
	if err != nil {
		reporter.Error(" Generating server file: %s fail: %v", filename, err)
		return nil
	}

	reporter.Notice(" Generating server file: %s successful", filename)
	return root
}
