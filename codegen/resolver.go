package codegen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/getkin/kin-openapi/openapi3"
)

// Resolver resolves all swagger spec
type Resolver struct {
	Cache map[string]*TypeDescriptor
}

// Resolve resolves the spec
func (r *Resolver) Resolve(swagger *openapi3.Swagger) *Spec {
	spec := &Spec{
		Schemas:       r.resolveSchemas(swagger.Components.Schemas),
		Parameters:    r.resolveParameters(swagger.Components.Parameters),
		Headers:       r.resolveHeaders(swagger.Components.Headers),
		RequestBodies: r.resolveRequestBodies(swagger.Components.RequestBodies),
		Responses:     r.resolveResponses(swagger.Components.Responses),
		Controllers:   r.resolveControllers(swagger),
	}

	spew.Dump(spec)
	return spec
}

func (r *Resolver) resolveSchemas(schemas map[string]*openapi3.SchemaRef) []*TypeDescriptor {
	descriptors := []*TypeDescriptor{}

	for name, ref := range schemas {
		descriptor := r.resolveType(ref)
		descriptor.Name = name
		descriptors = append(descriptors, descriptor)

		// spew.Dump(descriptor)
	}

	return descriptors
}

func (r *Resolver) resolveControllers(spec *openapi3.Swagger) []*ControllerDescriptor {
	var (
		hashmap    = make(map[string]*ControllerDescriptor)
		operations = r.resolveOperations(spec.Paths)
	)

	for _, operation := range operations {
		tag := "default"

		if len(operation.Tags) > 0 {
			tag = operation.Tags[0]
		}

		controller, ok := hashmap[tag]
		if !ok {
			controller = &ControllerDescriptor{
				Name: tag,
			}

			if descriptor := spec.Components.Tags.Get(tag); descriptor != nil {
				controller.Description = descriptor.Description
			}

			hashmap[tag] = controller
		}

		controller.Operations = append(controller.Operations, operation)
	}

	controllers := []*ControllerDescriptor{}

	for _, controller := range hashmap {
		controllers = append(controllers, controller)
	}

	return controllers
}

func (r *Resolver) resolveOperations(schemas openapi3.Paths) []*OperationDescriptor {
	descriptors := []*OperationDescriptor{}

	for path, item := range schemas {
		kv := map[string]*openapi3.Operation{
			"CONNECT": item.Connect,
			"DELETE":  item.Delete,
			"GET":     item.Get,
			"HEAD":    item.Head,
			"OPTIONS": item.Options,
			"PATCH":   item.Patch,
			"POST":    item.Post,
			"PUT":     item.Put,
			"TRACE":   item.Trace,
		}

		for method, op := range kv {
			if op == nil {
				continue
			}

			descriptor := &OperationDescriptor{
				Path:          path,
				Method:        method,
				Name:          op.OperationID,
				Summary:       op.Summary,
				Description:   op.Description,
				Deprecated:    op.Deprecated,
				Tags:          op.Tags,
				Parameters:    r.resolveOperationParameters(item.Parameters),
				Responses:     r.resolveResponses(op.Responses),
				RequestBodies: r.resolveOperationRequestBody(op.RequestBody),
			}

			descriptors = append(descriptors, descriptor)

			// spew.Dump(descriptor)
		}
	}

	return descriptors
}

func (r *Resolver) resolveOperationParameters(params []*openapi3.ParameterRef) []*ParameterDescriptor {
	parameters := map[string]*openapi3.ParameterRef{}

	for _, param := range params {
		parameters[param.Value.Name] = param
	}

	return r.resolveParameters(parameters)
}

func (r *Resolver) resolveOperationRequestBody(body *openapi3.RequestBodyRef) []*RequestBodyDescriptor {
	m := map[string]*openapi3.RequestBodyRef{}

	if body != nil {
		m["body"] = body
	}

	return r.resolveRequestBodies(m)
}

func (r *Resolver) resolveParameters(schemas map[string]*openapi3.ParameterRef) []*ParameterDescriptor {
	descriptors := []*ParameterDescriptor{}

	for name, ref := range schemas {
		descriptor := &ParameterDescriptor{
			Name:          name,
			In:            ref.Value.In,
			Description:   ref.Value.Description,
			Required:      ref.Value.Required,
			Deprecated:    ref.Value.Deprecated,
			ParameterType: r.resolveType(ref.Value.Schema),
		}

		descriptors = append(descriptors, descriptor)
	}

	return descriptors
}

func (r *Resolver) resolveHeaders(schemas map[string]*openapi3.HeaderRef) []*HeaderDescriptor {
	descriptors := []*HeaderDescriptor{}

	for name, ref := range schemas {
		descriptor := &HeaderDescriptor{
			Name:       name,
			HeaderType: r.resolveType(ref.Value.Schema),
		}

		descriptors = append(descriptors, descriptor)

		// spew.Dump(descriptor)
	}

	return descriptors
}

func (r *Resolver) resolveRequestBodies(schemas map[string]*openapi3.RequestBodyRef) []*RequestBodyDescriptor {
	descriptors := []*RequestBodyDescriptor{}

	for name, ref := range schemas {
		descriptor := &RequestBodyDescriptor{
			Name:        name,
			Description: ref.Value.Description,
			Required:    ref.Value.Required,
			Contents:    r.resolveMediaType(ref.Value.Content),
		}

		descriptors = append(descriptors, descriptor)

		// spew.Dump(descriptor)
	}

	return descriptors
}

func (r *Resolver) resolveResponses(schemas map[string]*openapi3.ResponseRef) []*ResponseDescriptor {
	descriptors := []*ResponseDescriptor{}

	for name, ref := range schemas {
		descriptor := &ResponseDescriptor{
			Description: ref.Value.Description,
			Headers:     r.resolveHeaders(ref.Value.Headers),
			Contents:    r.resolveMediaType(ref.Value.Content),
		}

		if code, err := strconv.Atoi(name); err == nil {
			descriptor.Code = code
		}

		descriptors = append(descriptors, descriptor)

		// spew.Dump(descriptor)
	}

	return descriptors
}

func (r *Resolver) resolveMediaType(content map[string]*openapi3.MediaType) []*ContentDescriptor {
	descriptors := []*ContentDescriptor{}

	for name, ref := range content {
		descriptor := &ContentDescriptor{
			Name:        name,
			ContentType: r.resolveType(ref.Schema),
		}

		descriptors = append(descriptors, descriptor)
	}

	return descriptors
}

func (r *Resolver) resolveType(schemaRef *openapi3.SchemaRef) *TypeDescriptor {
	var descriptor *TypeDescriptor

	if descriptor = r.lookup(schemaRef); descriptor != nil {
		return descriptor
	}

	if descriptor = r.resolveObjectType(schemaRef); descriptor != nil {
		r.store(schemaRef, descriptor)
		return descriptor
	}

	if descriptor = r.resolveEnumType(schemaRef); descriptor != nil {
		r.store(schemaRef, descriptor)
		return descriptor
	}

	if descriptor = r.resolvePrimitiveType(schemaRef); descriptor != nil {
		return descriptor
	}

	return nil
}

func (r *Resolver) lookup(schemaRef *openapi3.SchemaRef) *TypeDescriptor {
	if descriptor, ok := r.Cache[schemaRef.Ref]; ok {
		return descriptor
	}

	return nil
}

func (r *Resolver) store(schemaRef *openapi3.SchemaRef, descriptor *TypeDescriptor) {
	if key := schemaRef.Ref; key != "" {
		r.Cache[key] = descriptor
	}
}

func (r *Resolver) resolveObjectType(schemaRef *openapi3.SchemaRef) *TypeDescriptor {
	switch schemaRef.Value.Type {
	case "object":
		return &TypeDescriptor{
			IsClass:     true,
			Properties:  r.resolveProperties(schemaRef),
			Description: schemaRef.Value.Description,
		}
	case "array":
		if descriptor := r.resolveObjectType(schemaRef.Value.Items); descriptor != nil {
			descriptor = descriptor.Clone()
			descriptor.IsArray = true
			return descriptor
		}
	}

	return nil
}

func (r *Resolver) resolveProperties(schemaRef *openapi3.SchemaRef) []*PropertyDescriptor {
	descriptors := []*PropertyDescriptor{}

	for name, ref := range schemaRef.Value.Properties {
		property := &PropertyDescriptor{
			Name:         name,
			Description:  ref.Value.Description,
			Nullable:     ref.Value.Nullable,
			PropertyType: r.resolveType(ref),
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
	case "array":
		if descriptor := r.resolvePrimitiveType(schemaRef.Value.Items); descriptor != nil {
			descriptor = descriptor.Clone()
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
