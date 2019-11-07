package codegen

import (
	"fmt"
	"go/token"
	"strings"

	"github.com/dave/dst"
	"github.com/go-openapi/inflect"
)

// ControllerBuilder builds a controller
type ControllerBuilder struct {
	Package string
}

// Build builds a type
func (builder *ControllerBuilder) Build(descriptor *ControllerDescriptor) []dst.Decl {
	var (
		name = camelize(descriptor.Name) + "API"
		tree = []dst.Decl{}
		node dst.Decl
	)

	// generate controller type
	node = builder.controller(name, descriptor)
	tree = append(tree, node)

	// generate mount type
	node = builder.mount(name, descriptor)
	tree = append(tree, node)

	// generate operations
	for _, operation := range descriptor.Operations {
		//TODO: generate input and output
		node := builder.operation(name, operation)
		tree = append(tree, node)
	}

	return tree
}

func (builder *ControllerBuilder) controller(name string, descriptor *ControllerDescriptor) dst.Decl {
	node := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: &dst.Ident{
					Name: name,
				},
				Type: &dst.StructType{
					Fields: &dst.FieldList{
						List: []*dst.Field{},
					},
					Incomplete: true,
				},
			},
		},
	}

	node.Decs.Before = dst.NewLine
	node.Decs.Start.Append(commentf("%s is a struct type auto-generated from OpenAPI spec", name))

	if descriptor.Description != "" {
		node.Decs.Start.Append(commentf(descriptor.Description))
	}

	node.Decs.Start.Append(commentf("stride:generate"))
	node.Decs.After = dst.EmptyLine

	return node
}

func (builder *ControllerBuilder) mount(receiver string, descriptor *ControllerDescriptor) dst.Decl {
	node := &dst.FuncDecl{
		// receiver param
		Recv: &dst.FieldList{
			List: []*dst.Field{
				&dst.Field{
					Names: []*dst.Ident{
						&dst.Ident{
							Name: "controller",
						},
					},
					Type: &dst.Ident{
						Name: pointer(receiver),
					},
				},
			},
		},
		// function name
		Name: &dst.Ident{
			Name: "Mount",
		},
		// function parameters
		Type: &dst.FuncType{
			Params: &dst.FieldList{
				List: []*dst.Field{
					&dst.Field{
						Names: []*dst.Ident{
							&dst.Ident{
								Name: "r",
							},
						},
						Type: &dst.Ident{
							Name: "chi.Router",
						},
					},
				},
			},
		},
		Body: &dst.BlockStmt{},
	}

	// function body
	for _, operation := range descriptor.Operations {
		stmt := &dst.ExprStmt{
			X: &dst.CallExpr{
				Fun: &dst.SelectorExpr{
					X: &dst.Ident{
						Name: "r",
					},
					Sel: &dst.Ident{
						Name: inflect.Camelize(strings.ToLower(operation.Method)),
					},
				},
				Args: []dst.Expr{
					&dst.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("%q", operation.Path),
					},
					&dst.SelectorExpr{
						X: &dst.Ident{
							Name: "controller",
						},
						Sel: &dst.Ident{
							Name: camelize(operation.Name),
						},
					},
				},
			},
		}

		node.Body.List = append(node.Body.List, stmt)
	}

	node.Decs.Before = dst.NewLine
	node.Decs.Start.Append(commentf("Mount mounts all operations to the corresponding paths"))
	node.Decs.Start.Append(commentf("stride:generate"))

	return node
}

func (builder *ControllerBuilder) operation(receiver string, descriptor *OperationDescriptor) dst.Decl {
	name := camelize(descriptor.Name)

	node := &dst.FuncDecl{
		// receiver param
		Recv: &dst.FieldList{
			List: []*dst.Field{
				&dst.Field{
					Names: []*dst.Ident{
						&dst.Ident{
							Name: "controller",
						},
					},
					Type: &dst.Ident{
						Name: pointer(receiver),
					},
				},
			},
		},
		// function name
		Name: &dst.Ident{
			Name: name,
		},
		// function parameters
		Type: &dst.FuncType{
			Params: &dst.FieldList{
				List: []*dst.Field{
					&dst.Field{
						Names: []*dst.Ident{
							&dst.Ident{
								Name: "w",
							},
						},
						Type: &dst.Ident{
							Name: "http.ResponseWriter",
						},
					},
					&dst.Field{
						Names: []*dst.Ident{
							&dst.Ident{
								Name: "r",
							},
						},
						Type: &dst.Ident{
							Name: "*http.Request",
						},
					},
				},
			},
		},
		Body: &dst.BlockStmt{},
	}

	node.Decs.Before = dst.NewLine
	node.Decs.Start.Append(commentf("%s handles endpoint %s %s", name, descriptor.Method, descriptor.Path))

	if descriptor.Deprecated {
		node.Decs.Start.Append(commentf("Deprecated: The operation is obsolete"))
	}

	if descriptor.Description != "" {
		node.Decs.Start.Append(commentf(descriptor.Description))
	}

	if descriptor.Summary != "" {
		node.Decs.Start.Append(commentf(descriptor.Summary))
	}

	node.Decs.Start.Append(commentf("stride:generate"))

	return node
}
