package codegen

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/go-openapi/inflect"
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

	// write types
	if err := g.write(g.filename("contract"), g.types(spec.Types)); err != nil {
		return err
	}

	builder := &ControllerBuilder{}

	// write controllers
	for _, descriptor := range spec.Controllers {
		name := descriptor.Name + "_api"
		if err := g.write(g.filename(name), builder.Build(descriptor)); err != nil {
			return err
		}

		// spec := descriptor.Name + "_api_test"
		// if err := g.write(g.filename(spec), g.spec(descriptor)); err != nil {
		// 	return err
		// }
	}

	return nil
}

func (g *Generator) write(name string, decls []dst.Decl) error {
	pkg := filepath.Base(filepath.Dir(name))

	if strings.HasSuffix(name, "_test.go") {
		pkg = pkg + "_test"
	}

	root := &dst.File{
		Name: &dst.Ident{
			Name: pkg,
		},
		Decls: decls,
	}

	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := decorator.Fprint(file, root); err != nil {
		return err
	}

	return nil
}

func (g *Generator) types(descriptors TypeDescriptorCollection) []dst.Decl {
	var (
		tree    = []dst.Decl{}
		builder = &TypeBuilder{}
	)

	for _, descriptor := range descriptors {
		node := builder.Build(descriptor)
		tree = append(tree, node)
	}

	return tree
}

func (g *Generator) filename(name string) string {
	name = inflect.Underscore(name) + ".go"
	return filepath.Join(g.Path, "service", name)
}

func element(descriptor *TypeDescriptor) *TypeDescriptor {
	element := descriptor

	for element.IsAlias {
		element = element.Element
	}

	return element
}
