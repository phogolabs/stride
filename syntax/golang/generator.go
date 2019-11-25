package golang

import (
	"os"
	"path/filepath"

	"github.com/phogolabs/stride/codedom"
)

// FileGenerator is a file generator
type FileGenerator interface {
	// Generate generates the file
	Generate() *File
}

// Generator generates the source code
type Generator struct {
	Path string
}

// Generate generates the source code
func (g *Generator) Generate(spec *codedom.SpecDescriptor) error {
	var (
		dirPkg    = filepath.Join(g.Path, "service")
		dirCmd    = filepath.Join(g.Path, "cmd", filepath.Base(g.Path))
		generator FileGenerator
	)

	generator = &SchemaGenerator{
		Path:       filepath.Join(g.Path, "service"),
		Collection: spec.Types,
	}

	// writes the schema
	if err := g.sync(generator); err != nil {
		return err
	}

	// write the controller's schema
	for _, descriptor := range spec.Controllers {
		generator = &ControllerGenerator{
			Mode:       ControllerGeneratorModeSchema,
			Path:       dirPkg,
			Controller: descriptor,
		}

		if err := g.sync(generator); err != nil {
			return err
		}
	}

	// write the controller's api
	for _, descriptor := range spec.Controllers {
		generator = &ControllerGenerator{
			Mode:       ControllerGeneratorModeAPI,
			Path:       dirPkg,
			Controller: descriptor,
		}

		if err := g.sync(generator); err != nil {
			return err
		}
	}

	// write the controller's spec
	for _, descriptor := range spec.Controllers {
		generator = &ControllerGenerator{
			Mode:       ControllerGeneratorModeSpec,
			Path:       dirPkg,
			Controller: descriptor,
		}

		if err := g.sync(generator); err != nil {
			return err
		}
	}

	// write the server
	generator = &ServerGenerator{
		Path:        dirPkg,
		Controllers: spec.Controllers,
	}

	if err := g.sync(generator); err != nil {
		return err
	}

	// write the application main
	generator = &MainGenerator{
		Path: dirCmd,
	}

	if err := g.sync(generator); err != nil {
		return err
	}

	return nil
}

func (g *Generator) sync(generator FileGenerator) error {
	if target := generator.Generate(); target != nil {
		// merge if the file exist
		if source, err := OpenFile(target.Name()); err == nil {
			if err := target.Merge(source); err != nil {
				return err
			}
		}

		// mkdir create the directory
		dir := filepath.Dir(target.Name())

		// prepare the service package directory
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		// write the file
		if err := target.Sync(); err != nil {
			return err
		}
	}

	return nil
}
