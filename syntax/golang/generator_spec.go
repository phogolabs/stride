package golang

import (
	"bytes"
	"path/filepath"

	"github.com/phogolabs/stride/contract"
	"github.com/phogolabs/stride/syntax"
)

// SpecGenerator builds the main
type SpecGenerator struct {
	Path     string
	Reporter contract.Reporter
}

// Generate generates a file
func (g *SpecGenerator) Generate() *File {
	filename := filepath.Join(g.Path, "service", "suite_test.go")

	reporter := g.Reporter.With(contract.SeverityHigh)
	reporter.Notice(" Generating spec suite file: %s...", filename)

	writer := &syntax.TemplateWriter{
		Path:    "syntax/golang/spec_suite.go.tpl",
		Context: map[string]interface{}{},
	}

	buffer := &bytes.Buffer{}
	if _, err := writer.WriteTo(buffer); err != nil {
		reporter.Error(" Generating spec suite file: %s fail: %v", filename, err)
		return nil
	}

	root, err := ReadFile(filename, buffer)
	if err != nil {
		reporter.Error(" Generating spec suite file: %s fail: %v", filename, err)
		return nil
	}

	reporter.Notice(" Generating spec suite file: %s successful", filename)
	return root
}
