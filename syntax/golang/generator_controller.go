package golang

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/phogolabs/stride/codegen"
	"github.com/phogolabs/stride/inflect"
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
	Controller *codegen.ControllerDescriptor
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
		name := inflect.Camelize(operation.Name)

		// input
		input := root.Struct(name + "Input")
		input.Commentf("It is the input of %s operation", name)

		// path input
		g.param("Path", root, input, operation.Parameters)
		// query input
		g.param("Query", root, input, operation.Parameters)
		// header input
		g.param("Header", root, input, operation.Parameters)
		// cookie input
		g.param("Cookie", root, input, operation.Parameters)

		for _, request := range operation.Requests {
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
			output.AddField("Body", response.ResponseType.Kind(), g.tagOfArg("Body"))

			// output status method
			method := output.
				Function("Status").
				AddReturn("int")

			if body := method.Body(); body != nil {
				body.AppendComment()
				body.Append("return %d", response.Code)

				if err := body.Build(); err != nil {
					panic(err)
				}
			}

			method.Commentf("Status returns the response status code")
			// NOTE: we handle the first response for now
			break
		}
	}
}

func (g *ControllerGenerator) param(kind string, root *File, parent *StructType, parameters codegen.ParameterDescriptorCollection) {
	builder := root.Struct(parent.Name() + kind)
	builder.Commentf("It is the %s of %s", strings.ToLower(kind), parent.Name())

	for _, param := range parameters {
		if strings.EqualFold(param.In, kind) {
			// add a import if needed
			root.AddImport(param.ParameterType.Namespace())

			// add a field
			builder.AddField(param.Name, param.ParameterType.Kind(), param.Tags()...)
		}
	}

	parent.AddField(kind, inflect.Pointer(builder.Name()), g.tagOfArg(kind))
}

func (g *ControllerGenerator) controller(root *File) {
	// add a import if needed
	root.AddImport("github.com/go-chi/chi")
	root.AddImport("github.com/phogolabs/restify")
	root.AddImport("net/http")

	// struct
	builder := root.Struct(g.name())
	builder.Commentf(g.Controller.Description)

	// // method mount
	// method := builder.Function("Mount").AddParam("r", "chi.Router")
	// method.Commentf("Mount mounts all operations to the corresponding paths")

	// // mount method block
	// g.mount(method)

	// operations
	for _, operation := range g.Controller.Operations {

		writer := &TemplateWriter{
			Path: "syntax/golang/operation.go.tpl",
			Context: map[string]interface{}{
				"controller":  builder.Name(),
				"operation":   operation.Name,
				"method":      operation.Method,
				"path":        operation.Path,
				"description": operation.Description,
				"summary":     operation.Summary,
				"deprecated":  operation.DeprecationMessage(),
			},
		}

		buffer := &bytes.Buffer{}

		if _, err := writer.WriteTo(buffer); err != nil {
			panic(err)
		}

		fmt.Println(operation.Name)
		fmt.Println("===========")
		fmt.Println(buffer.String())
	}
}

func (g *ControllerGenerator) mount(builder *FunctionType) {
	body := builder.Body()

	for _, operation := range g.Controller.Operations {
		var (
			path    = operation.Path
			method  = inflect.Camelize(strings.ToLower(operation.Method))
			handler = inflect.Camelize(operation.Name)
		)

		body.Append("r.%s(%q, x.%s)", method, path, handler)
	}

	if err := body.Build(); err != nil {
		panic(err)
	}
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
	name := inflect.Camelize(g.Controller.Name) + "API"
	return name
}

func (g *ControllerGenerator) tagOfArg(kind string) *codegen.TagDescriptor {
	return &codegen.TagDescriptor{
		Key:  strings.ToLower(kind),
		Name: "~",
	}
}
