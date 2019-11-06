package codegen

import (
	"bytes"
	"fmt"
	"go/token"
	"io"
	"os"
	"path/filepath"
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
type Generator struct {
	Path string
}

// Generate generates the source code
func (g *Generator) Generate(spec *SpecDescriptor) error {
	path := filepath.Join(g.Path, "service")

	// prepare the service package directory
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	// write types
	if err := g.write(g.filename("contract"), g.types(spec.Types)); err != nil {
		return err
	}

	// write controllers
	for _, descriptor := range spec.Controllers {
		name := descriptor.Name + "_api"
		if err := g.write(g.filename(name), g.controller(descriptor)); err != nil {
			return err
		}

		spec := descriptor.Name + "_api_test"
		if err := g.write(g.filename(spec), g.spec(descriptor)); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) write(name string, decls []dst.Decl) error {
	pkg := filepath.Base(filepath.Dir(name))

	if strings.HasSuffix(name, "_test.go") {
		pkg = pkg + "_test"
	}

	root := &dst.File{
		Name: &dst.Ident{
			Name: pkg,
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

func (g *Generator) controller(descriptor *ControllerDescriptor) []dst.Decl {
	node := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: &dst.Ident{
					Name: descriptor.Name,
				},
				Type: &dst.StructType{
					Fields: &dst.FieldList{
						List: []*dst.Field{},
					},
					Incomplete: true,
				},
			},
		},
	}

	return []dst.Decl{node}
}

func (g *Generator) spec(descriptor *ControllerDescriptor) []dst.Decl {
	return []dst.Decl{}
}

func (g *Generator) filename(name string) string {
	name = inflect.Underscore(name) + ".go"
	return filepath.Join(g.Path, name)
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

	var (
		name = builder.camelize(descriptor.Name)
		node = &dst.GenDecl{
			Tok: token.TYPE,
			Specs: []dst.Spec{
				&dst.TypeSpec{
					Name: &dst.Ident{
						Name: name,
					},
					Type: expr,
				},
			},
		}
	)

	node.Decs.Before = dst.NewLine
	node.Decs.Start.Append(builder.commentf("%s is a struct type auto-generated from OpenAPI spec", name))

	if descriptor.Description != "" {
		node.Decs.Start.Append(builder.commentf(descriptor.Description))
	}

	node.Decs.Start.Append(builder.commentf("stride:generate"))

	return node
}

func (builder *TypeBuilder) camelize(text string) string {
	var (
		field  = inflect.Camelize(text)
		buffer = &bytes.Buffer{}
		suffix = "Id"
	)

	switch {
	case field == suffix:
		buffer.WriteString(strings.ToUpper(field))
	case strings.HasSuffix(field, suffix):
		buffer.WriteString(strings.TrimSuffix(field, suffix))
		buffer.WriteString(strings.ToUpper(suffix))
	default:
		buffer.WriteString(field)
	}

	return buffer.String()
}

func (builder *TypeBuilder) commentf(text string, args ...interface{}) string {
	text = fmt.Sprintf(text, args...)
	return fmt.Sprintf("// %s", text)
}

func (builder *TypeBuilder) kind(descriptor *TypeDescriptor) string {
	var (
		elem = element(descriptor)
		name = descriptor.Name
	)

	switch descriptor.Name {
	case "date-time":
		name = "time.Time"
	case "date":
		name = "time.Time"
	case "uuid":
		name = "schema.UUID"
	}

	if elem.IsClass {
		if len(elem.Properties) == 0 {
			return "interface{}"
		}
	}

	if elem.IsNullable {
		return fmt.Sprintf("*%s", name)
	}

	return name
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

	// TODO: uncomment when you add xml support
	// tags.Set(&structtag.Tag{
	// 	Key:     "xml",
	// 	Name:    property.Name,
	// 	Options: builder.omitempty(property),
	// })

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

	if len(options) == 0 {
		options = append(options, "-")
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
