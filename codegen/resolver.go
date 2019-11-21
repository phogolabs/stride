package codegen

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/inflect"
)

// Resolver resolves all swagger spec
type Resolver struct {
	cache TypeDescriptorMap
}

// NewResolver creates a new resolver
func NewResolver() *Resolver {
	return &Resolver{
		cache: TypeDescriptorMap{},
	}
}

// Resolve resolves the spec
func (r *Resolver) Resolve(swagger *openapi3.Swagger) *SpecDescriptor {
	defer r.cache.Clear()

	var (
		components  = swagger.Components
		ctx         = emptyCtx
		controllers = r.operations(ctx, swagger.Paths)
	)

	r.schemas(ctx, components.Schemas)
	r.parameters(ctx, components.Parameters)
	r.headers(ctx, components.Headers)
	r.requests(ctx, components.RequestBodies)
	r.responses(ctx, components.Responses)

	return &SpecDescriptor{
		Types:       r.cache.Collection(),
		Controllers: controllers,
	}
}

func (r *Resolver) schemas(ctx *ResolverContext, schemas map[string]*openapi3.SchemaRef) TypeDescriptorCollection {
	descriptors := TypeDescriptorCollection{}

	for name, schema := range schemas {
		cctx := ctx.Child(name, schema)
		descriptors = append(descriptors, r.resolve(cctx))
	}

	return descriptors
}

func (r *Resolver) operations(ctx *ResolverContext, operations map[string]*openapi3.PathItem) ControllerDescriptorCollection {
	descriptors := ControllerDescriptorMap{}

	key := func(tags []string) string {
		key := "default"

		if len(tags) > 0 {
			key = tags[0]
		}

		return key
	}

	for path, spec := range operations {
		for method, spec := range spec.Operations() {
			var (
				controller = descriptors.Get(key(spec.Tags))
				cctx       = ctx.Child(spec.OperationID, nil)
			)

			var (
				parameters = make(map[string]*openapi3.ParameterRef)
				requests   = make(map[string]*openapi3.RequestBodyRef)
				responses  = spec.Responses
			)

			requests["request"] = spec.RequestBody

			for _, param := range spec.Parameters {
				parameters[param.Value.Name] = param
			}

			operation := &OperationDescriptor{
				Path:        path,
				Method:      method,
				Name:        spec.OperationID,
				Description: spec.Description,
				Summary:     spec.Summary,
				Deprecated:  spec.Deprecated,
				Tags:        spec.Tags,
				Requests:    r.requests(cctx, requests, r.parameters(cctx, parameters)...),
				Responses:   r.responses(cctx, responses),
			}

			controller.Operations = append(controller.Operations, operation)
		}
	}

	return descriptors.Collection()
}

func (r *Resolver) requests(ctx *ResolverContext, bodies map[string]*openapi3.RequestBodyRef, params ...*ParameterDescriptor) RequestDescriptorCollection {
	descriptors := RequestDescriptorCollection{}

	for name, spec := range bodies {
		if spec == nil {
			continue
		}

		for contentType, content := range spec.Value.Content {
			if !strings.EqualFold(contentType, "application/json") {
				//TODO: at some point we must support all content-types
				continue
			}

			schema := content.Schema

			if reference := spec.Ref; reference != "" {
				if spec.Ref != "" {
					schema = &openapi3.SchemaRef{
						Ref:   reference,
						Value: schema.Value,
					}
				}
			}

			var (
				cctx    = ctx.Child(name, schema)
				request = &RequestDescriptor{
					ContentType: contentType,
					Description: spec.Value.Description,
					Required:    spec.Value.Required,
					RequestType: r.resolve(cctx),
					Parameters:  params,
				}
			)

			descriptors = append(descriptors, request)
		}
	}

	// sort descriptors
	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) responses(ctx *ResolverContext, responses map[string]*openapi3.ResponseRef) ResponseDescriptorCollection {
	descriptors := ResponseDescriptorCollection{}

	for name, spec := range responses {
		code, err := strconv.Atoi(name)
		if err == nil {
			name = inflect.Dasherize(http.StatusText(code)) + "-response"
		}

		for contentType, content := range spec.Value.Content {
			if !strings.EqualFold(contentType, "application/json") {
				//TODO: at some point we must support all content-types
				continue
			}

			schema := content.Schema

			if reference := spec.Ref; reference != "" {
				if spec.Ref != "" {
					schema = &openapi3.SchemaRef{
						Ref:   reference,
						Value: schema.Value,
					}
				}
			}

			var (
				cctx     = ctx.Child(name, schema)
				response = &ResponseDescriptor{
					Code:         code,
					ContentType:  contentType,
					Description:  spec.Value.Description,
					ResponseType: r.resolve(cctx),
					Parameters:   r.headers(cctx, spec.Value.Headers),
				}
			)

			descriptors = append(descriptors, response)
		}
	}

	// sort descriptors
	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) parameters(ctx *ResolverContext, parameters map[string]*openapi3.ParameterRef) ParameterDescriptorCollection {
	descriptors := ParameterDescriptorCollection{}

	for name, spec := range parameters {
		var (
			cctx      = ctx.Child(name, spec.Value.Schema)
			parameter = &ParameterDescriptor{
				Name:          spec.Value.Name,
				In:            spec.Value.In,
				Style:         spec.Value.Style,
				Description:   spec.Value.Description,
				Required:      spec.Value.Required,
				Deprecated:    spec.Value.Deprecated,
				ParameterType: r.resolve(cctx),
			}
		)

		//TODO: handle the explode value propertly as per spec
		if value := spec.Value.Explode; value != nil {
			parameter.Explode = *value
		}

		descriptors = append(descriptors, parameter)
	}

	// sort parameters
	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) headers(ctx *ResolverContext, headers map[string]*openapi3.HeaderRef) ParameterDescriptorCollection {
	descriptors := ParameterDescriptorCollection{}

	for name, spec := range headers {
		var (
			cctx   = ctx.Child(name, spec.Value.Schema)
			header = &ParameterDescriptor{
				Name:          name,
				In:            "header",
				ParameterType: r.resolve(cctx),
			}
		)

		descriptors = append(descriptors, header)
	}

	// sort descriptors
	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) resolve(ctx *ResolverContext) *TypeDescriptor {
	if descriptor := r.cache.Get(ctx.Name); descriptor != nil {
		return descriptor
	}

	// reference type descriptor
	if reference := ctx.Schema.Ref; reference != "" {
		descriptor := r.resolve(ctx.Dereference())

		if ctx.Parent.IsRoot() {
			descriptor = &TypeDescriptor{
				Name:        inflect.Dasherize(ctx.Name),
				Description: ctx.Schema.Value.Description,
				IsAlias:     true,
				Element:     descriptor,
			}

			// add the descriptor to the cache
			r.cache.Add(descriptor)
		}

		return descriptor
	}

	// class type descriptor
	if kind := r.kind(ctx.Schema.Value); kind == "object" {
		descriptor := &TypeDescriptor{
			Name:        inflect.Dasherize(ctx.Name),
			Description: ctx.Schema.Value.Description,
			IsClass:     true,
			IsNullable:  true,
		}

		//TODO: handle min and max properites somehow
		//TODO: handle pattern properties
		//TODO: handle discriminator

		required := func(name string) bool {
			for _, key := range ctx.Schema.Value.Required {
				if strings.EqualFold(key, name) {
					return true
				}
			}

			return false
		}

		for field, schema := range ctx.Schema.Value.Properties {
			property := &PropertyDescriptor{
				Name:         field,
				Description:  schema.Value.Description,
				ReadOnly:     schema.Value.ReadOnly,
				WriteOnly:    schema.Value.WriteOnly,
				Required:     required(field),
				PropertyType: r.resolve(ctx.Child(field, schema)),
			}

			descriptor.Properties = append(descriptor.Properties, property)
		}

		if extra := ctx.Schema.Value.AdditionalProperties; extra != nil {
			//TODO: handle additonal properties
		}

		// sort properties by name
		sort.Sort(descriptor.Properties)

		// add the descriptor to the cache
		r.cache.Add(descriptor)

		return descriptor
	}

	// array descriptor
	if kind := r.kind(ctx.Schema.Value); kind == "array" {
		descriptor := &TypeDescriptor{
			Name:        inflect.Dasherize(ctx.Name),
			Description: ctx.Schema.Value.Description,
			Default:     ctx.Schema.Value.Default,
			IsNullable:  ctx.Schema.Value.Nullable,
			IsArray:     true,
			Element:     r.resolve(ctx.Array()),
			Metadata: Metadata{
				"unique": ctx.Schema.Value.UniqueItems,
				"min":    ctx.Schema.Value.MinLength,
				"max":    ctx.Schema.Value.MaxLength,
			},
		}

		// add the descriptor to the cache
		r.cache.Add(descriptor)

		return descriptor
	}

	// enum type descriptor
	if kind := r.kind(ctx.Schema.Value); kind == "string" {
		if values := ctx.Schema.Value.Enum; len(values) > 0 {
			descriptor := &TypeDescriptor{
				Name:        inflect.Dasherize(ctx.Name),
				Description: ctx.Schema.Value.Description,
				Default:     ctx.Schema.Value.Default,
				IsNullable:  ctx.Schema.Value.Nullable,
				IsEnum:      true,
				Metadata: Metadata{
					"values": values,
				},
			}

			// add the descriptor to the cache
			r.cache.Add(descriptor)

			return descriptor
		}
	}

	descriptor := &TypeDescriptor{
		Name:        r.kind(ctx.Schema.Value),
		Default:     ctx.Schema.Value.Default,
		IsNullable:  ctx.Schema.Value.Nullable,
		IsPrimitive: true,
	}

	switch ctx.Schema.Value.Type {
	case "string":
		descriptor.Metadata = Metadata{
			"min":     &ctx.Schema.Value.MinLength,
			"max":     ctx.Schema.Value.MaxLength,
			"pattern": ctx.Schema.Value.Pattern,
		}
	case "number", "integer":
		descriptor.Metadata = Metadata{
			"min":           ctx.Schema.Value.Min,
			"max":           ctx.Schema.Value.Max,
			"min_exclusive": ctx.Schema.Value.ExclusiveMin,
			"max_exclusive": ctx.Schema.Value.ExclusiveMax,
			"multiple_of":   ctx.Schema.Value.MultipleOf,
		}
	}

	if ctx.Parent.IsRoot() {
		descriptor = &TypeDescriptor{
			Name:        ctx.Name,
			Description: ctx.Schema.Value.Description,
			IsAlias:     true,
			Element:     descriptor,
		}

		// add the descriptor to the cache
		r.cache.Add(descriptor)
	}

	return descriptor
}

func (r *Resolver) kind(schema *openapi3.Schema) string {
	var (
		kind   = schema.Type
		format = schema.Format
	)

	if kind == "" {
		return "object"
	}

	switch kind {
	case "string":
		switch format {
		case "binary":
			return "binary"
		case "byte":
			return "byte"
		case "uuid":
			return "uuid"
		case "date":
			return "date"
		case "date-time":
			return "date-time"
		default:
			return "string"
		}
	case "integer":
		switch format {
		case "int64":
			return "int64"
		case "int", "int32":
			fallthrough
		default:
			return "int32"
		}
	case "number":
		switch format {
		case "double", "float64":
			return "float64"
		case "float", "float32":
			fallthrough
		default:
			return "float32"
		}
	default:
		return kind
	}
}
