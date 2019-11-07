package codegen

import (
	"os"
	"path/filepath"

	"github.com/dave/dst/decorator"
)

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

	contract := &ContractGenerator{
		Path:       path,
		Collection: spec.Types,
	}

	// write the common types
	if err := g.write(contract.Generate()); err != nil {
		return err
	}

	// write the controller's schema
	for _, descriptor := range spec.Controllers {
		controller := &ControllerGenerator{
			Mode:       ControllerGeneratorModeSchema,
			Path:       path,
			Controller: descriptor,
		}

		if err := g.write(controller.Generate()); err != nil {
			return err
		}
	}

	// write the controller's api
	for _, descriptor := range spec.Controllers {
		controller := &ControllerGenerator{
			Mode:       ControllerGeneratorModeAPI,
			Path:       path,
			Controller: descriptor,
		}

		if err := g.write(controller.Generate()); err != nil {
			return err
		}
	}

	// write the controller's spec
	for _, descriptor := range spec.Controllers {
		controller := &ControllerGenerator{
			Mode:       ControllerGeneratorModeSpec,
			Path:       path,
			Controller: descriptor,
		}

		if err := g.write(controller.Generate()); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) write(file *File) error {
	writer, err := os.Create(file.Name)
	if err != nil {
		return err
	}
	defer writer.Close()

	if err := decorator.Fprint(writer, file.Content); err != nil {
		return err
	}

	return nil
}

func element(descriptor *TypeDescriptor) *TypeDescriptor {
	element := descriptor

	for element.IsAlias {
		element = element.Element
	}

	return element
}
