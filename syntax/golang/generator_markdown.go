package golang

import (
	"go/build"
	"os"
	"path/filepath"

	"github.com/phogolabs/stride/codedom"
	"github.com/phogolabs/stride/contract"
)

// MarkdownGenerator builds the main
type MarkdownGenerator struct {
	Path     string
	Reporter contract.Reporter
	Info     *codedom.InfoDescriptor
}

// Generate generates a file
func (g *MarkdownGenerator) Generate() error {
	var (
		command  = filepath.Base(g.Path)
		filename = filepath.Join(g.Path, "README.md")
	)

	reporter := g.Reporter.With(contract.SeverityHigh)
	reporter.Notice(" Generating markdown file: %s...", filename)

	project, err := filepath.Rel(filepath.Join(build.Default.GOPATH, "src"), g.Path)

	if err != nil {
		reporter.Error(" Generating markdown file: %s fail: %v", filename, err)
		return err
	}

	writer := &TemplateWriter{
		Path: "syntax/golang/readme.md.tpl",
		Context: map[string]interface{}{
			"command":     command,
			"project":     project,
			"title":       g.Info.Title,
			"description": g.Info.Description,
			"version":     g.Info.Version,
		},
	}

	file, err := os.Create(filename)
	if err != nil {
		reporter.Error(" Generating markdown file: %s fail: %v", filename, err)
		return err
	}

	defer file.Close()

	if _, err := writer.WriteTo(file); err != nil {
		reporter.Error(" Generating markdown file: %s fail: %v", filename, err)
		return err
	}

	reporter.Notice(" Generating markdown file: %s successful", filename)
	return nil
}
