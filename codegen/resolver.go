package codegen

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

// Resolver resolves all swagger spec
type Resolver struct {
	Schemas map[string]*TypeDescriptor
}

// Resolve resolves the spec
func (r *Resolver) Resolve(spec *openapi3.Swagger) {
	r.resolveSchemas(spec.Components.Schemas)
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
		ref.Ref = fmt.Sprintf(schema, name)

		descriptor := r.resolveType(ref)
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

func (r *Resolver) resolveType(schemaRef *openapi3.SchemaRef) *TypeDescriptor {
	var descriptor *TypeDescriptor

	if descriptor = r.resolveObjectType(schemaRef); descriptor != nil {
		return descriptor
	}

	if descriptor = r.resolveEnumType(schemaRef); descriptor != nil {
		return descriptor
	}

	if descriptor = r.resolvePrimitiveType(schemaRef); descriptor != nil {
		return descriptor
	}

	return nil
}

func (r *Resolver) resolveObjectType(schemaRef *openapi3.SchemaRef) *TypeDescriptor {
	switch schemaRef.Value.Type {
	case "object":
		return &TypeDescriptor{
			IsClass:     true,
			Properties:  r.resolveProperties(schemaRef),
			Path:        schemaRef.Ref,
			Description: schemaRef.Value.Description,
			Ref:         schemaRef,
		}
	case "array":
		if descriptor := r.resolveObjectType(schemaRef.Value.Items); descriptor != nil {
			descriptor.IsArray = true
			return descriptor
		}
	}

	return nil
}

func (r *Resolver) resolveProperties(schemaRef *openapi3.SchemaRef) []*PropertyDescriptor {
	descriptors := []*PropertyDescriptor{}

	for name, ref := range schemaRef.Value.Properties {
		if ref.Ref == "" {
			ref.Ref = schemaRef.Ref + "/" + name
		}

		property := &PropertyDescriptor{
			Name:         name,
			Description:  ref.Value.Description,
			Nullable:     ref.Value.Nullable,
			PropertyType: r.resolveType(ref),
			// ComponentType: schemaref,
			Ref: ref,
		}

		for _, field := range schemaRef.Value.Required {
			if strings.EqualFold(name, field) {
				property.Required = true
				break
			}
		}

		descriptors = append(descriptors, property)
	}

	return descriptors
}

func (r *Resolver) resolvePrimitiveType(schemaRef *openapi3.SchemaRef) *TypeDescriptor {
	switch schemaRef.Value.Format {
	case "int32":
		return &TypeDescriptor{
			Name:        "int32",
			IsPrimitive: true,
			Ref:         schemaRef,
		}
	case "int64":
		return &TypeDescriptor{
			Name:        "int64",
			IsPrimitive: true,
			Ref:         schemaRef,
		}
	case "float":
		return &TypeDescriptor{
			Name:        "float",
			IsPrimitive: true,
			Ref:         schemaRef,
		}
	case "double":
		return &TypeDescriptor{
			Name:        "double",
			IsPrimitive: true,
			Ref:         schemaRef,
		}
	case "byte":
		return &TypeDescriptor{
			Name:        "byte",
			IsPrimitive: true,
			Ref:         schemaRef,
		}
	case "binary":
		return &TypeDescriptor{
			Name:        "binary",
			IsPrimitive: true,
			Ref:         schemaRef,
		}
	case "date":
		return &TypeDescriptor{
			Name:        "date-time",
			IsPrimitive: true,
			Ref:         schemaRef,
		}
	case "date-time":
		return &TypeDescriptor{
			Name:        "date-time",
			IsPrimitive: true,
			Ref:         schemaRef,
		}
	}

	switch schemaRef.Value.Type {
	case "string":
		return &TypeDescriptor{
			Name:        "string",
			IsPrimitive: true,
			Ref:         schemaRef,
		}
	case "array":
		if descriptor := r.resolvePrimitiveType(schemaRef.Value.Items); descriptor != nil {
			descriptor.IsArray = true
			return descriptor
		}
	}

	return nil
}

func (r *Resolver) resolveEnumType(schemaRef *openapi3.SchemaRef) *TypeDescriptor {
	if len(schemaRef.Value.Enum) > 0 {
		descriptor := &TypeDescriptor{
			IsEnum:      true,
			Path:        schemaRef.Ref,
			Description: schemaRef.Value.Description,
			Ref:         schemaRef,
		}

		for _, value := range schemaRef.Value.Enum {
			property := &PropertyDescriptor{
				Name:          fmt.Sprintf("%v", value),
				ComponentType: descriptor,
			}

			descriptor.Properties = append(descriptor.Properties, property)
		}

		return descriptor
	}

	return nil
}
