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
	Path       string
	Mode       ControllerGeneratorMode
	Controller *ControllerDescriptor
}

// Generate generates a file
func (g *ControllerGenerator) Generate() *File {
	var (
		filename = filepath.Join(g.Path, g.filename())
		root     = NewFile(filename)
	)

	switch g.Mode {
	case ControllerGeneratorModeAPI:
		g.controller(root)
	case ControllerGeneratorModeSchema:
		g.schema(root)
	case ControllerGeneratorModeSpec:
		g.spec(root)
	}

	return root
}

func (g *ControllerGenerator) schema(root *File) {
	for _, operation := range g.Controller.Operations {
		name := camelize(operation.Name)

		// input
		input := root.Struct(name + "Input")
		input.Commentf("It is the input of %s operation", name)

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
			input.AddField("Body", request.RequestType.Kind(), g.tagOfArg("Body"))

			// NOTE: we handle the first request for now
			break
		}

		// output
		output := root.Struct(name + "Output")
		output.Commentf("It is the output of %s operation", name)

		for _, response := range operation.Responses {
			// output header
			g.param("Header", root, output, response.Parameters)

			// output body
			input.AddField("Body", response.ResponseType.Kind(), g.tagOfArg("Body"))

			// output status method
			method := output.
				Function("Status").
				AddReturn("int").
				Body("return %d", response.Code)

			method.Commentf("Status returns the response status code")
			// NOTE: we handle the first response for now
			break
		}
	}
}

func (g *ControllerGenerator) param(kind string, root *File, parent *StructType, parameters ParameterDescriptorCollection) {
	builder := root.Struct(parent.Name() + kind)
	builder.Commentf("It is the %s of %s", strings.ToLower(kind), parent.Name())

	for _, param := range parameters {
		if strings.EqualFold(param.In, kind) {
			builder.AddField(param.Name, param.ParameterType.Kind(), param.Tags()...)
		}
	}

	parent.AddField(kind, pointer(builder.Name()), g.tagOfArg(kind))
}

func (g *ControllerGenerator) controller(root *File) {
	// struct
	builder := root.Struct(g.name())
	builder.Commentf(g.Controller.Description)

	// method mount
	method := builder.Function("Mount").AddParam("r", "chi.Router")
	method.Commentf("Mount mounts all operations to the corresponding paths")

	// mount method block
	g.mount(method)

	// operations
	for _, operation := range g.Controller.Operations {
		name := camelize(operation.Name)

		method = builder.Function(operation.Name).
			AddParam("w", "http.ResponseWriter").
			AddParam("r", "*http.Request")

		method.Commentf("%s handles endpoint %s %s", name, operation.Method, operation.Path)

		if operation.Deprecated {
			method.Commentf("Deprecated: The operation is obsolete")
		}

		method.Commentf(operation.Description)
		method.Commentf(operation.Summary)

		g.operation(name, method)
	}
}

func (g *ControllerGenerator) mount(builder *FunctionType) {
	buffer := NewBlockWriter()

	for _, operation := range g.Controller.Operations {
		var (
			path    = operation.Path
			method  = camelize(strings.ToLower(operation.Method))
			handler = camelize(operation.Name)
		)

		buffer.Write("r.%s(%q, x.%s)", method, path, handler)
	}

	builder.Body(buffer.String())
}

func (g *ControllerGenerator) operation(name string, builder *FunctionType) {
	buffer := NewBlockWriter()

	buffer.Write("reactor := restify.NewReactor(w, r)")
	buffer.Write("")
	buffer.Write("var (")
	buffer.Write("   input  = &%sInput{}", name)
	buffer.Write("   output = &%sOutput{}", name)
	buffer.Write(")")
	buffer.Write("")
	buffer.Write("if err := reactor.Bind(input); err != nil {")
	buffer.Write("   reactor.Render(err)")
	buffer.Write("   return")
	buffer.Write("}")
	buffer.Write("")
	buffer.Write("// stride:define:block:start body")
	buffer.Write("// TODO: Please add your implementation here")
	buffer.Write("// stride:define:block:end body")
	buffer.Write("")
	buffer.Write("if err := reactor.Render(output); err != nil {")
	buffer.Write("   reactor.Render(err)")
	buffer.Write("}")

	// define the block
	builder.Body(buffer.String())
}

func (g *ControllerGenerator) spec(root *File) {
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

func (g *ControllerGenerator) tagOfArg(kind string) *TagDescriptor {
	return &TagDescriptor{
		Key:  strings.ToLower(kind),
		Name: "~",
	}
}
