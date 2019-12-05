package golang

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/phogolabs/stride/codedom"
	"github.com/phogolabs/stride/contract"
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
	Reporter   contract.Reporter
}

// Generate generates a file
func (g *ControllerGenerator) Generate() *File {
	var (
		filename = filepath.Join(g.Path, g.filename())
		root     = NewFile(filename)
	)

	reporter := g.Reporter.With(contract.SeverityHigh)

	reporter.Notice(" Generating controller: %s file: %s...",
		inflect.Dasherize(g.name()),
		root.Name(),
	)

	defer reporter.Notice(" Generating controller: %s file: %s successful",
		inflect.Dasherize(g.name()),
		root.Name(),
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
	reporter := g.Reporter.With(contract.SeverityHigh)

	reporter.Notice(" Generating controller: %s schema...", inflect.Dasherize(g.name()))
	defer reporter.Success("ﳑ Generating controller: %s schema successful", inflect.Dasherize(g.name()))

	for _, operation := range g.Controller.Operations {
		name := inflect.Camelize(operation.Name)

		g.Reporter.Info("ﳑ Generating controller: %s operation: %s schema...",
			inflect.Dasherize(g.name()),
			inflect.Dasherize(operation.Name),
		)

		g.Reporter.Info("ﳑ Generating controller: %s operation: %s schema input...",
			inflect.Dasherize(g.name()),
			inflect.Dasherize(operation.Name),
		)

		for index, request := range operation.Requests {
			reporter := g.Reporter.With(contract.SeverityLow)

			if index > 0 {
				reporter.Warn("ﳑ Generating request content-type: %s skipped. More than one request per operation is not supported",
					inflect.Dasherize(request.ContentType),
				)

				continue
			}

			// input
			input := NewStructType(name + "Input")
			input.Commentf("It is the input of %s operation", name)
			// add the input to the file
			root.AddNode(input)

			// path input
			g.param("Path", root, input, request.Parameters)
			// query input
			g.param("Query", root, input, request.Parameters)
			// header input
			g.param("Header", root, input, request.Parameters)
			// cookie input
			g.param("Cookie", root, input, request.Parameters)

			if request.RequestType != nil {
				reporter.Info("ﳑ Generating type: %s field: %s content-type: %s...",
					inflect.Dasherize(input.Name()),
					inflect.Dasherize("body"),
					inflect.Dasherize(request.ContentType),
				)

				// input body
				input.AddField("Body", request.RequestType.Kind(), g.tagOfArg("Body"))

				reporter.Info("ﳑ Generating type: %s field: %s content-type: %s successful",
					inflect.Dasherize(input.Name()),
					inflect.Dasherize("body"),
					inflect.Dasherize(request.ContentType),
				)

				g.function(root, "decode", map[string]interface{}{
					"receiver": input.Name(),
					"function": "decode",
					"body":     request.RequestType.Name,
				})
			}
		}

		g.Reporter.Success("ﳑ Generating controller: %s operation: %s schema input successful",
			inflect.Dasherize(g.name()),
			inflect.Dasherize(operation.Name),
		)

		g.Reporter.Info("ﳑ Generating controller: %s operation: %s schema output...",
			inflect.Dasherize(g.name()),
			inflect.Dasherize(operation.Name),
		)

		for _, response := range operation.Responses {
			reporter := g.Reporter.With(contract.SeverityLow)

			// output
			output := NewStructType(name + inflect.Camelize(http.StatusText(response.Code)) + "Output")
			output.Commentf("It is the output of %s operation with code: %d", name, response.Code)
			// add the output to the file
			root.AddNode(output)

			// output header
			g.param("Header", root, output, response.Parameters)

			reporter.Info("ﳑ Generating type: %s field: %s content-type: %s code: %d...",
				inflect.Dasherize(output.Name()),
				inflect.Dasherize("body"),
				inflect.Dasherize(response.ContentType),
				response.Code,
			)

			if response.ResponseType != nil {
				// output body
				output.AddField("Body", response.ResponseType.Kind(), g.tagOfArg("Body"))

				g.function(root, "encode", map[string]interface{}{
					"receiver": output.Name(),
					"function": "encode",
					"body":     response.ResponseType.Name,
				})
			}

			reporter.Info("ﳑ Generating type: %s field: %s content-type: %s code: %d successful",
				inflect.Dasherize(output.Name()),
				inflect.Dasherize("body"),
				inflect.Dasherize(response.ContentType),
				response.Code,
			)

			g.function(root, "status", map[string]interface{}{
				"receiver": output.Name(),
				"function": "status",
				"code":     response.Code,
			})

			if response.IsDefault {
				// output
				alias := NewLiteralType(name + "Output")
				alias.Element(output.Name())
				alias.Commentf("It is the alias to the default output of %s operation", name)

				reporter.Info("ﳑ Generating type: %s...", inflect.Dasherize(alias.Name()))

				// add the output to the file
				root.AddNode(alias)

				reporter.Info("ﳑ Generating type: %s successful", inflect.Dasherize(alias.Name()))
			}
		}

		g.Reporter.Success("ﳑ Generating controller: %s operation: %s schema output successful",
			inflect.Dasherize(g.Controller.Name),
			inflect.Dasherize(operation.Name),
		)

		g.Reporter.Success("ﳑ Generating controller: %s operation: %s schema successful",
			inflect.Dasherize(g.Controller.Name),
			inflect.Dasherize(operation.Name),
		)
	}
}

func (g *ControllerGenerator) param(kind string, root *File, parent *StructType, parameters codedom.ParameterDescriptorCollection) {
	spec := NewStructType(parent.Name() + kind)
	spec.Commentf("It is the %s of %s", strings.ToLower(kind), parent.Name())

	reporter := g.Reporter.With(contract.SeverityLow)
	reporter.Info("ﳑ Generating type: %s...", inflect.Dasherize(spec.Name()))
	defer reporter.Success("ﳑ Generating type: %s successful", inflect.Dasherize(spec.Name()))

	for _, param := range parameters {
		if strings.EqualFold(param.In, kind) {
			reporter.Info("ﳑ Generating type: %s field: %s...",
				inflect.Dasherize(spec.Name()),
				inflect.Dasherize(param.Name),
			)

			// add a import if needed
			root.AddImport(param.ParameterType.Namespace())

			// add a field
			spec.AddField(param.Name, param.ParameterType.Kind(), param.Tags()...)

			reporter.Success("ﳑ Generating type: %s field: %s success...",
				inflect.Dasherize(spec.Name()),
				inflect.Dasherize(param.Name),
			)
		}
	}

	if spec.HasFields() {
		// add the spec to the file
		root.AddNode(spec)

		reporter.Info("ﳑ Generating type: %s field: %s...",
			inflect.Dasherize(parent.Name()),
			inflect.Dasherize(kind),
		)

		// add the spec as property to the parent
		parent.AddField(kind, inflect.Pointer(spec.Name()), g.tagOfArg(kind))

		reporter.Success("ﳑ Generating type: %s field: %s successful",
			inflect.Dasherize(parent.Name()),
			inflect.Dasherize(kind),
		)
	}
}

func (g *ControllerGenerator) controller(root *File) {
	reporter := g.Reporter.With(contract.SeverityHigh)

	reporter.Notice("ﳑ Generating controller: %s...", inflect.Dasherize(g.name()))
	defer reporter.Success("ﳑ Generating controller: %s successful", inflect.Dasherize(g.name()))

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
	g.function(root, "mount", map[string]interface{}{
		"receiver":   spec.Name(),
		"function":   "mount",
		"operations": g.Controller.Operations,
	})

	// operations
	for _, operation := range g.Controller.Operations {
		g.function(root, "operation", map[string]interface{}{
			"receiver":    spec.Name(),
			"function":    operation.Name,
			"method":      operation.Method,
			"path":        operation.Path,
			"description": operation.Description,
			"summary":     operation.Summary,
			"deprecated":  operation.DeprecationMessage(),
		})

	}
}

func (g *ControllerGenerator) spec(root *File) {
	//TODO:
}

func (g *ControllerGenerator) function(root *File, name string, ctx map[string]interface{}) {
	var (
		receiver  = ctx["receiver"].(string)
		operation = ctx["function"].(string)
	)

	g.Reporter.Info("ﳑ Generating type: %s function: %s...",
		inflect.Dasherize(receiver),
		inflect.Dasherize(operation),
	)

	// mount method
	writer := &TemplateWriter{
		Path:    fmt.Sprintf("syntax/golang/%s.go.tpl", name),
		Context: ctx,
	}

	buffer := &bytes.Buffer{}

	if _, err := writer.WriteTo(buffer); err != nil {
		g.Reporter.Error("ﳑ Generating type: %s function: %s fail: %v",
			inflect.Dasherize(receiver),
			inflect.Dasherize(operation),
			err,
		)

		return
	}

	if err := root.AddFunction(buffer.String()); err != nil {
		g.Reporter.Error("ﳑ Generating type: %s function: %s fail: %v",
			inflect.Dasherize(receiver),
			inflect.Dasherize(operation),
			err,
		)

		return
	}

	g.Reporter.Success("ﳑ Generating type: %s function: %s successful",
		inflect.Dasherize(receiver),
		inflect.Dasherize(operation),
	)
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
