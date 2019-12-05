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

//go:generate counterfeiter -fake-name CodeGenerator -o ../fake/code_generator.go . CodeGenerator

// CodeGenerator generates the code
type CodeGenerator interface {
	// Generate generates a source code from spec
	Generate(spec *codedom.SpecDescriptor) error
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

	spec, err := g.Resolver.Resolve(swagger)
	if err != nil {
		return err
	}

	if err := g.Generator.Generate(spec); err != nil {
		return err
	}

	return nil
}
