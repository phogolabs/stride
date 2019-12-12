package service

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/phogolabs/stride/codedom"
)

//go:generate counterfeiter -fake-name SpecResolver -o ../fake/spec_resolver.go . SpecResolver

// SpecResolver resolves the spec
type SpecResolver interface {
	// Resolve resolves the spec
	Resolve(spec *openapi3.Swagger) (*codedom.SpecDescriptor, error)
}

//go:generate counterfeiter -fake-name SyntaxGenerator -o ../fake/syntax_generator.go . SyntaxGenerator

// SyntaxGenerator generates the code
type SyntaxGenerator interface {
	// Generate generates a source code from spec
	Generate(spec *codedom.SpecDescriptor) error
}

var _ SyntaxGenerator = CompositeGenerator{}

// CompositeGenerator represents a composite generator
type CompositeGenerator []SyntaxGenerator

// Generate generates the source code
func (items CompositeGenerator) Generate(spec *codedom.SpecDescriptor) error {
	for _, generator := range items {
		if err := generator.Generate(spec); err != nil {
			return err
		}
	}

	return nil
}

// Generator generates the code
type Generator struct {
	Path      string
	Generator SyntaxGenerator
	Resolver  SpecResolver
}

// Generate generates the source code
func (g *Generator) Generate() error {
	loader := openapi3.NewSwaggerLoader()

	swagger, err := loader.LoadSwaggerFromFile(g.Path)
	if err != nil {
		return err
	}

	spec, err := g.Resolver.Resolve(swagger)
	if err != nil {
		return err
	}

	if err := g.Generator.Generate(spec); err != nil {
		return err
	}

	return nil
}
