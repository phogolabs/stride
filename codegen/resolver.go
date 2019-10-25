package codegen

import (
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/inflect"
	"github.com/phogolabs/log"
)

// ResolverState is the current resolver state
type ResolverState struct {
	Path      string
	SchemaRef *openapi3.SchemaRef
}

// Resolver resolves all swagger spec
type Resolver struct {
	TypeMapping   map[string]string
	ImportMapping map[string]string
	Cache         map[string]*TypeDescriptor
}

// Resolve resolves the spec
func (r *Resolver) Resolve(swagger *openapi3.Swagger) *SpecDescriptor {
	spec := &SpecDescriptor{
		Schemas: r.resolveSchemas(swagger.Components.Schemas),
		// Parameters:    r.resolveParameters(root, swagger.Components.Parameters),
		// Headers:       r.resolveHeaders(root, swagger.Components.Headers),
		// RequestBodies: r.resolveRequestBodies(root, swagger.Components.RequestBodies),
		// Responses:     r.resolveResponses(root, swagger.Components.Responses),
		// Operations:    r.resolveOperations(root, swagger.Paths),
	}

	return spec
}

func (r *Resolver) resolveSchemas(schemas map[string]*openapi3.SchemaRef) TypeDescriptorCollection {
	descriptors := TypeDescriptorCollection{}

	for name, schemaRef := range schemas {
		state := &ResolverState{
			Path:      name,
			SchemaRef: schemaRef,
		}

		descriptors = append(descriptors, r.resolveType(state))
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
			// log.Infof("Resolving operation: %v", path)

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
	parameters := map[string]*openapi3.ParameterRef{}

	for _, param := range params {
		parameters[param.Value.Name] = param
	}

	return r.resolveParameters(parent, parameters)
}

func (r *Resolver) resolveParameters(parent string, schemas map[string]*openapi3.ParameterRef) ParameterDescriptorCollection {
	descriptors := ParameterDescriptorCollection{}

	// for name, ref := range schemas {
	// 	path := join(parent, "parameters", name)
	// 	log.Infof("Resolving parameter: %v", path)

	// 	descriptor := &ParameterDescriptor{
	// 		Name:          ref.Value.Name,
	// 		In:            ref.Value.In,
	// 		Description:   ref.Value.Description,
	// 		Required:      ref.Value.Required,
	// 		Deprecated:    ref.Value.Deprecated,
	// 		ParameterType: r.resolveType(join(path, "schema"), ref.Value.Schema),
	// 	}

	// 	descriptors = append(descriptors, descriptor)
	// }

	// sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) resolveHeaders(parent string, schemas map[string]*openapi3.HeaderRef) HeaderDescriptorCollection {
	descriptors := HeaderDescriptorCollection{}

	// for name, ref := range schemas {
	// 	path := join(parent, "headers", name)
	// 	log.Infof("Resolving header: %v", path)

	// 	descriptor := &HeaderDescriptor{
	// 		Name:       name,
	// 		HeaderType: r.resolveType(join(path, "schema"), ref.Value.Schema),
	// 	}

	// 	descriptors = append(descriptors, descriptor)
	// }

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
	var (
		err         error
		descriptors = ResponseDescriptorCollection{}
	)

	for name, ref := range schemas {
		code := 0

		if code, err = strconv.Atoi(name); err == nil {
			name = strings.ToLower(http.StatusText(code))
		}

		name = inflect.Dasherize(name + "-response")

		path := join(parent, "responses", name)
		// log.Infof("Resolving response: %v", path)

		descriptor := &ResponseDescriptor{
			Name:        name,
			Description: ref.Value.Description,
			Headers:     r.resolveHeaders(path, ref.Value.Headers),
			Contents:    r.resolveMediaType(path, ref.Value.Content),
		}

		descriptors = append(descriptors, descriptor)
	}

	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) resolveMediaType(parent string, content map[string]*openapi3.MediaType) ContentDescriptorCollection {
	descriptors := ContentDescriptorCollection{}

	// for name, ref := range content {
	// 	path := join(parent, "contents", name)
	// 	log.Infof("Resolving media type: %v", path)

	// 	descriptor := &ContentDescriptor{
	// 		Name:        name,
	// 		ContentType: r.resolveType(join(path, "schema"), ref.Schema),
	// 	}

	// 	descriptors = append(descriptors, descriptor)
	// }

	// sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) resolveType(state *ResolverState) *TypeDescriptor {
	// log.Infof("Resolving path: %v", state.Path)

	// if descriptor := r.resolveObjectType(state); descriptor != nil {
	// r.Cache[key] = descriptor
	// return descriptor
	// }

	if descriptor := r.resolveEnumType(state); descriptor != nil {
		// r.Cache[key] = descriptor
		return descriptor
	}

	if descriptor := r.resolvePrimitiveType(state); descriptor != nil {
		// 	// r.Cache[key] = descriptor
		return descriptor
	}

	return nil
}

func (r *Resolver) resolvePrimitiveType(state *ResolverState) *TypeDescriptor {
	declaration := NewTypeDeclarationPrimitive(state.SchemaRef)

	if declaration == nil {
		return nil
	}

	return NewTypeDescriptorPrimitive(declaration)
}

func (r *Resolver) resolveEnumType(state *ResolverState) *TypeDescriptor {
	if len(state.SchemaRef.Value.Enum) == 0 {
		return nil
	}

	if state.SchemaRef.Value.Type != "string" {
		//NOTE: we support only string enums for now
		//TODO: write to the standard log
		return nil
	}

	declaration := &TypeDeclaration{
		Name:      state.Path,
		SchemaRef: state.SchemaRef,
	}

	return NewTypeDescriptorEnum(declaration)
}

func (r *Resolver) resolveObjectType(parent string, schemaRef *openapi3.SchemaRef) *TypeDescriptor {
	switch schemaRef.Value.Type {
	case "object":
		return &TypeDescriptor{
			IsClass:     true,
			Name:        r.name(parent, schemaRef),
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
	descriptors := PropertyDescriptorCollection{}

	// for name, ref := range schemaRef.Value.Properties {
	// 	path := join(parent, "properties", name)

	// 	property := &PropertyDescriptor{
	// 		Name:         name,
	// 		Description:  ref.Value.Description,
	// 		Nullable:     ref.Value.Nullable,
	// 		PropertyType: r.resolveType(join(path, "schema"), ref),
	// 	}

	// 	for _, field := range schemaRef.Value.Required {
	// 		if strings.EqualFold(name, field) {
	// 			property.Required = true
	// 			break
	// 		}
	// 	}

	// 	descriptors = append(descriptors, property)
	// }

	// sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) name(parent string, schemaRef *openapi3.SchemaRef) string {
	key := schemaRef.Ref

	if key == "" {
		key = parent
	}

	key = strings.ReplaceAll(key, "components", "")
	key = strings.ReplaceAll(key, "schemas", "")
	key = strings.ReplaceAll(key, "schema", "")
	key = strings.ReplaceAll(key, "operations", "")
	key = strings.ReplaceAll(key, "headers", "")
	key = strings.ReplaceAll(key, "responses", "")
	key = strings.ReplaceAll(key, "request-bodies", "")
	key = strings.ReplaceAll(key, "contents", "")
	key = strings.ReplaceAll(key, "application", "")
	key = strings.ReplaceAll(key, "#", "")

	var parts []string

	for _, part := range strings.Split(key, "/") {
		part = strings.TrimSpace(part)

		if part == "" {
			continue
		}

		parts = append(parts, part)
	}

	key = strings.Join(parts, "-")
	return key
}

func camelize(key string) string {
	key = inflect.Camelize(key)
	key = strings.ReplaceAll(key, "Xml", "XML")
	key = strings.ReplaceAll(key, "Json", "JSON")
	key = strings.ReplaceAll(key, "Id", "ID")
	key = strings.ReplaceAll(key, "Ok", "OK")
	return key
}

func join(paths ...string) string {
	for index, part := range paths {
		part = strings.ToLower(inflect.Dasherize(part))
		paths[index] = part
	}

	return filepath.Join(paths...)
}
