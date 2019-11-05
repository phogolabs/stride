package codegen

import (
	"bytes"
	"fmt"
	"go/token"
	"io"
	"os"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/fatih/structtag"
	"github.com/go-openapi/inflect"
)

// Renderer renders the dom
type Renderer interface {
	Render(w io.Writer)
}

// Generator generates the source code
type Generator struct{}

// Generate generates the source code
func (g *Generator) Generate(spec *SpecDescriptor) error {
	if err := g.write("contract.go", g.types(spec.Types)); err != nil {
		return err
	}

	return nil
}

func (g *Generator) write(name string, decls []dst.Decl) error {
	root := &dst.File{
		Name: &dst.Ident{
			Name: "service",
		},
		Decls: decls,
	}

	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := decorator.Fprint(file, root); err != nil {
		return err
	}

	return nil
}

func (g *Generator) types(descriptors TypeDescriptorCollection) []dst.Decl {
	var (
		tree    = []dst.Decl{}
		builder = &TypeBuilder{}
	)

	for _, descriptor := range descriptors {
		node := builder.Build(descriptor)
		tree = append(tree, node)
	}

	return tree
}

// TypeBuilder builds a type
type TypeBuilder struct{}

// Build builds a type
func (builder *TypeBuilder) Build(descriptor *TypeDescriptor) dst.Decl {
	var (
		expr dst.Expr
		tag  = &TagBuilder{}
	)

	switch {
	case descriptor.IsAlias:
		expr = &dst.Ident{
			Name: descriptor.Element.Name,
		}
	case descriptor.IsArray:
		expr = &dst.ArrayType{
			Elt: &dst.Ident{
				Name: builder.kind(descriptor.Element),
			},
		}
	case descriptor.IsClass:
		spec := &dst.StructType{
			Fields: &dst.FieldList{},
		}
		expr = spec

		for _, property := range descriptor.Properties {
			field := &dst.Field{
				Names: []*dst.Ident{
					&dst.Ident{
						Name: builder.camelize(property.Name),
					},
				},
				Type: &dst.Ident{
					Name: builder.kind(property.PropertyType),
				},
				Tag: &dst.BasicLit{
					Kind:  token.STRING,
					Value: tag.Build(property),
				},
			}

			spec.Fields.List = append(spec.Fields.List, field)

			if property.Description != "" {
				field.Decs.Type.Append(builder.commentf(property.Description))
			}
		}
	case descriptor.IsEnum:
		//TODO: generate values
		fallthrough
	default:
		return nil
	}

	node := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: &dst.Ident{
					Name: descriptor.Name,
				},
				Type: expr,
			},
		},
	}

	node.Decs.Before = dst.NewLine
	node.Decs.Start.Append(builder.commentf("%s is a struct type auto-generated from OpenAPI spec", descriptor.Name))

	if descriptor.Description != "" {
		node.Decs.Start.Append(builder.commentf(descriptor.Description))
	}

	node.Decs.Start.Append(builder.commentf("stride:generate"))

	return node
}

func (builder *TypeBuilder) camelize(text string) string {
	field := inflect.Camelize(text)

	switch {
	case field == "Id":
		field = strings.ToUpper(field)
	case strings.HasSuffix(field, "Id"):
		field = fmt.Sprintf("%vID", strings.TrimSuffix(field, "Id"))
	}

	return field
}

func (builder *TypeBuilder) commentf(text string, args ...interface{}) string {
	text = fmt.Sprintf(text, args...)
	return fmt.Sprintf("// %s", text)
}

func (builder *TypeBuilder) kind(descriptor *TypeDescriptor) string {
	item := element(descriptor)

	if item.IsClass || item.IsNullable {
		return fmt.Sprintf("*%s", descriptor.Name)
	}

	return descriptor.Name
}

// TagBuilder builds the tags for given type
type TagBuilder struct{}

// Build returns the tag for given property
func (builder *TagBuilder) Build(property *PropertyDescriptor) string {
	tags := &structtag.Tags{}

	tags.Set(&structtag.Tag{
		Key:     "json",
		Name:    property.Name,
		Options: builder.omitempty(property),
	})

	tags.Set(&structtag.Tag{
		Key:     "xml",
		Name:    property.Name,
		Options: builder.omitempty(property),
	})

	if value := property.PropertyType.Default; value != nil {
		tags.Set(&structtag.Tag{
			Key:  "default",
			Name: fmt.Sprintf("%v", value),
		})
	}

	if options := builder.validate(property); len(options) > 0 {
		tags.Set(&structtag.Tag{
			Key:     "validate",
			Name:    options[0],
			Options: options[1:],
		})
	}

	return fmt.Sprintf("`%s`", tags.String())
}

func (builder *TagBuilder) omitempty(property *PropertyDescriptor) []string {
	options := []string{}
	if !property.Required {
		options = append(options, "omitempty")
	}

	return options
}

func (builder *TagBuilder) validate(property *PropertyDescriptor) []string {
	var (
		options  = []string{}
		metadata = element(property.PropertyType).Metadata
	)

	if property.Required {
		options = append(options, "required")
	}

	for k, v := range metadata {
		switch k {
		case "unique":
			if unique, ok := v.(bool); ok {
				if unique {
					options = append(options, "unique")
				}
			}
		case "pattern":
			if value, ok := v.(string); ok {
				if value != "" {
					// TODO: add support for regex
					// options = append(options, fmt.Sprintf("regex=%v", value))
				}
			}
		case "multiple_of":
			if value, ok := v.(*float64); ok {
				if value != nil {
					// TODO: add support for multileof
					// options = append(options, fmt.Sprintf("multipleof=%v", value))
				}
			}
		case "min":
			if value, ok := v.(*float64); ok {
				if value != nil {
					if exclusive, ok := metadata["min_exclusive"].(bool); ok {
						if exclusive {
							options = append(options, fmt.Sprintf("gt=%v", *value))
						} else {
							options = append(options, fmt.Sprintf("gte=%v", *value))
						}
					}
				}
			}
		case "max":
			if value, ok := v.(*float64); ok {
				if value != nil {
					if exclusive, ok := metadata["max_exclusive"].(bool); ok {
						if exclusive {
							options = append(options, fmt.Sprintf("lt=%v", *value))
						} else {
							options = append(options, fmt.Sprintf("lte=%v", *value))
						}
					}
				}
			}
		case "values":
			if values, ok := v.([]interface{}); ok {
				if len(values) > 0 {
					options = append(options, fmt.Sprintf("oneof=%v", builder.oneof(values)))
				}
			}
		}
	}

	return options
}

func (builder *TagBuilder) oneof(values []interface{}) string {
	buffer := &bytes.Buffer{}

	for index, value := range values {
		if index > 0 {
			fmt.Fprint(buffer, " ")
		}

		fmt.Fprintf(buffer, "%v", value)
	}

	return buffer.String()
}

func element(descriptor *TypeDescriptor) *TypeDescriptor {
	element := descriptor

	for element.IsAlias {
		element = element.Element
	}

	return element
}
