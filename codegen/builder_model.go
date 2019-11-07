package codegen

import (
	"strings"

	"github.com/dave/dst"
	"github.com/go-openapi/inflect"
)

// ModelBuilder builds a controller
type ModelBuilder struct {
	Package string
}

// Build builds a type
func (builder *ModelBuilder) Build(descriptor *ControllerDescriptor) []dst.Decl {
	tree := []dst.Decl{}

	// generate operations
	for _, operation := range descriptor.Operations {
		if branch := builder.input(operation); branch != nil {
			tree = append(tree, branch...)
		}

		if node := builder.output(operation); node != nil {
			tree = append(tree, node)
		}
	}

	return tree
}

func (builder *ModelBuilder) input(descriptor *OperationDescriptor) []dst.Decl {
	parent := &StructTypeBuilder{
		Name: descriptor.Name + "Input",
	}

	parent.Comments = append(parent.Comments, commentf("%s represents the input of %s operation", parent.Name, descriptor.Name))

	var (
		tree = []dst.Decl{}
	)

	if node := builder.param(parent, "path", descriptor); node != nil {
		tree = append(tree, node)
	}

	if node := builder.param(parent, "query", descriptor); node != nil {
		tree = append(tree, node)
	}

	if node := builder.param(parent, "header", descriptor); node != nil {
		tree = append(tree, node)
	}

	if node := builder.param(parent, "cookie", descriptor); node != nil {
		tree = append(tree, node)
	}

	// add the request property
	builder.request(parent, descriptor)

	if len(parent.Fields) > 0 {
		node := parent.Build()
		tree = append(tree, node)
	}

	return tree
}

func (builder *ModelBuilder) param(parent *StructTypeBuilder, kind string, descriptor *OperationDescriptor) dst.Decl {
	child := &StructTypeBuilder{
		Name: descriptor.Name + inflect.Camelize(kind) + "Input",
	}

	child.Comments = append(child.Comments,
		commentf("%s represents the %s input of %s operation", child.Name, kind, descriptor.Name))

	for _, parameter := range descriptor.Parameters {
		if strings.EqualFold(parameter.In, kind) {
			field := &Field{
				Name: parameter.Name,
				Type: parameter.ParameterType.Kind(),
				Tags: parameter.Tags(),
			}

			child.Fields = append(child.Fields, field)
		}
	}

	// no parameters no type
	if len(child.Fields) == 0 {
		return nil
	}

	field := &Field{
		Name: inflect.Camelize(kind),
		Type: pointer(child.Name),
	}

	parent.Fields = append(parent.Fields, field)
	return child.Build()
}

func (builder *ModelBuilder) request(parent *StructTypeBuilder, descriptor *OperationDescriptor) {
	for _, request := range descriptor.Requests {
		field := &Field{
			Name: "Body",
			Type: request.RequestType.Kind(),
		}

		parent.Fields = append(parent.Fields, field)
		//NOTE: for now we generate only the first request
		return
	}
}

func (builder *ModelBuilder) output(descriptor *OperationDescriptor) dst.Decl {
	child := &StructTypeBuilder{
		Name: descriptor.Name + "Output",
	}

	child.Comments = append(child.Comments, commentf("%s represents the output of %s operation", child.Name, descriptor.Name))

	for _, response := range descriptor.Responses {
		field := &Field{
			Name: "Body",
			Type: response.ResponseType.Kind(),
		}

		child.Fields = append(child.Fields, field)
		//NOTE: for now we generate only the first response
		return child.Build()
	}

	return nil
}
