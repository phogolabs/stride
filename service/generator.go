package service

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/phogolabs/stride/codegen"
)

// SpecResolver resolves the spec
type SpecResolver interface {
	// Resolve resolves the spec
	Resolve(spec *openapi3.Swagger) *codegen.SpecDescriptor
}

// CodeGenerator generates the code
type CodeGenerator interface {
	// Generate generates a source code from spec
	Generate(spec *codegen.SpecDescriptor) error
}

// Generator generates the code
type Generator struct {
	Path      string
	Generator CodeGenerator
	Resolver  SpecResolver
}

// Generate generates the source code
func (g *Generator) Generate() error {
	loader := openapi3.NewSwaggerLoader()

	swagger, err := loader.LoadSwaggerFromFile(g.Path)
	if err != nil {
		return err
	}

	spec := g.Resolver.Resolve(swagger)

	if err := g.Generator.Generate(spec); err != nil {
		return err
	}

	return nil
}
