package codegen

import (
	"bytes"
	"fmt"
	"go/token"
	"sort"
	"strings"

	"github.com/dave/dst"
	"github.com/fatih/structtag"
	"github.com/go-openapi/inflect"
)

// Builder builds a node from the code
type Builder interface {
	Build() dst.Decl
}

// File represents the generated file
type File struct {
	Name    string
	Content *dst.File
}

// FileBuilder builds a file
type FileBuilder struct {
	Package string
}

// Build builds the file
func (b *FileBuilder) Build() *dst.File {
	root := &dst.File{
		Name: &dst.Ident{
			Name: b.Package,
		},
	}

	return root
}

// Field represents a field
type Field struct {
	Name string
	Type string
	Tags []*Tag
}

// Tag represents a tag
type Tag struct {
	Key     string
	Name    string
	Options []string
}

// StructTypeBuilder builds a struct
type StructTypeBuilder struct {
	Name     string
	Comments []string
	Fields   []*Field
}

// Commentf adds a comment
func (b *StructTypeBuilder) Commentf(pattern string, args ...interface{}) {
	if pattern == "" {
		return
	}

	b.Comments = append(b.Comments, commentf(pattern, args...))
}

// Build builds the type
func (b *StructTypeBuilder) Build() dst.Decl {
	var expr dst.Expr

	if spec := b.spec(); spec != nil {
		expr = spec
	}

	node := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: &dst.Ident{
					Name: b.Name,
				},
				Type: expr,
			},
		},
	}

	node.Decs.Before = dst.EmptyLine
	node.Decs.After = dst.EmptyLine

	for _, text := range b.Comments {
		node.Decs.Start.Append(text)
	}

	node.Decs.Start.Append(commentf("stride:generate"))

	return node
}

func (b *StructTypeBuilder) spec() dst.Expr {
	var (
		tags = &structtag.Tags{}
		spec = &dst.StructType{
			Fields: &dst.FieldList{},
		}
	)

	for _, descriptor := range b.Fields {
		field := &dst.Field{
			Names: []*dst.Ident{
				&dst.Ident{
					Name: descriptor.Name,
				},
			},
			Type: &dst.Ident{
				Name: descriptor.Type,
			},
		}

		for _, descriptor := range descriptor.Tags {
			tags.Set(&structtag.Tag{
				Key:     descriptor.Key,
				Name:    descriptor.Name,
				Options: descriptor.Options,
			})
		}

		if value := tags.String(); value != "" {
			field.Tag = &dst.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("`%s`", value),
			}
		}

		spec.Fields.List = append(spec.Fields.List, field)
	}

	return spec
}

// LiteralTypeBuilder builds a literal type
type LiteralTypeBuilder struct {
	Name     string
	Element  string
	Comments []string
}

// Build builds the type
func (b *LiteralTypeBuilder) Build() dst.Decl {
	node := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: &dst.Ident{
					Name: b.Name,
				},
				Type: &dst.ArrayType{
					Elt: &dst.Ident{
						Name: b.Element,
					},
				},
			},
		},
	}

	node.Decs.Before = dst.EmptyLine
	node.Decs.After = dst.EmptyLine

	for _, text := range b.Comments {
		node.Decs.Start.Append(text)
	}

	node.Decs.Start.Append(commentf("stride:generate"))

	return node
}

// ArrayTypeBuilder builds an array type
type ArrayTypeBuilder struct {
	Name     string
	Element  string
	Comments []string
}

// Build builds the type
func (b *ArrayTypeBuilder) Build() dst.Decl {
	node := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: &dst.Ident{
					Name: b.Name,
				},
				Type: &dst.Ident{
					Name: b.Element,
				},
			},
		},
	}

	node.Decs.Before = dst.EmptyLine
	node.Decs.After = dst.EmptyLine

	for _, text := range b.Comments {
		node.Decs.Start.Append(text)
	}

	node.Decs.Start.Append(commentf("stride:generate"))

	return node
}

// revise
func camelize(text string) string {
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

func sorted(m map[string]string) []string {
	keys := []string{}

	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

func commentf(text string, args ...interface{}) string {
	text = fmt.Sprintf(text, args...)
	return fmt.Sprintf("// %s", text)
}

func pointer(text string) string {
	return fmt.Sprintf("*%s", text)
}
