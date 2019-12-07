package codedom

import (
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/phogolabs/flaw"
	"github.com/phogolabs/stride/contract"
	"github.com/phogolabs/stride/inflect"
)

// Resolver resolves all swagger spec
type Resolver struct {
	Cache    TypeDescriptorMap
	Reporter contract.Reporter
}

// Resolve resolves the spec
func (r *Resolver) Resolve(swagger *openapi3.Swagger) (*SpecDescriptor, error) {
	reporter := r.Reporter.With(contract.SeverityVeryHigh)
	reporter.Notice("Resolving spec...")

	defer r.Cache.Clear()

	var (
		ctx         = &ResolverContext{}
		components  = swagger.Components
		controllers = r.operations(ctx, swagger.Paths)
	)

	r.schemas(ctx, components.Schemas)
	r.parameters(ctx, components.Parameters)
	r.headers(ctx, components.Headers)
	r.requests(ctx, components.RequestBodies)
	r.responses(ctx, components.Responses)

	if err := ctx.Collector; len(err) > 0 {
		reporter.Error("Resolving spec fail!")
		return nil, flaw.Errorf("Please check the error log for more details")
	}

	reporter.Success("Resolving spec complete")

	return &SpecDescriptor{
		Types:       r.Cache.Collection(),
		Controllers: controllers,
	}, nil
}

func (r *Resolver) schemas(ctx *ResolverContext, schemas map[string]*openapi3.SchemaRef) TypeDescriptorCollection {
	var (
		collector = flaw.ErrorCollector{}
		reporter  = r.Reporter.With(contract.SeverityHigh)
	)

	reporter.Notice("Resolving schemas...")

	defer func() {
		if err := collector; len(err) > 0 {
			ctx.Collector.Wrap(err)
			reporter.Error("Resolving schemas fail")
		} else {
			reporter.Success("Resolving schemas successful")
		}
	}()

	descriptors := TypeDescriptorCollection{}

	for name, schema := range schemas {
		cctx := ctx.Child(name, schema)
		descriptors = append(descriptors, r.resolve(cctx))

		if err := cctx.Collector; len(err) > 0 {
			collector.Wrap(err)
		}
	}

	return descriptors
}

func (r *Resolver) operations(ctx *ResolverContext, operations map[string]*openapi3.PathItem) ControllerDescriptorCollection {
	var (
		collector = flaw.ErrorCollector{}
		reporter  = r.Reporter.With(contract.SeverityHigh)
	)

	reporter.Notice("Resolving operations...")

	defer func() {
		if err := collector; len(err) > 0 {
			ctx.Collector.Wrap(err)
			reporter.Error("Resolving operations fail")
		} else {
			reporter.Success("Resolving operations successful")
		}
	}()

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
			r.Reporter.Info("Resolving operation: %s method: %v path: %v...",
				inflect.Dasherize(spec.OperationID),
				inflect.UpperCase(method),
				inflect.LowerCase(path))

			var (
				controller = descriptors.Get(key(spec.Tags))
				cctx       = ctx.Child(spec.OperationID, nil)
			)

			var (
				parameterMap = make(map[string]*openapi3.ParameterRef)
				requestMap   = make(map[string]*openapi3.RequestBodyRef)
				responses    = spec.Responses
			)

			requestMap["request"] = spec.RequestBody

			for _, param := range spec.Parameters {
				parameterMap[param.Value.Name] = param
			}

			operation := &OperationDescriptor{
				Path:        path,
				Method:      method,
				Name:        inflect.Dasherize(spec.OperationID),
				Description: spec.Description,
				Summary:     spec.Summary,
				Deprecated:  spec.Deprecated,
				Tags:        spec.Tags,
				Requests:    r.requests(cctx, requestMap),
				Responses:   r.responses(cctx, responses),
			}

			parameters := r.parameters(cctx, parameterMap)

			if len(operation.Requests) == 0 {
				request := &RequestDescriptor{
					ContentType: "application/unknown",
					Description: spec.Description,
				}

				operation.Requests = append(operation.Requests, request)
			}

			for _, request := range operation.Requests {
				request.Parameters = parameters
			}

			controller.Operations = append(controller.Operations, operation)

			if err := cctx.Collector; len(err) > 0 {
				collector.Wrap(err)
				r.Reporter.Error("Resolving operation: %s method: %v path: %v fail",
					inflect.Dasherize(spec.OperationID),
					inflect.UpperCase(method),
					inflect.LowerCase(path))
			} else {
				r.Reporter.Info("Resolving operation: %s method: %v path: %v successful",
					inflect.Dasherize(spec.OperationID),
					inflect.UpperCase(method),
					inflect.LowerCase(path))
			}
		}
	}

	return descriptors.Collection()
}

func (r *Resolver) requests(ctx *ResolverContext, bodies map[string]*openapi3.RequestBodyRef) RequestDescriptorCollection {
	var (
		collector = flaw.ErrorCollector{}
		reporter  = r.Reporter.With(contract.SeverityHigh)
	)

	reporter.Notice("Resolving requests...")

	defer func() {
		if err := collector; len(err) > 0 {
			ctx.Collector.Wrap(err)
			reporter.Error("Resolving requests fail")
		} else {
			reporter.Success("Resolving requests successful")
		}
	}()

	descriptors := RequestDescriptorCollection{}

	for name, spec := range bodies {
		if spec == nil {
			continue
		}

		rcollector := flaw.ErrorCollector{}
		r.Reporter.Info("Resolving request body: %s....", inflect.Dasherize(name))

		for contentType, content := range spec.Value.Content {
			r.Reporter.Info("Resolving request body: %s content-type: %s...",
				inflect.Dasherize(name),
				inflect.LowerCase(contentType))

			schema := content.Schema

			if reference := spec.Ref; reference != "" {
				if schema.Ref == "" {
					schema = &openapi3.SchemaRef{
						Ref:   reference,
						Value: schema.Value,
					}
				}
			}

			var (
				cctx       = ctx.Child(name, schema)
				descriptor = &RequestDescriptor{
					ContentType: contentType,
					Description: spec.Value.Description,
					Required:    spec.Value.Required,
					RequestType: r.resolve(cctx),
				}
			)

			descriptors = append(descriptors, descriptor)

			if err := cctx.Collector; len(err) > 0 {
				rcollector.Wrap(err)
				r.Reporter.Error("Resolving request body: %s content-type: %s fail",
					inflect.Dasherize(name),
					inflect.LowerCase(contentType))
			} else {
				r.Reporter.Info("Resolving request body: %s content-type: %s successful",
					inflect.Dasherize(name),
					inflect.LowerCase(contentType))
			}
		}

		if err := rcollector; len(err) > 0 {
			collector.Wrap(err)
			r.Reporter.Error("Resolving request body: %s fail", inflect.Dasherize(name))
		} else {
			r.Reporter.Info("Resolving request body: %s successful", inflect.Dasherize(name))
		}
	}

	// sort descriptors
	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) responses(ctx *ResolverContext, responses map[string]*openapi3.ResponseRef) ResponseDescriptorCollection {
	var (
		collector = flaw.ErrorCollector{}
		reporter  = r.Reporter.With(contract.SeverityHigh)
	)

	reporter.Notice("Resolving responses...")

	defer func() {
		if err := collector; len(err) > 0 {
			ctx.Collector.Wrap(err)
			reporter.Error("Resolving responses fail")
		} else {
			reporter.Success("Resolving responses successful")
		}
	}()

	var (
		descriptors = ResponseDescriptorCollection{}
		defaultSpec = r.responsesOf(responses)
	)

	for name, spec := range responses {
		var (
			text       = name
			rcollector = flaw.ErrorCollector{}
		)

		code, err := strconv.Atoi(name)
		if err == nil {
			name = inflect.Dasherize(http.StatusText(code)) + "-response"
			text = inflect.Dasherize(ctx.Name) + "-" + name
		} else {
			code = -1
		}

		r.Reporter.Info("Resolving response: %s...", inflect.Dasherize(text))

		if len(spec.Value.Content) == 0 {
			response := &ResponseDescriptor{
				Code:        code,
				ContentType: "application/unknown",
				Description: spec.Value.Description,
				IsDefault:   spec == defaultSpec,
			}

			descriptors = append(descriptors, response)
		}

		for contentType, content := range spec.Value.Content {
			r.Reporter.Info("Resolving response: %s content-type: %s...",
				inflect.Dasherize(text),
				inflect.LowerCase(contentType))

			schema := content.Schema

			if reference := spec.Ref; reference != "" {
				if schema.Ref == "" {
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
					IsDefault:    spec == defaultSpec,
				}
			)

			if length := len(descriptors); code >= 0 && length > 0 {
				prev := descriptors[length-1]

				if !reflect.DeepEqual(prev.ResponseType, response.ResponseType) {
					err := fmt.Errorf("Expecting response: %s content-type: %s body: %s to equal content-type: %s body: %s",
						inflect.Dasherize(text),
						inflect.LowerCase(response.ContentType),
						inflect.Dasherize(response.ResponseType.Name),
						inflect.LowerCase(prev.ContentType),
						inflect.Dasherize(prev.ResponseType.Name),
					)

					reporter := r.Reporter.With(contract.SeverityVeryHigh)
					reporter.Error(err.Error())
					reporter.Error("You cannot have a response with different content-type. The response body should be the same for all content-type declarations")

					cctx.Collector.Wrap(err)
				}

				continue
			} else {
				descriptors = append(descriptors, response)
			}

			if err := cctx.Collector; len(err) > 0 {
				reporter.Error("Resolving response: %s content-type: %s fail", inflect.Dasherize(text), inflect.LowerCase(response.ContentType))
				rcollector.Wrap(err)
			} else {
				r.Reporter.Info("Resolving response: %s content-type: %s successful",
					inflect.Dasherize(text),
					inflect.LowerCase(contentType))
			}
		}

		if err := rcollector; len(err) > 0 {
			collector.Wrap(err)
			r.Reporter.Error("Resolving response: %s fail", inflect.Dasherize(text))
		} else {
			r.Reporter.Info("Resolving response: %s successful", inflect.Dasherize(text))
		}
	}

	// sort descriptors
	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) responsesOf(responses map[string]*openapi3.ResponseRef) *openapi3.ResponseRef {
	if spec, ok := responses["default"]; ok {
		return spec
	}

	var (
		prefix = "2"
		codes  = []string{}
	)

	for name, spec := range responses {
		codes = append(codes, name)

		if strings.HasPrefix(name, prefix) {
			return spec
		}
	}

	sort.Strings(codes)

	if len(codes) > 0 {
		return responses[codes[0]]
	}

	return nil
}

func (r *Resolver) parameters(ctx *ResolverContext, parameters map[string]*openapi3.ParameterRef) ParameterDescriptorCollection {
	var (
		collector = flaw.ErrorCollector{}
		reporter  = r.Reporter.With(contract.SeverityHigh)
	)

	reporter.Notice("Resolving parameters...")

	defer func() {
		if err := collector; len(err) > 0 {
			ctx.Collector.Wrap(err)
			reporter.Error("Resolving parameters fail")
		} else {
			reporter.Success("Resolving parameters successful")
		}
	}()

	descriptors := ParameterDescriptorCollection{}

	for name, spec := range parameters {
		r.Reporter.Info("Resolving parameter: %s...", inflect.Dasherize(name))

		schema := spec.Value.Schema

		if reference := spec.Ref; reference != "" {
			if schema.Ref == "" {
				schema = &openapi3.SchemaRef{
					Ref:   reference,
					Value: schema.Value,
				}
			}
		}

		var (
			cctx      = ctx.Child(name, schema)
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

		if value := spec.Value.Style; value == "" {
			switch spec.Value.In {
			case "header":
				parameter.Style = "form"
			case "query":
				parameter.Style = "form"
			case "path":
				parameter.Style = "simple"
			case "cookie":
				parameter.Style = "form"
			}
		} else {
			parameter.Style = value
		}

		if value := spec.Value.Explode; value == nil {
			switch spec.Value.In {
			case "header":
				parameter.Explode = false
			case "query":
				parameter.Explode = true
			case "path":
				parameter.Explode = false
			case "cookie":
				parameter.Explode = true
			}
		} else {
			parameter.Explode = *value
		}

		descriptors = append(descriptors, parameter)

		if err := cctx.Collector; len(err) > 0 {
			collector.Wrap(err)
			r.Reporter.Error("Resolving parameter: %s fail", inflect.Dasherize(name))
		} else {
			r.Reporter.Info("Resolving parameter: %s successful", inflect.Dasherize(name))
		}
	}

	// sort parameters
	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) headers(ctx *ResolverContext, headers map[string]*openapi3.HeaderRef) ParameterDescriptorCollection {
	var (
		collector = flaw.ErrorCollector{}
		reporter  = r.Reporter.With(contract.SeverityHigh)
	)

	reporter.Notice("Resolving headers...")

	defer func() {
		if err := collector; len(err) > 0 {
			ctx.Collector.Wrap(err)
			reporter.Error("Resolving headers fail")
		} else {
			reporter.Success("Resolving headers successful")
		}
	}()

	descriptors := ParameterDescriptorCollection{}

	for name, spec := range headers {
		r.Reporter.Info("Resolving header: %s...", inflect.Dasherize(name))

		schema := spec.Value.Schema

		if reference := spec.Ref; reference != "" {
			if schema.Ref == "" {
				schema = &openapi3.SchemaRef{
					Ref:   reference,
					Value: schema.Value,
				}
			}
		}

		var (
			cctx   = ctx.Child(name, schema)
			header = &ParameterDescriptor{
				Name:          name,
				In:            "header",
				ParameterType: r.resolve(cctx),
			}
		)

		descriptors = append(descriptors, header)

		if err := cctx.Collector; len(err) > 0 {
			collector.Wrap(err)
			r.Reporter.Error("Resolving header: %s fail", inflect.Dasherize(name))
		} else {
			r.Reporter.Info("Resolving header: %s successful", inflect.Dasherize(name))
		}
	}

	// sort descriptors
	sort.Sort(descriptors)

	return descriptors
}

func (r *Resolver) add(descriptor *TypeDescriptor) error {
	if err := r.Cache.Add(descriptor); err != nil {
		reporter := r.Reporter.With(contract.SeverityVeryHigh)
		reporter.Error("Resolving type: %s fail: %v ", inflect.Dasherize(descriptor.Name), err)
		reporter.Error("Please check your OpenAPI spec for duplicated name: '%v'", descriptor.Name)
		reporter.Error("The requests, responses, parameters, headers should have unique names across the whole document.")
		return err
	}

	return nil
}

func (r *Resolver) resolve(ctx *ResolverContext) *TypeDescriptor {
	var (
		collector = flaw.ErrorCollector{}
		reporter  = r.Reporter.With(contract.SeverityLow)
	)

	reporter.Notice("Resolving type: %s...", inflect.Dasherize(ctx.Name))

	defer func() {
		if err := collector; len(err) > 0 {
			ctx.Collector.Wrap(err)
			reporter.Error("Resolving type: %s fail", inflect.Dasherize(ctx.Name))
		} else {
			reporter.Success("Resolving type: %s successful", inflect.Dasherize(ctx.Name))
		}
	}()

	switch {
	case ctx.Schema == nil:
	case ctx.Schema.Value.OneOf != nil:
		reporter.Warn("Resolving type: %s does not support 'one-of' clause. Reverting to generic type", inflect.Dasherize(ctx.Name))
		ctx.Schema = nil
	case ctx.Schema.Value.AnyOf != nil:
		reporter.Warn("Resolving type: %s does not support 'any-of' clause. Reverting to generic type", inflect.Dasherize(ctx.Name))
		ctx.Schema = nil
	case ctx.Schema.Value.AllOf != nil:
		reporter.Warn("Resolving type: %s does not support 'all-of' clause. Reverting to generic type", inflect.Dasherize(ctx.Name))
		ctx.Schema = nil
	case ctx.Schema.Value.Not != nil:
		reporter.Warn("Resolving type: %s does not support 'not' clause. Reverting to generic type", inflect.Dasherize(ctx.Name))
		ctx.Schema = nil
	}

	if ctx.Schema == nil {
		descriptor := &TypeDescriptor{
			IsAny: true,
		}

		if ctx.Parent.IsRoot() {
			descriptor = &TypeDescriptor{
				Name:    inflect.Dasherize(ctx.Name),
				IsAlias: true,
				Element: descriptor,
			}

			// add the descriptor to the cache
			if err := r.add(descriptor); err != nil {
				collector.Wrap(err)
			}
		}

		return descriptor
	}

	// reference type descriptor
	if reference := ctx.Schema.Ref; reference != "" {
		reporter.Info("Resolving type: %s to alias...", inflect.Dasherize(ctx.Name))

		var (
			cctx       = ctx.Dereference()
			descriptor = r.resolve(cctx)
		)

		if ctx.Parent.IsRoot() {
			descriptor = &TypeDescriptor{
				Name:        inflect.Dasherize(ctx.Name),
				Description: ctx.Schema.Value.Description,
				IsAlias:     true,
				Element:     descriptor,
			}

			// add the descriptor to the cache
			if err := r.add(descriptor); err != nil {
				cctx.Collector.Wrap(err)
			}
		}

		if err := cctx.Collector; len(err) > 0 {
			collector.Wrap(err)
			reporter.Error("Resolving type: %s to alias fail", inflect.Dasherize(ctx.Name))
		} else {
			reporter.Info("Resolving type: %s to alias successful", inflect.Dasherize(ctx.Name))
		}

		return descriptor
	}

	// class type descriptor
	if kind := r.kind(ctx.Schema.Value); kind == "object" {
		reporter.Info("Resolving type: %s to class...", inflect.Dasherize(ctx.Name))

		descriptor := &TypeDescriptor{
			Name:        inflect.Dasherize(ctx.Name),
			Description: ctx.Schema.Value.Description,
			IsClass:     true,
			IsNullable:  true,
		}

		//TODO: handle min and max properties somehow
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
			reporter.Info("Resolving type: %s field: %s...",
				inflect.Dasherize(ctx.Name),
				inflect.Dasherize(field))

			var (
				cctx     = ctx.Child(field, schema)
				property = &PropertyDescriptor{
					Name:         field,
					Description:  schema.Value.Description,
					ReadOnly:     schema.Value.ReadOnly,
					WriteOnly:    schema.Value.WriteOnly,
					Required:     required(field),
					PropertyType: r.resolve(cctx),
				}
			)

			descriptor.Properties = append(descriptor.Properties, property)

			if err := cctx.Collector; len(err) > 0 {
				collector.Wrap(err)
				reporter.Error("Resolving type: %s field: %s fail",
					inflect.Dasherize(ctx.Name),
					inflect.Dasherize(field))
			} else {
				reporter.Info("Resolving type: %s field: %s successful",
					inflect.Dasherize(ctx.Name),
					inflect.Dasherize(field))
			}
		}

		switch {
		case ctx.Schema.Value.AdditionalPropertiesAllowed != nil:
			fallthrough
		case ctx.Schema.Value.AdditionalProperties != nil:
			var (
				schema   = ctx.Schema.Value.AdditionalProperties
				kctx     = ctx.Child("key", schemaOf("string"))
				pctx     = ctx.Child("properties", schema)
				property = &PropertyDescriptor{
					Name: "properties",
					PropertyType: &TypeDescriptor{
						Key:     r.resolve(kctx),
						Element: r.resolve(pctx),
						Metadata: Metadata{
							"min":           uint64Ptr(&ctx.Schema.Value.MinProps),
							"max":           uint64Ptr(ctx.Schema.Value.MaxProps),
							"min_exclusive": false,
							"max_exclusive": false,
						},
						IsMap: true,
					},
					IsEmbedded: true,
				}
			)

			descriptor.Properties = append(descriptor.Properties, property)

			if err := kctx.Collector; len(err) > 0 {
				collector.Wrap(err)
			}

			if err := pctx.Collector; len(err) > 0 {
				collector.Wrap(err)
			}
		case !descriptor.HasProperties():
			var (
				kctx = ctx.Child("key", schemaOf("string"))
				pctx = ctx.Child("properties", nil)
			)

			descriptor = &TypeDescriptor{
				Name:    inflect.Dasherize(ctx.Name),
				Key:     r.resolve(kctx),
				Element: r.resolve(pctx),
				Metadata: Metadata{
					"min":           uint64Ptr(&ctx.Schema.Value.MinProps),
					"max":           uint64Ptr(ctx.Schema.Value.MaxProps),
					"min_exclusive": false,
					"max_exclusive": false,
				},
				IsMap: true,
			}

			if err := kctx.Collector; len(err) > 0 {
				collector.Wrap(err)
			}

			if err := pctx.Collector; len(err) > 0 {
				collector.Wrap(err)
			}
		}

		// sort properties by name
		sort.Sort(descriptor.Properties)

		// add the descriptor to the cache
		if err := r.add(descriptor); err != nil {
			collector.Wrap(err)
		}

		if err := collector; len(err) > 0 {
			reporter.Error("Resolving type: %s to class fail", inflect.Dasherize(ctx.Name))
		} else {
			reporter.Info("Resolving type: %s to class successful", inflect.Dasherize(ctx.Name))
		}

		return descriptor
	}

	// array descriptor
	if kind := r.kind(ctx.Schema.Value); kind == "array" {
		reporter.Info("Resolving type: %s to array...", inflect.Dasherize(ctx.Name))

		cctx := ctx.Array()
		descriptor := &TypeDescriptor{
			Name:        inflect.Dasherize(ctx.Name),
			Description: ctx.Schema.Value.Description,
			Default:     ctx.Schema.Value.Default,
			IsNullable:  ctx.Schema.Value.Nullable,
			IsArray:     true,
			Element:     r.resolve(cctx),
			Metadata: Metadata{
				"unique":        ctx.Schema.Value.UniqueItems,
				"min":           uint64Ptr(&ctx.Schema.Value.MinLength),
				"max":           uint64Ptr(ctx.Schema.Value.MaxLength),
				"min_exclusive": ctx.Schema.Value.ExclusiveMin,
				"max_exclusive": ctx.Schema.Value.ExclusiveMax,
			},
		}

		if err := cctx.Collector; len(err) > 0 {
			collector.Wrap(err)
		}

		// add the descriptor to the cache
		if err := r.add(descriptor); err != nil {
			collector.Wrap(err)
		}

		if err := collector; len(err) > 0 {
			reporter.Error("Resolving type: %s to array fail", inflect.Dasherize(ctx.Name))
		} else {
			reporter.Info("Resolving type: %s to array successful", inflect.Dasherize(ctx.Name))
		}

		return descriptor
	}

	// enum type descriptor
	if kind := r.kind(ctx.Schema.Value); kind == "string" {
		if values := ctx.Schema.Value.Enum; len(values) > 0 {
			reporter.Info("Resolving type: %s to enum...", inflect.Dasherize(ctx.Name))

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
			if err := r.add(descriptor); err != nil {
				collector.Wrap(err)
			}

			if err := collector; len(err) > 0 {
				reporter.Error("Resolving type: %s to enum fail", inflect.Dasherize(ctx.Name))
			} else {
				reporter.Info("Resolving type: %s to enum successful", inflect.Dasherize(ctx.Name))
			}

			return descriptor
		}
	}

	reporter.Info("Resolving type: %s to primitive...", inflect.Dasherize(ctx.Name))

	descriptor := &TypeDescriptor{
		Name:        r.kind(ctx.Schema.Value),
		Default:     ctx.Schema.Value.Default,
		IsNullable:  ctx.Schema.Value.Nullable,
		IsPrimitive: true,
	}

	switch r.kind(ctx.Schema.Value) {
	case "string", "byte", "binary":
		descriptor.Metadata = Metadata{
			"min":           uint64Ptr(&ctx.Schema.Value.MinLength),
			"max":           uint64Ptr(ctx.Schema.Value.MaxLength),
			"pattern":       ctx.Schema.Value.Pattern,
			"min_exclusive": ctx.Schema.Value.ExclusiveMin,
			"max_exclusive": ctx.Schema.Value.ExclusiveMax,
		}
	case "int32", "int64", "float32", "float64":
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
		if err := r.add(descriptor); err != nil {
			collector.Wrap(err)
		}
	}

	if err := collector; len(err) > 0 {
		reporter.Error("Resolving type: %s to primitive fail", inflect.Dasherize(ctx.Name))
	} else {
		reporter.Info("Resolving type: %s to primitive successful", inflect.Dasherize(ctx.Name))
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

func uint64Ptr(v *uint64) *float64 {
	if v == nil {
		return nil
	}

	f := float64(*v)
	return &f
}
