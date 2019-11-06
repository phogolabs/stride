package codegen

import (
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/inflect"
)

// ResolverContext is the current resolver context
type ResolverContext struct {
	Name   string
	Stage  string
	Schema *openapi3.SchemaRef
}

// Referenced returns the referenced context
func (r *ResolverContext) Referenced() *ResolverContext {
	ctx := &ResolverContext{
		Name:   filepath.Base(r.Schema.Ref),
		Schema: &openapi3.SchemaRef{Value: r.Schema.Value},
		Stage:  "reference",
	}

	return ctx
}

// Property returns the property context
func (r *ResolverContext) Property(name string, schema *openapi3.SchemaRef) *ResolverContext {
	ctx := &ResolverContext{
		Name:   inflect.Camelize(r.Name + "_" + name),
		Stage:  "property",
		Schema: schema,
	}

	return ctx
}

// Array returns the array context
func (r *ResolverContext) Array() *ResolverContext {
	ctx := &ResolverContext{
		Name:   inflect.Singularize(r.Name),
		Stage:  "array",
		Schema: r.Schema.Value.Items,
	}

	return ctx
}

// Resolver resolves all swagger spec
type Resolver struct{}

// Resolve resolves the spec
func (r *Resolver) Resolve(swagger *openapi3.Swagger) *SpecDescriptor {
	spec := &SpecDescriptor{
		Controllers: r.operations(swagger.Paths),
	}

	descriptors := TypeDescriptorMap{}
	descriptors.CollectFromSchemas(r.schemas(swagger.Components.Schemas))
	descriptors.CollectFromParameters(r.parameters(swagger.Components.Parameters))
	descriptors.CollectFromHeaders(r.headers(swagger.Components.Headers))
	descriptors.CollectFromRequests(r.requests(swagger.Components.RequestBodies))
	descriptors.CollectFromResponses(r.responses(swagger.Components.Responses))
	// descriptors.CollectFromControllers(spec.Controllers)

	spec.Types = descriptors.Collection()

	return spec
}

func (r *Resolver) schemas(schemas map[string]*openapi3.SchemaRef) TypeDescriptorCollection {
	descriptors := TypeDescriptorCollection{}

	for name, schema := range schemas {
		ctx := &ResolverContext{
			Name:   name,
			Stage:  "schema",
			Schema: schema,
		}

		descriptors = append(descriptors, r.resolve(ctx))
	}

	return descriptors
}

func (r *Resolver) operations(operations map[string]*openapi3.PathItem) ControllerDescriptorCollection {
	descriptors := ControllerDescriptorMap{}

	for path, spec := range operations {
		for method, spec := range spec.Operations() {
			controller := descriptors.Get(spec.Tags)

			var (
				parameters = make(map[string]*openapi3.ParameterRef)
				requests   = make(map[string]*openapi3.RequestBodyRef)
				responses  = spec.Responses
			)

			requests[spec.OperationID] = spec.RequestBody

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
				Parameters:  r.parameters(parameters),
				Requests:    r.requests(requests),
				Responses:   r.responses(responses),
			}

			controller.Operations = append(controller.Operations, operation)
		}
	}

	return descriptors.Collection()
}

func (r *Resolver) requests(bodies map[string]*openapi3.RequestBodyRef) RequestDescriptorCollection {
	descriptors := RequestDescriptorCollection{}

	for name, spec := range bodies {
		if spec == nil {
			continue
		}

		for contentType, content := range spec.Value.Content {
			ctx := &ResolverContext{
				Name:   name,
				Stage:  "request",
				Schema: content.Schema,
			}

			request := &RequestDescriptor{
				ContentType: contentType,
				Description: spec.Value.Description,
				Required:    spec.Value.Required,
				RequestType: r.resolve(ctx),
			}

			descriptors = append(descriptors, request)
		}
	}

	// sort descriptors
	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) responses(responses map[string]*openapi3.ResponseRef) ResponseDescriptorCollection {
	descriptors := ResponseDescriptorCollection{}

	for name, spec := range responses {
		code, err := strconv.Atoi(name)
		if err != nil {
			code = 0
		}

		for contentType, content := range spec.Value.Content {
			ctx := &ResolverContext{
				Name:   name,
				Stage:  "response",
				Schema: content.Schema,
			}

			response := &ResponseDescriptor{
				Code:         code,
				ContentType:  contentType,
				Description:  spec.Value.Description,
				ResponseType: r.resolve(ctx),
				Headers:      r.headers(spec.Value.Headers),
			}

			descriptors = append(descriptors, response)
		}
	}

	// sort descriptors
	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) parameters(parameters map[string]*openapi3.ParameterRef) ParameterDescriptorCollection {
	descriptors := ParameterDescriptorCollection{}

	for name, spec := range parameters {
		ctx := &ResolverContext{
			Name:   name,
			Stage:  "parameter",
			Schema: spec.Value.Schema,
		}

		parameter := &ParameterDescriptor{
			Name:          spec.Value.Name,
			In:            spec.Value.In,
			Style:         spec.Value.Style,
			Explode:       spec.Value.Explode,
			Description:   spec.Value.Description,
			Required:      spec.Value.Required,
			Deprecated:    spec.Value.Deprecated,
			ParameterType: r.resolve(ctx),
		}

		descriptors = append(descriptors, parameter)
	}

	// sort parameters
	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) headers(headers map[string]*openapi3.HeaderRef) HeaderDescriptorCollection {
	descriptors := HeaderDescriptorCollection{}

	for name, spec := range headers {
		ctx := &ResolverContext{
			Name:   name,
			Stage:  "header",
			Schema: spec.Value.Schema,
		}

		header := &HeaderDescriptor{
			Name:       name,
			HeaderType: r.resolve(ctx),
		}

		descriptors = append(descriptors, header)
	}

	// sort descriptors
	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) resolve(ctx *ResolverContext) *TypeDescriptor {
	// reference type descriptor
	if reference := ctx.Schema.Ref; reference != "" {
		descriptor := r.resolve(ctx.Referenced())

		switch ctx.Stage {
		case "array":
			return descriptor
		case "property":
			return descriptor
		default:
			return &TypeDescriptor{
				Name:        ctx.Name,
				Description: ctx.Schema.Value.Description,
				IsAlias:     true,
				Element:     descriptor,
			}
		}
	}

	// class type descriptor
	if kind := kind(ctx.Schema.Value); kind == "object" {
		descriptor := &TypeDescriptor{
			Name:        ctx.Name,
			Description: ctx.Schema.Value.Description,
			IsClass:     true,
			IsNullable:  true,
		}

		//TODO: handle min and max properites somehow
		//TODO: handle additional properties
		//TODO: handle pattern properties
		//TODO: handle discriminator
		//TODO: handle read and write only

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
				Required:     required(field),
				PropertyType: r.resolve(ctx.Property(field, schema)),
			}

			descriptor.Properties = append(descriptor.Properties, property)
		}

		// sort properties by name
		sort.Sort(descriptor.Properties)

		return descriptor
	}

	// array descriptor
	if kind := kind(ctx.Schema.Value); kind == "array" {
		descriptor := &TypeDescriptor{
			Name:        ctx.Name,
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

		return descriptor
	}

	// enum type descriptor
	if kind := kind(ctx.Schema.Value); kind == "string" {
		if values := ctx.Schema.Value.Enum; len(values) > 0 {
			descriptor := &TypeDescriptor{
				Name:        ctx.Name,
				Description: ctx.Schema.Value.Description,
				Default:     ctx.Schema.Value.Default,
				IsNullable:  ctx.Schema.Value.Nullable,
				IsEnum:      true,
				Metadata: Metadata{
					"values": values,
				},
			}

			return descriptor
		}
	}

	descriptor := &TypeDescriptor{
		Name:        kind(ctx.Schema.Value),
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

	switch ctx.Stage {
	case "property":
		return descriptor
	case "array":
		return descriptor
	default:
		return &TypeDescriptor{
			Name:        ctx.Name,
			Description: ctx.Schema.Value.Description,
			IsAlias:     true,
			Element:     descriptor,
		}
	}
}

func name(names ...string) string {
	return strings.Join(names, "_")
}

func kind(schema *openapi3.Schema) string {
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
		case "uuid":
			return "schema.UUID"
		case "duration":
			return "time.Duration"
		case "date", "date-time":
			return "time.Time"
		default:
			return "string"
		}
	case "integer":
		switch format {
		case "int64":
			return "int64"
		default:
			return "int32"
		}
	case "number":
		switch format {
		case "float64":
			return "float64"
		default:
			return "float32"
		}
	default:
		return kind
	}
}
