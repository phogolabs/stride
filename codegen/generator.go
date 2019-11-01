package codegen

import (
	"go/ast"
	"go/printer"
	"go/token"
	"os"
)

// Generator generates the source code
type Generator struct{}

// Generate generates the source code
func (g *Generator) Generate(spec *SpecDescriptor) error {
	fileSet := token.NewFileSet()

	file := &ast.File{
		Name:  ast.NewIdent("service"),
		Scope: ast.NewScope(nil),
		Decls: []ast.Decl{},
	}

	// for _, schema := range spec.Schemas {
	// 	declaration := g.declare(schema)
	// 	file.Decls = append(file.Decls, declaration)
	// }

	printer.Fprint(os.Stdout, fileSet, file)
	return nil
}

func (g *Generator) declare(descriptor *TypeDescriptor) *ast.GenDecl {
	declaration := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(descriptor.Name),
				Type: &ast.StructType{
					Fields:     &ast.FieldList{},
					Incomplete: true,
				},
			},
		},
	}

	return declaration
}
