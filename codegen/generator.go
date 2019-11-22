package codegen

import (
	"os"
	"path/filepath"
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
func (g *Generator) Generate(spec *SpecDescriptor) error {
	path := filepath.Join(g.Path, "service")

	// prepare the service package directory
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	generator := &SchemaGenerator{
		Path:       path,
		Collection: spec.Types,
	}

	if err := g.sync(generator); err != nil {
		return err
	}

	// write the controller's schema
	for _, descriptor := range spec.Controllers {
		generator := &ControllerGenerator{
			Mode:       ControllerGeneratorModeSchema,
			Path:       path,
			Controller: descriptor,
		}

		if err := g.sync(generator); err != nil {
			return err
		}
	}

	// write the controller's api
	for _, descriptor := range spec.Controllers {
		generator := &ControllerGenerator{
			Mode:       ControllerGeneratorModeAPI,
			Path:       path,
			Controller: descriptor,
		}

		if err := g.sync(generator); err != nil {
			return err
		}
	}

	// write the controller's spec
	for _, descriptor := range spec.Controllers {
		generator := &ControllerGenerator{
			Mode:       ControllerGeneratorModeSpec,
			Path:       path,
			Controller: descriptor,
		}

		if err := g.sync(generator); err != nil {
			return err
		}
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

		// write the file
		if err := target.Sync(); err != nil {
			return err
		}
	}

	return nil
}
