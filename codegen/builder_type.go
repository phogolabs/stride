package codegen

import (
	"go/token"

	"github.com/dave/dst"
)

// TypeBuilder builds a type
type TypeBuilder struct{}

// Build builds a type
func (builder *TypeBuilder) Build(descriptor *TypeDescriptor) dst.Decl {
	var expr dst.Expr

	switch {
	case descriptor.IsAlias:
		expr = &dst.Ident{
			Name: descriptor.Element.Name,
		}
	case descriptor.IsArray:
		expr = &dst.ArrayType{
			Elt: &dst.Ident{
				Name: descriptor.Element.Kind(),
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
						Name: property.Name,
					},
				},
				Type: &dst.Ident{
					Name: property.PropertyType.Kind(),
				},
				// Tag: &dst.BasicLit{
				// 	Kind: token.STRING,
				//  Value: property.Tags(),
				// },
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
		node = &dst.GenDecl{
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
	)

	node.Decs.Before = dst.NewLine
	node.Decs.Start.Append(commentf("%s is a struct type auto-generated from OpenAPI spec", descriptor.Name))

	if descriptor.Description != "" {
		node.Decs.Start.Append(commentf(descriptor.Description))
	}

	node.Decs.Start.Append(commentf("stride:generate"))

	return node
}
