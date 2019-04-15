package codegen

import (
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/inflect"
	"github.com/phogolabs/log"
)

// Resolver resolves all swagger spec
type Resolver struct{}

// Resolve resolves the spec
func (r *Resolver) Resolve(swagger *openapi3.Swagger) *SpecDescriptor {
	const root = "#/components"

	spec := &SpecDescriptor{
		Schemas:       r.resolveSchemas(root, swagger.Components.Schemas),
		Parameters:    r.resolveParameters(root, swagger.Components.Parameters),
		Headers:       r.resolveHeaders(root, swagger.Components.Headers),
		RequestBodies: r.resolveRequestBodies(root, swagger.Components.RequestBodies),
		Responses:     r.resolveResponses(root, swagger.Components.Responses),
		Operations:    r.resolveOperations(root, swagger.Paths),
	}

	return spec
}

func (r *Resolver) resolveSchemas(parent string, schemas map[string]*openapi3.SchemaRef) TypeDescriptorCollection {
	descriptors := TypeDescriptorCollection{}

	for name, ref := range schemas {
		path := join(parent, "schemas", name)
		log.Infof("Resolving schema: %v", path)

		descriptor := r.resolveType(path, ref)
		descriptor.Name = name
		descriptors = append(descriptors, descriptor)
	}

	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) resolveOperations(parent string, schemas openapi3.Paths) OperationDescriptorCollection {
	descriptors := OperationDescriptorCollection{}

	for uri, item := range schemas {
		parameters := item.Parameters

		for method, op := range item.Operations() {
			if op == nil {
				continue
			}

			params := append(parameters, op.Parameters...)

			path := join(parent, "operations", op.OperationID)
			log.Infof("Resolving operation: %v", path)

			descriptor := &OperationDescriptor{
				Path:        uri,
				Method:      method,
				Name:        op.OperationID,
				Summary:     op.Summary,
				Description: op.Description,
				Deprecated:  op.Deprecated,
				Tags:        op.Tags,
				Parameters:  r.resolveOperationParameters(path, params),
				Responses:   r.resolveResponses(path, op.Responses),
				// RequestBody: r.resolveRequestBody(join(path, "request-bodies", "default"), op.RequestBody),
			}

			descriptors = append(descriptors, descriptor)
		}
	}

	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) resolveOperationParameters(parent string, params openapi3.Parameters) ParameterDescriptorCollection {
	// PATH: "#/components/operations/op.OperationID/parameters/{parameter_name}/...
	parameters := map[string]*openapi3.ParameterRef{}

	for _, param := range params {
		parameters[param.Value.Name] = param
	}

	return r.resolveParameters(parent, parameters)
}

func (r *Resolver) resolveParameters(parent string, schemas map[string]*openapi3.ParameterRef) ParameterDescriptorCollection {
	descriptors := ParameterDescriptorCollection{}

	for name, ref := range schemas {
		path := join(parent, "parameters", ref.Value.Name)
		log.Infof("Resolving parameter: %v (%v)", path, name)

		descriptor := &ParameterDescriptor{
			Name:          ref.Value.Name,
			In:            ref.Value.In,
			Description:   ref.Value.Description,
			Required:      ref.Value.Required,
			Deprecated:    ref.Value.Deprecated,
			ParameterType: r.resolveType(join(path, "schema"), ref.Value.Schema),
		}

		descriptors = append(descriptors, descriptor)
	}

	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) resolveHeaders(parent string, schemas map[string]*openapi3.HeaderRef) HeaderDescriptorCollection {
	descriptors := HeaderDescriptorCollection{}

	for name, ref := range schemas {
		path := join(parent, "headers", name)
		log.Infof("Resolving header: %v", path)

		descriptor := &HeaderDescriptor{
			Name:       name,
			HeaderType: r.resolveType(join(path, "schema"), ref.Value.Schema),
		}

		descriptors = append(descriptors, descriptor)
	}

	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) resolveRequestBodies(parent string, schemas map[string]*openapi3.RequestBodyRef) RequestBodyDescriptorCollection {
	descriptors := RequestBodyDescriptorCollection{}

	for name, ref := range schemas {
		path := join(parent, "request-bodies", name)
		log.Infof("Resolving request body: %v", path)

		descriptor := &RequestBodyDescriptor{
			Name:        name,
			Description: ref.Value.Description,
			Required:    ref.Value.Required,
			Contents:    r.resolveMediaType(path, ref.Value.Content),
		}

		descriptors = append(descriptors, descriptor)
	}

	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) resolveResponses(parent string, schemas map[string]*openapi3.ResponseRef) ResponseDescriptorCollection {
	descriptors := ResponseDescriptorCollection{}

	for name, ref := range schemas {
		path := join(parent, "responses", name)
		log.Infof("Resolving response: %v", path)

		descriptor := &ResponseDescriptor{
			Name:        name,
			Description: ref.Value.Description,
			Headers:     r.resolveHeaders(path, ref.Value.Headers),
			Contents:    r.resolveMediaType(path, ref.Value.Content),
		}

		if code, err := strconv.Atoi(name); err == nil {
			descriptor.Code = code
			descriptor.Name = http.StatusText(code)
		}

		descriptors = append(descriptors, descriptor)
	}

	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) resolveMediaType(parent string, content map[string]*openapi3.MediaType) ContentDescriptorCollection {
	descriptors := ContentDescriptorCollection{}

	for name, ref := range content {
		path := join(parent, "contents", name)
		log.Infof("Resolving media type: %v", path)

		descriptor := &ContentDescriptor{
			Name:        name,
			ContentType: r.resolveType(join(path, "schema"), ref.Schema),
		}

		descriptors = append(descriptors, descriptor)
	}

	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) resolveType(parent string, schemaRef *openapi3.SchemaRef) *TypeDescriptor {
	log.Infof("Resolving type: %v", parent)

	var descriptor *TypeDescriptor

	if descriptor = r.resolveObjectType(parent, schemaRef); descriptor != nil {
		return descriptor
	}

	if descriptor = r.resolveEnumType(parent, schemaRef); descriptor != nil {
		return descriptor
	}

	if descriptor = r.resolvePrimitiveType(parent, schemaRef); descriptor != nil {
		return descriptor
	}

	return nil
}

func (r *Resolver) resolveObjectType(parent string, schemaRef *openapi3.SchemaRef) *TypeDescriptor {
	switch schemaRef.Value.Type {
	case "object":
		return &TypeDescriptor{
			IsClass:     true,
			Name:        schemaRef.Ref,
			Properties:  r.resolveProperties(parent, schemaRef),
			Description: schemaRef.Value.Description,
		}
	case "array":
		if descriptor := r.resolveObjectType(parent, schemaRef.Value.Items); descriptor != nil {
			descriptor = descriptor.Clone()
			descriptor.IsArray = true
			return descriptor
		}
	}

	return nil
}

func (r *Resolver) resolveProperties(parent string, schemaRef *openapi3.SchemaRef) PropertyDescriptorCollection {
	// PATH: #/components/schemas/{schema-name}/properties/...
	// PATH: #/components/schemas/{schema-name}/properties/{property-name}/schema/properties/...

	descriptors := PropertyDescriptorCollection{}

	for name, ref := range schemaRef.Value.Properties {
		path := join(parent, "properties", name, "schema")

		property := &PropertyDescriptor{
			Name:         name,
			Description:  ref.Value.Description,
			Nullable:     ref.Value.Nullable,
			PropertyType: r.resolveType(path, ref),
		}

		for _, field := range schemaRef.Value.Required {
			if strings.EqualFold(name, field) {
				property.Required = true
				break
			}
		}

		descriptors = append(descriptors, property)
	}

	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) resolvePrimitiveType(parent string, schemaRef *openapi3.SchemaRef) *TypeDescriptor {
	switch schemaRef.Value.Format {
	case "int32":
		return &TypeDescriptor{
			Name:        "int32",
			IsPrimitive: true,
		}
	case "int64":
		return &TypeDescriptor{
			Name:        "int64",
			IsPrimitive: true,
		}
	case "float":
		return &TypeDescriptor{
			Name:        "float",
			IsPrimitive: true,
		}
	case "double":
		return &TypeDescriptor{
			Name:        "double",
			IsPrimitive: true,
		}
	case "byte":
		return &TypeDescriptor{
			Name:        "byte",
			IsPrimitive: true,
		}
	case "uuid":
		return &TypeDescriptor{
			Name:        "uuid",
			IsPrimitive: true,
		}
	case "binary":
		return &TypeDescriptor{
			Name:        "binary",
			IsPrimitive: true,
		}
	case "date":
		return &TypeDescriptor{
			Name:        "date-time",
			IsPrimitive: true,
		}
	case "date-time":
		return &TypeDescriptor{
			Name:        "date-time",
			IsPrimitive: true,
		}
	}

	switch schemaRef.Value.Type {
	case "string":
		return &TypeDescriptor{
			Name:        "string",
			IsPrimitive: true,
		}
	case "number":
		return &TypeDescriptor{
			Name:        "double",
			IsPrimitive: true,
		}
	case "integer":
		return &TypeDescriptor{
			Name:        "int64",
			IsPrimitive: true,
		}
	case "array":
		if descriptor := r.resolvePrimitiveType(parent, schemaRef.Value.Items); descriptor != nil {
			descriptor = descriptor.Clone()
			descriptor.IsArray = true
			return descriptor
		}
	}

	return nil
}

func (r *Resolver) resolveEnumType(parent string, schemaRef *openapi3.SchemaRef) *TypeDescriptor {
	if len(schemaRef.Value.Enum) > 0 {
		descriptor := &TypeDescriptor{
			IsEnum:      true,
			Name:        schemaRef.Ref,
			Description: schemaRef.Value.Description,
		}

		for _, value := range schemaRef.Value.Enum {
			property := &PropertyDescriptor{
				Name: fmt.Sprintf("%v", value),
			}

			descriptor.Properties = append(descriptor.Properties, property)
		}

		return descriptor
	}

	return nil
}

func join(paths ...string) string {
	for index, part := range paths {
		part = strings.ToLower(inflect.Dasherize(part))
		paths[index] = part
	}

	return filepath.Join(paths...)
}
