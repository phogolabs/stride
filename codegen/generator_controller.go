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
	builder := NewFileBuilder("service")

	switch g.Mode {
	case ControllerGeneratorModeAPI:
		g.controller(builder)
	case ControllerGeneratorModeSchema:
		g.schema(builder)
	case ControllerGeneratorModeSpec:
		g.spec(builder)
	}

	return builder.Build(filepath.Join(g.Path, g.filename()))
}

func (g *ControllerGenerator) schema(root *FileBuilder) {
	for _, operation := range g.Controller.Operations {
		name := camelize(operation.Name)

		// input
		input := root.Type(name + "Input")
		input.Commentf("%s is the input of %s operation", input.Name(), name)

		for _, request := range operation.Requests {
			// path input
			g.param("Path", root, input, request.Parameters)
			// query input
			g.param("Query", root, input, request.Parameters)
			// header input
			g.param("Header", root, input, request.Parameters)
			// cookie input
			g.param("Cookie", root, input, request.Parameters)

			// input body
			input.Field("Body", request.RequestType.Kind(), g.tagOfArg("Body"))

			// NOTE: we handle the first request for now
			break
		}

		// output
		output := root.Type(operation.Name + "Output")
		output.Commentf("%s is the output of %s operation", output.Name(), name)

		for _, response := range operation.Responses {
			// output header
			g.param("Header", root, output, response.Parameters)

			// output body
			input.Field("Body", response.ResponseType.Kind(), g.tagOfArg("Body"))

			// output status method
			method := output.Method("Status").Return("int")
			method.Commentf("Status returns the response status code")

			// output status mmethod body
			// method.Block(fmt.Sprintf("return %d", response.Code))

			// NOTE: we handle the first response for now
			break
		}
	}
}

func (g *ControllerGenerator) param(kind string, root *FileBuilder, parent *StructTypeBuilder, parameters ParameterDescriptorCollection) {
	builder := root.Type(parent.Name() + kind)
	builder.Commentf("%s is the %s of %s", builder.Name(), strings.ToLower(kind), parent.Name())

	for _, param := range parameters {
		if strings.EqualFold(param.In, kind) {
			builder.Field(param.Name, param.ParameterType.Kind(), param.Tags()...)
		}
	}

	parent.Field(kind, pointer(builder.Name()), g.tagOfArg(kind))
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

	// mount method block
	g.mount(method)

	// operations
	for _, operation := range g.Controller.Operations {
		name := camelize(operation.Name)

		method = builder.Method(operation.Name)
		method.Commentf("%s handles endpoint %s %s", name, operation.Method, operation.Path)

		if operation.Deprecated {
			method.Commentf("Deprecated: The operation is obsolete")
		}

		method.Commentf(operation.Description)
		method.Commentf(operation.Summary)

		method.Param("w", "http.ResponseWriter")
		method.Param("r", "*http.Request")

		g.operation(name, method)
	}
}

func (g *ControllerGenerator) mount(builder *MethodTypeBuilder) {
	// for _, operation := range g.Controller.Operations {
	// var (
	// 	path    = operation.Path
	// 	method  = camelize(strings.ToLower(operation.Method))
	// 	handler = camelize(operation.Name)
	// )

	// builder.Block("r.%s(%q, x.%s)", method, path, handler)
	// }
}

func (g *ControllerGenerator) operation(name string, builder *MethodTypeBuilder) {
	// builder.Block("reactor := restify.NewReactor(w, r)")
	// builder.Block("")
	// builder.Block("var (")
	// builder.Block("   input  = &%sInput{}", name)
	// builder.Block("   output = &%sOutput{}", name)
	// builder.Block(")")
	// builder.Block("")
	// builder.Block("if err := reactor.Bind(input); err != nil {")
	// builder.Block("   reactor.Render(err)")
	// builder.Block("   return")
	// builder.Block("}")
	// builder.Block("")
	// builder.Block("// stride:block open")
	// builder.Block("// TODO: Please add your implementation here")
	// builder.Block("// stride:block close")
	// builder.Block("")
	// builder.Block("if err := reactor.Render(output); err != nil {")
	// builder.Block("   reactor.Render(err)")
	// builder.Block("}")
}

func (g *ControllerGenerator) spec(root *FileBuilder) {
	//TODO:
}

func (g *ControllerGenerator) filename() string {
	name := inflect.Underscore(g.Controller.Name) + "_api"

	switch g.Mode {
	case ControllerGeneratorModeAPI:
		name = name + ".go"
	case ControllerGeneratorModeSchema:
		name = name + "_model.go"
	case ControllerGeneratorModeSpec:
		name = name + "_test.go"
	}

	return name
}

func (g *ControllerGenerator) name() string {
	name := camelize(g.Controller.Name) + "API"
	return name
}

func (g *ControllerGenerator) tagOfArg(kind string) *Tag {
	return &Tag{
		Key:  strings.ToLower(kind),
		Name: "~",
	}
}
