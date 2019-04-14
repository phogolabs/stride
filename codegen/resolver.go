package codegen

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/phogolabs/log"
)

// Resolver resolves all swagger spec
type Resolver struct{}

// Resolve resolves the spec
func (r *Resolver) Resolve(swagger *openapi3.Swagger) *SpecDescriptor {
	spec := &SpecDescriptor{
		Schemas:       r.resolveSchemas(swagger.Components.Schemas),
		Parameters:    r.resolveParameters(swagger.Components.Parameters),
		Headers:       r.resolveHeaders(swagger.Components.Headers),
		RequestBodies: r.resolveRequestBodies(swagger.Components.RequestBodies),
		Responses:     r.resolveResponses(swagger.Components.Responses),
		Controllers:   r.resolveControllers(swagger),
	}

	return spec
}

func (r *Resolver) resolveSchemas(schemas map[string]*openapi3.SchemaRef) TypeDescriptorCollection {
	descriptors := TypeDescriptorCollection{}

	for name, ref := range schemas {
		log.Infof("Resolving schema: %v", name)

		descriptor := r.resolveType(ref)
		descriptor.Name = name
		descriptors = append(descriptors, descriptor)
	}

	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) resolveControllers(spec *openapi3.Swagger) ControllerDescriptorCollection {
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

	controllers := ControllerDescriptorCollection{}

	for _, controller := range hashmap {
		controllers = append(controllers, controller)
	}

	sort.Sort(controllers)

	return controllers
}

func (r *Resolver) resolveOperations(schemas openapi3.Paths) OperationDescriptorCollection {
	descriptors := OperationDescriptorCollection{}

	for path, item := range schemas {
		parameters := item.Parameters

		for method, op := range item.Operations() {
			if op == nil {
				continue
			}

			log.Infof("Resolving operation: %v %v", method, path)

			params := append(parameters, op.Parameters...)

			descriptor := &OperationDescriptor{
				Path:        path,
				Method:      method,
				Name:        op.OperationID,
				Summary:     op.Summary,
				Description: op.Description,
				Deprecated:  op.Deprecated,
				Tags:        op.Tags,
				Parameters:  r.resolveOperationParameters(params),
				Responses:   r.resolveResponses(op.Responses),
			}

			if bodies := r.resolveOperationRequestBody(op.RequestBody); len(bodies) > 0 {
				descriptor.RequestBody = bodies[0]
			}

			descriptors = append(descriptors, descriptor)
		}
	}

	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) resolveOperationParameters(params openapi3.Parameters) ParameterDescriptorCollection {
	parameters := map[string]*openapi3.ParameterRef{}

	for _, param := range params {
		parameters[param.Value.Name] = param
	}

	return r.resolveParameters(parameters)
}

func (r *Resolver) resolveOperationRequestBody(body *openapi3.RequestBodyRef) RequestBodyDescriptorCollection {
	m := map[string]*openapi3.RequestBodyRef{}

	if body != nil {
		m[body.Ref] = body
	}

	return r.resolveRequestBodies(m)
}

func (r *Resolver) resolveParameters(schemas map[string]*openapi3.ParameterRef) ParameterDescriptorCollection {
	descriptors := ParameterDescriptorCollection{}

	for name, ref := range schemas {
		log.Infof("Resolving parameter: %v", name)

		descriptor := &ParameterDescriptor{
			Name:          ref.Value.Name,
			In:            ref.Value.In,
			Description:   ref.Value.Description,
			Required:      ref.Value.Required,
			Deprecated:    ref.Value.Deprecated,
			ParameterType: r.resolveType(ref.Value.Schema),
		}

		descriptors = append(descriptors, descriptor)
	}

	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) resolveHeaders(schemas map[string]*openapi3.HeaderRef) HeaderDescriptorCollection {
	descriptors := HeaderDescriptorCollection{}

	for name, ref := range schemas {
		log.Infof("Resolving header: %v", name)

		descriptor := &HeaderDescriptor{
			Name:       name,
			HeaderType: r.resolveType(ref.Value.Schema),
		}

		descriptors = append(descriptors, descriptor)
	}

	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) resolveRequestBodies(schemas map[string]*openapi3.RequestBodyRef) RequestBodyDescriptorCollection {
	descriptors := RequestBodyDescriptorCollection{}

	for name, ref := range schemas {
		log.Infof("Resolving request body: %v", name)

		descriptor := &RequestBodyDescriptor{
			Name:        name,
			Description: ref.Value.Description,
			Required:    ref.Value.Required,
			Contents:    r.resolveMediaType(ref.Value.Content),
		}

		descriptors = append(descriptors, descriptor)
	}

	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) resolveResponses(schemas map[string]*openapi3.ResponseRef) ResponseDescriptorCollection {
	descriptors := ResponseDescriptorCollection{}

	for name, ref := range schemas {
		log.Infof("Resolving response: %v", name)

		descriptor := &ResponseDescriptor{
			Name:        name,
			Description: ref.Value.Description,
			Headers:     r.resolveHeaders(ref.Value.Headers),
			Contents:    r.resolveMediaType(ref.Value.Content),
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

func (r *Resolver) resolveMediaType(content map[string]*openapi3.MediaType) ContentDescriptorCollection {
	descriptors := ContentDescriptorCollection{}

	for name, ref := range content {
		log.Infof("Resolving media type: %v", name)

		descriptor := &ContentDescriptor{
			Name:        name,
			ContentType: r.resolveType(ref.Schema),
		}

		descriptors = append(descriptors, descriptor)
	}

	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) resolveType(schemaRef *openapi3.SchemaRef) *TypeDescriptor {
	if schemaRef.Ref == "" {
		log.Infof("Resolving inline type")
	} else {
		log.Infof("Resolving referenced type: %v", schemaRef.Ref)
	}

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
			Name:        schemaRef.Ref,
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

func (r *Resolver) resolveProperties(schemaRef *openapi3.SchemaRef) PropertyDescriptorCollection {
	descriptors := PropertyDescriptorCollection{}

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

	sort.Sort(descriptors)

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
