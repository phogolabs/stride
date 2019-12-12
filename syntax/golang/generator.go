package golang

import (
	"os"
	"path/filepath"

	"github.com/phogolabs/stride/codedom"
	"github.com/phogolabs/stride/contract"
)

// FileGenerator is a file generator
type FileGenerator interface {
	// Generate generates the file
	Generate() *File
}

// Generator generates the source code
type Generator struct {
	Path     string
	Reporter contract.Reporter
}

// Generate generates the source code
func (g *Generator) Generate(spec *codedom.SpecDescriptor) error {
	var (
		dirPkg    = filepath.Join(g.Path, "service")
		dirCmd    = filepath.Join(g.Path, "cmd", filepath.Base(g.Path))
		generator FileGenerator
	)

	reporter := g.Reporter.With(contract.SeverityVeryHigh)
	reporter.Notice(" Generating spec...")

	generator = &SchemaGenerator{
		Path:       filepath.Join(g.Path, "service"),
		Collection: spec.Types,
		Reporter:   g.Reporter,
	}

	// writes the schema
	if err := g.sync(generator); err != nil {
		reporter.Error(" Generating spec fail")
		return err
	}

	// write the controller's schema
	for _, descriptor := range spec.Controllers {
		generator = &ControllerGenerator{
			Mode:       ControllerGeneratorModeSchema,
			Path:       dirPkg,
			Controller: descriptor,
			Reporter:   g.Reporter,
		}

		if err := g.sync(generator); err != nil {
			g.Reporter.Error(" Generating spec fail")
			return err
		}
	}

	// write the controller's api
	for _, descriptor := range spec.Controllers {
		generator = &ControllerGenerator{
			Mode:       ControllerGeneratorModeAPI,
			Path:       dirPkg,
			Controller: descriptor,
			Reporter:   g.Reporter,
		}

		if err := g.sync(generator); err != nil {
			reporter.Error(" Generating spec fail")
			return err
		}
	}

	// write the controller's spec
	for _, descriptor := range spec.Controllers {
		generator = &ControllerGenerator{
			Mode:       ControllerGeneratorModeSpec,
			Path:       dirPkg,
			Controller: descriptor,
			Reporter:   g.Reporter,
		}

		if err := g.sync(generator); err != nil {
			reporter.Error(" Generating spec fail")
			return err
		}
	}

	// write the server
	generator = &ServerGenerator{
		Path:        dirPkg,
		Controllers: spec.Controllers,
		Reporter:    reporter,
	}

	if err := g.sync(generator); err != nil {
		reporter.Error(" Generating spec fail")
		return err
	}

	// write the application main
	generator = &MainGenerator{
		Path:     dirCmd,
		Reporter: g.Reporter,
	}

	if err := g.sync(generator); err != nil {
		reporter.Error(" Generating spec fail")
		return err
	}

	markdown := &MarkdownGenerator{
		Path:     g.Path,
		Reporter: g.Reporter,
		Info:     spec.Info,
	}

	if err := markdown.Generate(); err != nil {
		reporter.Error(" Generating spec fail")
		return err
	}

	reporter.Success(" Generating spec complete!")
	return nil
}

func (g *Generator) sync(generator FileGenerator) error {
	reporter := g.Reporter.With(contract.SeverityLow)

	if target := generator.Generate(); target != nil {
		// merge if the file exist

		if source, err := OpenFile(target.Name()); err == nil {
			reporter.Info(" Merging file: %s...", target.Name())

			if err := target.Merge(source); err != nil {
				reporter.Error(" Merging file: %s fail: %v", target.Name(), err)
				return err
			}

			reporter.Success(" Merging file: %s successful", target.Name())
		}

		reporter.Info(" Sync file: %s...", target.Name())

		// mkdir create the directory
		dir := filepath.Dir(target.Name())

		// prepare the service package directory
		if err := os.MkdirAll(dir, 0755); err != nil {
			reporter.Error(" Sync file: %s fail: %v", target.Name(), err)
			return err
		}

		// write the file
		if err := target.Sync(); err != nil {
			reporter.Error("Sync file: %s fail: %v", target.Name(), err)
			return err
		}

		reporter.Success(" Sync file: %s successful", target.Name())
	}

	return nil
}
