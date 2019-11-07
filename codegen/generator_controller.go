package codegen

import (
	"path/filepath"
	"strings"

	"github.com/dave/dst"
	"github.com/go-openapi/inflect"
)

// ControllerGeneratorMode determines the mode of this generator
type ControllerGeneratorMode byte

const (
	// ControllerGeneratorModeSchema generates the schema for the controller
	ControllerGeneratorModeSchema ControllerGeneratorMode = 0
	// ControllerGeneratorModeAPI generates the api for the controller
	ControllerGeneratorModeAPI ControllerGeneratorMode = 1
	// ControllerGeneratorModeSpec generates the spec for the controller
	ControllerGeneratorModeSpec ControllerGeneratorMode = 2
)

// ControllerGenerator builds a controller
type ControllerGenerator struct {
	Mode       ControllerGeneratorMode
	Path       string
	Controller *ControllerDescriptor
}

// Generate generates a file
func (g *ControllerGenerator) Generate() *File {
	builder := &FileBuilder{
		Package: "service",
	}

	// create the file
	file := builder.Build()

	switch g.Mode {
	case ControllerGeneratorModeAPI:
		file.Decls = g.controller()
	case ControllerGeneratorModeSchema:
		file.Decls = g.schema()
	case ControllerGeneratorModeSpec:
		file.Decls = g.spec()
	}

	return &File{
		Name:    filepath.Join(g.Path, g.filename()),
		Content: file,
	}
}

func (g *ControllerGenerator) schema() []dst.Decl {
	tree := []dst.Decl{}

	for _, operation := range g.Controller.Operations {
		// input
		input := &StructTypeBuilder{
			Name: operation.Name + "Input",
		}

		input.Commentf("%s is the input of %s operation", input.Name, operation.Name)
		input.Commentf(operation.Description)

		param := func(name string) {
			if node := g.param(name, operation); node != nil {
				tree = append(tree, node)

				field := &Field{
					Name: name,
					Type: g.inputArg(operation.Name, name),
				}

				input.Fields = append(input.Fields, field)
			}
		}

		// path input
		param("Path")
		// query input
		param("Query")
		// header input
		param("Header")
		// cookie input
		param("Cookie")

		// input body
		for _, request := range operation.Requests {
			field := &Field{
				Name: "Body",
				Type: request.RequestType.Kind(),
			}

			input.Fields = append(input.Fields, field)
			// NOTE: we handle the first request for now
			break
		}

		tree = append(tree, input.Build())

		// output
		output := &StructTypeBuilder{
			Name: operation.Name + "Output",
		}

		output.Commentf("%s is the output of %s operation", output.Name, operation.Name)

		//TODO: output header

		// output body
		for _, response := range operation.Responses {
			field := &Field{
				Name: "Body",
				Type: response.ResponseType.Kind(),
			}

			output.Fields = append(input.Fields, field)
			// NOTE: we handle the first response for now
			break
		}

		tree = append(tree, output.Build())
	}

	return tree
}

func (g *ControllerGenerator) controller() []dst.Decl {
	builder := &StructTypeBuilder{
		Name: g.name(),
	}

	builder.Commentf("%s is a struct type auto-generated from OpenAPI spec", g.name())
	builder.Commentf(g.Controller.Description)

	return []dst.Decl{builder.Build()}
}

func (g *ControllerGenerator) spec() []dst.Decl {
	return nil
}

func (g *ControllerGenerator) param(kind string, operation *OperationDescriptor) dst.Decl {
	builder := &StructTypeBuilder{
		Name: operation.Name + strings.Title(kind) + "Input",
	}

	builder.Commentf("%s is the %s input of %s operation", builder.Name, kind, operation.Name)

	for _, param := range operation.Parameters {
		if strings.EqualFold(param.In, kind) {
			field := &Field{
				Name: param.Name,
				Type: param.ParameterType.Kind(),
				Tags: param.Tags(),
			}

			builder.Fields = append(builder.Fields, field)
		}
	}

	if len(builder.Fields) == 0 {
		return nil
	}

	return builder.Build()
}

func (g *ControllerGenerator) filename() string {
	name := g.Controller.Name + "_api"

	switch g.Mode {
	case ControllerGeneratorModeAPI:
		name = name + ".go"
	case ControllerGeneratorModeSchema:
		name = name + "_model.go"
	case ControllerGeneratorModeSpec:
		name = name + "_test.go"
	}

	return inflect.Underscore(name)
}

func (g *ControllerGenerator) name() string {
	name := g.Controller.Name + "API"
	return name
}

func (g *ControllerGenerator) inputArg(operation, kind string) string {
	return "*" + operation + strings.Title(kind) + "Input"
}

// 	var (
// 		name = descriptor.Name + "API"
// 		tree = []dst.Decl{}
// 		node dst.Decl
// 	)

// 	// generate controller type
// 	node = builder.controller(name, descriptor)
// 	tree = append(tree, node)

// 	// generate mount type
// 	node = builder.mount(name, descriptor)
// 	tree = append(tree, node)

// 	// generate operations
// 	for _, operation := range descriptor.Operations {
// 		node := builder.operation(name, operation)
// 		tree = append(tree, node)
// 	}

// 	return tree
// }

// func (builder *ControllerBuilder) controller(name string, descriptor *ControllerDescriptor) dst.Decl {
// 	parent := &StructTypeBuilder{
// 		Name: name,
// 	}

// 	parent.Comments = append(parent.Comments, commentf("%s is a struct type auto-generated from OpenAPI spec", name))

// 	if descriptor.Description != "" {
// 		parent.Comments = append(parent.Comments, commentf(descriptor.Description))
// 	}

// 	return parent.Build()
// }

// func (builder *ControllerBuilder) mount(receiver string, descriptor *ControllerDescriptor) dst.Decl {
// 	node := &dst.FuncDecl{
// 		// receiver param
// 		Recv: &dst.FieldList{
// 			List: []*dst.Field{
// 				&dst.Field{
// 					Names: []*dst.Ident{
// 						&dst.Ident{
// 							Name: "controller",
// 						},
// 					},
// 					Type: &dst.Ident{
// 						Name: pointer(receiver),
// 					},
// 				},
// 			},
// 		},
// 		// function name
// 		Name: &dst.Ident{
// 			Name: "Mount",
// 		},
// 		// function parameters
// 		Type: &dst.FuncType{
// 			Params: &dst.FieldList{
// 				List: []*dst.Field{
// 					&dst.Field{
// 						Names: []*dst.Ident{
// 							&dst.Ident{
// 								Name: "r",
// 							},
// 						},
// 						Type: &dst.Ident{
// 							Name: "chi.Router",
// 						},
// 					},
// 				},
// 			},
// 		},
// 		Body: &dst.BlockStmt{},
// 	}

// 	// function body
// 	for _, operation := range descriptor.Operations {
// 		stmt := &dst.ExprStmt{
// 			X: &dst.CallExpr{
// 				Fun: &dst.SelectorExpr{
// 					X: &dst.Ident{
// 						Name: "r",
// 					},
// 					Sel: &dst.Ident{
// 						Name: operation.Method,
// 					},
// 				},
// 				Args: []dst.Expr{
// 					&dst.BasicLit{
// 						Kind:  token.STRING,
// 						Value: fmt.Sprintf("%q", operation.Path),
// 					},
// 					&dst.SelectorExpr{
// 						X: &dst.Ident{
// 							Name: "controller",
// 						},
// 						Sel: &dst.Ident{
// 							Name: operation.Name,
// 						},
// 					},
// 				},
// 			},
// 		}

// 		node.Body.List = append(node.Body.List, stmt)
// 	}

// 	node.Decs.Before = dst.NewLine
// 	node.Decs.Start.Append(commentf("Mount mounts all operations to the corresponding paths"))
// 	node.Decs.Start.Append(commentf("stride:generate"))

// 	return node
// }

// func (builder *ControllerBuilder) operation(receiver string, descriptor *OperationDescriptor) dst.Decl {
// 	node := &dst.FuncDecl{
// 		// receiver param
// 		Recv: &dst.FieldList{
// 			List: []*dst.Field{
// 				&dst.Field{
// 					Names: []*dst.Ident{
// 						&dst.Ident{
// 							Name: "controller",
// 						},
// 					},
// 					Type: &dst.Ident{
// 						Name: pointer(receiver),
// 					},
// 				},
// 			},
// 		},
// 		// function name
// 		Name: &dst.Ident{
// 			Name: descriptor.Name,
// 		},
// 		// function parameters
// 		Type: &dst.FuncType{
// 			Params: &dst.FieldList{
// 				List: []*dst.Field{
// 					&dst.Field{
// 						Names: []*dst.Ident{
// 							&dst.Ident{
// 								Name: "w",
// 							},
// 						},
// 						Type: &dst.Ident{
// 							Name: "http.ResponseWriter",
// 						},
// 					},
// 					&dst.Field{
// 						Names: []*dst.Ident{
// 							&dst.Ident{
// 								Name: "r",
// 							},
// 						},
// 						Type: &dst.Ident{
// 							Name: "*http.Request",
// 						},
// 					},
// 				},
// 			},
// 		},
// 		Body: &dst.BlockStmt{},
// 	}

// 	node.Decs.Before = dst.NewLine
// 	node.Decs.Start.Append(commentf("%s handles endpoint %s %s", descriptor.Name, descriptor.Method, descriptor.Path))

// 	if descriptor.Deprecated {
// 		node.Decs.Start.Append(commentf("Deprecated: The operation is obsolete"))
// 	}

// 	if descriptor.Description != "" {
// 		node.Decs.Start.Append(commentf(descriptor.Description))
// 	}

// 	if descriptor.Summary != "" {
// 		node.Decs.Start.Append(commentf(descriptor.Summary))
// 	}

// 	node.Decs.Start.Append(commentf("stride:generate"))

// 	return node
// }
