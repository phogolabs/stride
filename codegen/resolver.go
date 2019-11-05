package codegen

import (
	"path/filepath"
	"sort"
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
		Schema: schema,
		Stage:  "property",
	}

	return ctx
}

// Array returns the array context
func (r *ResolverContext) Array() *ResolverContext {
	ctx := &ResolverContext{
		Name:   inflect.Singularize(r.Name),
		Schema: r.Schema.Value.Items,
		Stage:  "array",
	}

	return ctx
}

// Resolver resolves all swagger spec
type Resolver struct{}

// Resolve resolves the spec
func (r *Resolver) Resolve(swagger *openapi3.Swagger) *SpecDescriptor {
	descriptors := TypeDescriptorCollection{}

	descriptors = append(descriptors, r.schemas(swagger.Components.Schemas)...)
	descriptors = append(descriptors, r.parameters(swagger.Components.Parameters)...)
	descriptors = append(descriptors, r.bodies(swagger.Components.RequestBodies)...)
	descriptors = append(descriptors, r.responses(swagger.Components.Responses)...)

	hash := TypeDescriptorMap{}
	hash.CollectFrom(descriptors)

	spec := &SpecDescriptor{
		Types: hash.Collection(),
		// Operations:    r.resolveOperations(root, swagger.Paths),
	}

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

func (r *Resolver) bodies(bodies map[string]*openapi3.RequestBodyRef) TypeDescriptorCollection {
	descriptors := TypeDescriptorCollection{}

	for name, body := range bodies {
		content, ok := body.Value.Content["application/json"]
		if !ok {
			//TODO:
			continue
		}

		ctx := &ResolverContext{
			Name:   name,
			Stage:  "body",
			Schema: content.Schema,
		}

		descriptors = append(descriptors, r.resolve(ctx))
	}

	return descriptors
}

func (r *Resolver) responses(responses map[string]*openapi3.ResponseRef) TypeDescriptorCollection {
	descriptors := TypeDescriptorCollection{}

	for name, response := range responses {
		content, ok := response.Value.Content["application/json"]
		if !ok {
			//TODO:
			continue
		}

		ctx := &ResolverContext{
			Name:   name,
			Stage:  "response",
			Schema: content.Schema,
		}

		descriptors = append(descriptors, r.resolve(ctx))
	}

	return descriptors
}

func (r *Resolver) parameters(parameters map[string]*openapi3.ParameterRef) TypeDescriptorCollection {
	descriptors := TypeDescriptorCollection{}

	for name, parameter := range parameters {
		ctx := &ResolverContext{
			Name:   name,
			Stage:  "parameter",
			Schema: parameter.Value.Schema,
		}

		descriptors = append(descriptors, r.resolve(ctx))
	}
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
