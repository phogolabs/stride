package golang

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/phogolabs/stride/codedom"
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
	Controller *codedom.ControllerDescriptor
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
		input := NewStructType(name + "Input")
		input.Commentf("It is the input of %s operation", name)
		// add the input to the file
		root.AddNode(input)

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
		output := NewStructType(name + "Output")
		output.Commentf("It is the output of %s operation", name)
		// add the output to the file
		root.AddNode(output)

		for _, response := range operation.Responses {
			// output header
			g.param("Header", root, output, response.Parameters)

			// output body
			output.AddField("Body", response.ResponseType.Kind(), g.tagOfArg("Body"))

			writer := &TemplateWriter{
				Path: "syntax/golang/status.go.tpl",
				Context: map[string]interface{}{
					"code":   response.Code,
					"schema": output.Name(),
				},
			}

			buffer := &bytes.Buffer{}
			if _, err := writer.WriteTo(buffer); err != nil {
				panic(err)
			}

			if err := root.AddFunction(buffer.String()); err != nil {
				panic(err)
			}

			// NOTE: we handle the first response for now
			break
		}
	}
}

func (g *ControllerGenerator) param(kind string, root *File, parent *StructType, parameters codedom.ParameterDescriptorCollection) {
	spec := NewStructType(parent.Name() + kind)
	spec.Commentf("It is the %s of %s", strings.ToLower(kind), parent.Name())

	for _, param := range parameters {
		if strings.EqualFold(param.In, kind) {
			// add a import if needed
			root.AddImport(param.ParameterType.Namespace())

			// add a field
			spec.AddField(param.Name, param.ParameterType.Kind(), param.Tags()...)
		}
	}

	if spec.HasFields() {
		// add the spec to the file
		root.AddNode(spec)
		// add the spec as property to the parent
		parent.AddField(kind, inflect.Pointer(spec.Name()), g.tagOfArg(kind))
	}
}

func (g *ControllerGenerator) controller(root *File) {
	// add a import if needed
	root.AddImport("github.com/go-chi/chi")
	root.AddImport("github.com/phogolabs/restify")
	root.AddImport("net/http")

	// struct
	spec := NewStructType(g.name())
	spec.Commentf(g.Controller.Description)
	// add the spec to the file
	root.AddNode(spec)

	// mount method
	writer := &TemplateWriter{
		Path: "syntax/golang/mount.go.tpl",
		Context: map[string]interface{}{
			"controller": spec.Name(),
			"operations": g.Controller.Operations,
		},
	}

	buffer := &bytes.Buffer{}
	if _, err := writer.WriteTo(buffer); err != nil {
		panic(err)
	}

	if err := root.AddFunction(buffer.String()); err != nil {
		panic(err)
	}

	// operations
	for _, operation := range g.Controller.Operations {
		writer := &TemplateWriter{
			Path: "syntax/golang/operation.go.tpl",
			Context: map[string]interface{}{
				"controller":  spec.Name(),
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

		if err := root.AddFunction(buffer.String()); err != nil {
			panic(err)
		}
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

func (g *ControllerGenerator) tagOfArg(kind string) *codedom.TagDescriptor {
	return &codedom.TagDescriptor{
		Key:  strings.ToLower(kind),
		Name: "~",
	}
}
