package dom

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/getkin/kin-openapi/openapi3"
)

const (
	schema = "#/components/schemas/%v"
	param  = "#/components/parameters/%v"
)

// ResolverConfig is the resolver's configuration
type ResolverConfig struct {
}

// Resolver resolves all swagger spec
type Resolver struct {
	Spec *openapi3.Swagger

	Schemas map[string]*TypeDescriptor
}

// Resolve resolves the spec
func (r *Resolver) Resolve() {
	r.resolveSchemas(r.Spec.Components.Schemas)
	// r.resolveParameters(spec.Components.Parameters)
	// r.resolveHeaders(spec.Components.Headers)
	// r.resolveRequests(spec.Components.RequestBodies)
	// r.resolveResponses(spec.Components.Responses)

	// ops := r.resolveOperations(spec.Paths)
	// spew.Dump(ops)
}

func (r *Resolver) resolveSchemas(schemas map[string]*openapi3.SchemaRef) []*TypeDescriptor {
	descriptors := []*TypeDescriptor{}

	for name, ref := range schemas {
		descriptor := &TypeDescriptor{
			Path:        fmt.Sprintf(schema, name),
			Description: ref.Value.Description,
			Ref:         ref,
		}

		r.addType(descriptor)
		r.resolveType(descriptor)
		r.resolveProperties(descriptor)

		descriptors = append(descriptors, descriptor)
	}

	for _, descriptor := range descriptors {
		descriptor.Ref = nil

		for _, property := range descriptor.Properties {
			property.Ref = nil

			if property.PropertyType != nil {
				property.PropertyType.Ref = nil
			}
		}

		spew.Dump(descriptor)
	}

	return descriptors
}

func (r *Resolver) addType(descriptor *TypeDescriptor) {
	r.Schemas[descriptor.Path] = descriptor
}

func (r *Resolver) resolveType(parent *TypeDescriptor) {

}

func (r *Resolver) resolveProperties(parent *TypeDescriptor) {
	for name, ref := range parent.Ref.Value.Properties {
		descriptor := &TypeDescriptor{
			Parent:      parent,
			Path:        parent.Path + "/" + name,
			Description: ref.Value.Description,
			Ref:         ref,
		}

		r.addType(descriptor)
		r.resolveType(descriptor)
		r.resolveProperties(descriptor)

		property := &PropertyDescriptor{
			Name:          name,
			Description:   ref.Value.Description,
			Nullable:      ref.Value.Nullable,
			PropertyType:  descriptor,
			ComponentType: parent,
			Ref:           ref,
		}

		for _, field := range parent.Ref.Value.Required {
			if strings.EqualFold(name, field) {
				property.Required = true
				break
			}
		}

		parent.Properties = append(parent.Properties, property)
	}
}
