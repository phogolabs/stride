package codegen

import (
	"fmt"
	"go/token"

	"github.com/dave/dst"
)

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
						Name: camelize(property.Name),
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
				field.Decs.Type.Append(commentf(property.Description))
			}
		}
	case descriptor.IsEnum:
		//TODO: generate values
		fallthrough
	default:
		return nil
	}

	var (
		name = camelize(descriptor.Name)
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
	node.Decs.Start.Append(commentf("%s is a struct type auto-generated from OpenAPI spec", name))

	if descriptor.Description != "" {
		node.Decs.Start.Append(commentf(descriptor.Description))
	}

	node.Decs.Start.Append(commentf("stride:generate"))

	return node
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
