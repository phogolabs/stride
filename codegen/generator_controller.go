package codegen

import (
	"path/filepath"
	"strings"

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

	switch g.Mode {
	case ControllerGeneratorModeAPI:
		g.controller(builder)
	case ControllerGeneratorModeSchema:
		g.schema(builder)
	case ControllerGeneratorModeSpec:
		g.spec(builder)
	}

	return &File{
		Name:    filepath.Join(g.Path, g.filename()),
		Content: builder.Build(),
	}
}

func (g *ControllerGenerator) schema(root *FileBuilder) {
	param := func(kind string, parent *StructTypeBuilder, parameters ParameterDescriptorCollection) {
		builder := root.Type(parent.Name + kind)
		builder.Commentf("%s is the %s of %s", builder.Name, strings.ToLower(kind), parent.Name)

		for _, param := range parameters {
			if strings.EqualFold(param.In, kind) {
				builder.Field(param.Name, param.ParameterType.Kind(), param.Tags()...)
			}
		}

		parent.Field(kind, pointer(builder.Name), g.tagOfArg(kind))
	}

	for _, operation := range g.Controller.Operations {
		// input
		input := root.Type(operation.Name + "Input")
		input.Commentf("%s is the input of %s operation", input.Name, operation.Name)

		for _, request := range operation.Requests {
			// path input
			param("Path", input, request.Parameters)
			// query input
			param("Query", input, request.Parameters)
			// header input
			param("Header", input, request.Parameters)
			// cookie input
			param("Cookie", input, request.Parameters)

			// input body
			input.Field("Body", request.RequestType.Kind(), g.tagOfArg("Body"))

			// NOTE: we handle the first request for now
			break
		}

		// output
		output := root.Type(operation.Name + "Output")
		output.Commentf("%s is the output of %s operation", output.Name, operation.Name)

		for _, response := range operation.Responses {
			// output header
			param("Header", output, response.Parameters)

			// output body
			input.Field("Body", response.ResponseType.Kind(), g.tagOfArg("Body"))

			// NOTE: we handle the first response for now
			break
		}
	}
}

func (g *ControllerGenerator) controller(root *FileBuilder) {
	// struct
	builder := root.Type(g.name())
	builder.Commentf("%s is a struct type auto-generated from OpenAPI spec", g.name())
	builder.Commentf(g.Controller.Description)

	// method mount
	method := builder.Method("Mount")
	method.Commentf("Mount mounts all operations to the corresponding paths")
	method.Param("r", "chi.Router")

	// operations
	for _, operation := range g.Controller.Operations {
		method = builder.Method(operation.Name)
		method.Commentf("%s handles endpoint %s %s", operation.Name, operation.Method, operation.Path)

		if operation.Deprecated {
			method.Commentf("Deprecated: The operation is obsolete")
		}

		method.Commentf(operation.Description)
		method.Commentf(operation.Summary)

		method.Param("w", "http.ResponseWriter")
		method.Param("r", "*http.Request")
	}
}

// func (g *ControllerGenerator) mount() dst.Decl {
// 	// mount body
// 	block := &BlockBuilder{}

// 	for _, operation := range g.Controller.Operations {
// 		method := &Method{
// 			Receiver: "r",
// 			Name:     operation.Method,
// 			Parameters: []string{
// 				fmt.Sprintf("%q", operation.Path),
// 				fmt.Sprintf("controller.%s", operation.Name),
// 			},
// 		}

// 		block.Call(method)
// 	}

// 	node := builder.Build()
// 	node.Body = block.Build()

// 	return node
// }

func (g *ControllerGenerator) spec(root *FileBuilder) {
	//TODO:
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

func (g *ControllerGenerator) tagOfArg(kind string) *Tag {
	return &Tag{
		Key:  strings.ToLower(kind),
		Name: "~",
	}
}
