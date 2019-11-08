package codegen

import (
	"fmt"
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

		inputParam := func(name string, parameters ParameterDescriptorCollection) {
			if node := g.param("Input", name, operation.Name, parameters); node != nil {
				tree = append(tree, node)

				param := &Field{
					Name: name,
					Type: g.inputArg(operation.Name, name),
					Tags: []*Tag{g.tagOfArg(name)},
				}

				input.Fields = append(input.Fields, param)
			}
		}

		for _, request := range operation.Requests {
			// path input
			inputParam("Path", request.Parameters)
			// query input
			inputParam("Query", request.Parameters)
			// header input
			inputParam("Header", request.Parameters)
			// cookie input
			inputParam("Cookie", request.Parameters)

			// input body
			body := &Field{
				Name: "Body",
				Type: request.RequestType.Kind(),
				Tags: []*Tag{g.tagOfArg("Body")},
			}

			input.Fields = append(input.Fields, body)

			// NOTE: we handle the first request for now
			break
		}

		if len(input.Fields) > 0 {
			tree = append(tree, input.Build())
		}

		// output
		output := &StructTypeBuilder{
			Name: operation.Name + "Output",
		}

		output.Commentf("%s is the output of %s operation", output.Name, operation.Name)

		outputParam := func(name string, parameters ParameterDescriptorCollection) {
			if node := g.param("Output", name, operation.Name, parameters); node != nil {
				tree = append(tree, node)

				param := &Field{
					Name: name,
					Type: g.outputArg(operation.Name, name),
					Tags: []*Tag{g.tagOfArg(name)},
				}

				output.Fields = append(output.Fields, param)
			}
		}

		for _, response := range operation.Responses {
			// output header
			outputParam("Header", response.Parameters)

			// output body
			body := &Field{
				Name: "Body",
				Type: response.ResponseType.Kind(),
				Tags: []*Tag{g.tagOfArg("Body")},
			}

			output.Fields = append(output.Fields, body)

			// NOTE: we handle the first response for now
			break
		}

		if len(output.Fields) > 0 {
			tree = append(tree, output.Build())
		}
	}

	return tree
}

func (g *ControllerGenerator) controller() []dst.Decl {
	tree := []dst.Decl{}

	// struct
	builder := &StructTypeBuilder{
		Name: g.name(),
	}

	builder.Commentf("%s is a struct type auto-generated from OpenAPI spec", g.name())
	builder.Commentf(g.Controller.Description)

	tree = append(tree, builder.Build())

	// method mount
	node := g.mount()
	tree = append(tree, node)

	// operations
	for _, operation := range g.Controller.Operations {
		node := g.operation(operation)
		tree = append(tree, node)
	}

	return tree
}

func (g *ControllerGenerator) mount() dst.Decl {
	builder := &MethodTypeBuilder{
		Name: "Mount",
		Receiver: &Param{
			Name: "controller",
			Type: pointer(g.name()),
		},
		Parameters: []*Param{
			&Param{Name: "r", Type: "chi.Router"},
		},
	}

	builder.Commentf("Mount mounts all operations to the corresponding paths")

	// mount body
	block := &BlockBuilder{}

	for _, operation := range g.Controller.Operations {
		method := &Method{
			Receiver: "r",
			Name:     operation.Method,
			Parameters: []string{
				fmt.Sprintf("%q", operation.Path),
				fmt.Sprintf("controller.%s", operation.Name),
			},
		}

		block.Call(method)
	}

	node := builder.Build()
	node.Body = block.Build()

	return node
}

func (g *ControllerGenerator) operation(operation *OperationDescriptor) dst.Decl {
	builder := &MethodTypeBuilder{
		Name: operation.Name,
		Receiver: &Param{
			Name: "controller",
			Type: pointer(g.name()),
		},
		Parameters: []*Param{
			&Param{Name: "w", Type: "http.ResponseWriter"},
			&Param{Name: "r", Type: "*http.Request"},
		},
	}

	builder.Commentf("%s handles endpoint %s %s", operation.Name, operation.Method, operation.Path)

	if operation.Deprecated {
		builder.Commentf("Deprecated: The operation is obsolete")
	}

	builder.Commentf(operation.Description)
	builder.Commentf(operation.Summary)

	return builder.Build()
}

func (g *ControllerGenerator) spec() []dst.Decl {
	//TODO:
	return nil
}

func (g *ControllerGenerator) param(context, kind, name string, parameters ParameterDescriptorCollection) dst.Decl {
	builder := &StructTypeBuilder{
		Name: name + strings.Title(kind) + context,
	}

	builder.Commentf("%s is the %s %s of %s operation", builder.Name, strings.ToLower(context), kind, name)

	for _, param := range parameters {
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

func (g *ControllerGenerator) outputArg(operation, kind string) string {
	return "*" + operation + strings.Title(kind) + "Output"
}

func (g *ControllerGenerator) tagOfArg(kind string) *Tag {
	return &Tag{
		Key:  strings.ToLower(kind),
		Name: "~",
	}
}

// func (builder *ControllerBuilder) mount(receiver string, descriptor *ControllerDescriptor) dst.Decl {
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
