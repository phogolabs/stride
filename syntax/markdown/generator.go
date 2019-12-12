package markdown

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"

	"github.com/phogolabs/stride/codedom"
	"github.com/phogolabs/stride/contract"
	"github.com/phogolabs/stride/inflect"
	"github.com/phogolabs/stride/syntax"
)

// Generator builds the main
type Generator struct {
	Path     string
	Reporter contract.Reporter
}

// Generate generates the source code
func (g *Generator) Generate(spec *codedom.SpecDescriptor) error {
	reporter := g.Reporter.With(contract.SeverityVeryHigh)

	if spec.Info == nil {
		return nil
	}

	reporter.Notice(" Generating markdown documentation...")

	project, err := filepath.Rel(filepath.Join(build.Default.GOPATH, "src"), g.Path)

	if err != nil {
		reporter.Error(" Generating markdown documentation fail: ", err)
		return err
	}

	ctx := map[string]interface{}{
		"command":     filepath.Base(g.Path),
		"project":     project,
		"title":       strings.TrimSpace(spec.Info.Title),
		"description": strings.TrimSpace(spec.Info.Description),
		"version":     strings.TrimSpace(spec.Info.Version),
	}

	// generate README.md
	if err := g.sync(filepath.Join(g.Path, "README.md"), ctx); err != nil {
		reporter.Error(" Generating markdown documentation fail: ", err)
		return err
	}

	reporter.Success(" Generating markdown documentation successful!")
	return nil
}

func (g *Generator) sync(path string, ctx map[string]interface{}) error {
	reporter := g.Reporter.With(contract.SeverityHigh)

	var (
		name   = inflect.LowerCase(filepath.Base(path))
		writer = &syntax.TemplateWriter{
			Path:    fmt.Sprintf("syntax/markdown/%s.tpl", name),
			Context: ctx,
		}
	)

	file, err := os.Create(path)
	if err != nil {
		reporter.Error(" Generating markdown file: %s fail: %v", path, err)
		return err
	}

	defer file.Close()

	if _, err := writer.WriteTo(file); err != nil {
		reporter.Error(" Generating markdown file: %s fail: %v", path, err)
		return err
	}

	reporter.Notice(" Generating markdown file: %s successful", path)
	return nil

}
