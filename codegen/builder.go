package codegen

import (
	"bytes"
	"fmt"
	"go/token"
	"strings"

	"github.com/dave/dst"
	"github.com/fatih/structtag"
	"github.com/go-openapi/inflect"
)

// Builder builds a node from the code
type Builder interface {
	Build() []dst.Decl
	Name() string
	Commentf(string, ...interface{})
}

// File represents the generated file
type File struct {
	Name    string
	Content *dst.File
}

// FileBuilder builds a file
type FileBuilder struct {
	node *dst.File
}

// NewFileBuilder creates a new file
func NewFileBuilder(name string) *FileBuilder {
	node := &dst.File{
		Name: &dst.Ident{
			Name: name,
		},
	}

	return &FileBuilder{
		node: node,
	}
}

// Build builds the file
func (b *FileBuilder) Build(name string) *File {
	return &File{
		Name:    name,
		Content: b.node,
	}
}

// Type returns a struct type
func (b *FileBuilder) Type(name string) *StructTypeBuilder {
	builder := NewStructTypeBuilder(camelize(name))
	b.node.Decls = append(b.node.Decls, builder.node)
	return builder
}

// Literal returns a literal type
func (b *FileBuilder) Literal(name string) *LiteralTypeBuilder {
	builder := NewLiteralTypeBuilder(camelize(name))
	b.node.Decls = append(b.node.Decls, builder.node)
	return builder
}

// Array returns a array type
func (b *FileBuilder) Array(name string) *ArrayTypeBuilder {
	builder := NewArrayTypeBuilder(camelize(name))
	b.node.Decls = append(b.node.Decls, builder.node)
	return builder
}

// Tag represents a tag
type Tag struct {
	Key     string
	Name    string
	Options []string
}

// Tags represents a field tag list
type Tags []*Tag

func (tags Tags) String() string {
	builder := &structtag.Tags{}

	for _, descriptor := range tags {
		builder.Set(&structtag.Tag{
			Key:     descriptor.Key,
			Name:    descriptor.Name,
			Options: descriptor.Options,
		})
	}

	if value := builder.String(); value != "" {
		return fmt.Sprintf("`%s`", value)
	}

	return ""
}

var _ Builder = &StructTypeBuilder{}

// StructTypeBuilder builds a struct
type StructTypeBuilder struct {
	node    *dst.GenDecl
	methods []*dst.FuncDecl
}

// NewStructTypeBuilder creates a new struct type builder
func NewStructTypeBuilder(name string) *StructTypeBuilder {
	node := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: &dst.Ident{
					Name: camelize(name),
				},
				Type: &dst.StructType{
					Fields: &dst.FieldList{},
				},
			},
		},
	}

	// formatting
	node.Decs.Before = dst.EmptyLine
	node.Decs.After = dst.EmptyLine
	node.Decs.Start.Append("// stride:generate")

	return &StructTypeBuilder{
		node: node,
	}
}

// Name returns the type name
func (b *StructTypeBuilder) Name() string {
	return b.node.Specs[0].(*dst.TypeSpec).Name.Name
}

// Commentf adds a comment
func (b *StructTypeBuilder) Commentf(pattern string, args ...interface{}) {
	commentf(&b.node.Decs.Start, pattern, args...)
}

// Field defines a field
func (b *StructTypeBuilder) Field(name, kind string, tags ...*Tag) {
	field := property(camelize(name), kind)

	if options := Tags(tags).String(); options != "" {
		field.Tag = &dst.BasicLit{
			Kind:  token.STRING,
			Value: options,
		}
	}

	spec := b.node.Specs[0].(*dst.TypeSpec).Type.(*dst.StructType)
	spec.Fields.List = append(spec.Fields.List, field)
}

// Method returns a struct method
func (b *StructTypeBuilder) Method(name string) *MethodTypeBuilder {
	builder := NewMethodTypeBuilder(name).Receiver("x", pointer(b.Name()))
	b.methods = append(b.methods, builder.node)

	return builder
}

// Build builds the type
func (b *StructTypeBuilder) Build() []dst.Decl {
	// tree
	tree := []dst.Decl{}
	tree = append(tree, b.node)

	for _, method := range b.methods {
		tree = append(tree, method)
	}

	return tree
}

var _ Builder = &LiteralTypeBuilder{}

// LiteralTypeBuilder builds a literal type
type LiteralTypeBuilder struct {
	node *dst.GenDecl
}

// NewLiteralTypeBuilder creates a new literal type builder
func NewLiteralTypeBuilder(name string) *LiteralTypeBuilder {
	node := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: &dst.Ident{
					Name: camelize(name),
				},
			},
		},
	}

	// formatting
	node.Decs.Before = dst.EmptyLine
	node.Decs.After = dst.EmptyLine
	node.Decs.Start.Append("// stride:generate")

	return &LiteralTypeBuilder{
		node: node,
	}
}

// Name returns the type name
func (b *LiteralTypeBuilder) Name() string {
	return b.node.Specs[0].(*dst.TypeSpec).Name.Name
}

// Commentf adds a comment
func (b *LiteralTypeBuilder) Commentf(pattern string, args ...interface{}) {
	commentf(&b.node.Decs.Start, pattern, args...)
}

// Element sets the element
func (b *LiteralTypeBuilder) Element(name string) *LiteralTypeBuilder {
	b.node.Specs[0].(*dst.TypeSpec).Type = &dst.Ident{
		Name: camelize(name),
	}

	return b
}

// Build builds the type
func (b *LiteralTypeBuilder) Build() []dst.Decl {
	return []dst.Decl{b.node}
}

var _ Builder = &ArrayTypeBuilder{}

// ArrayTypeBuilder builds an array type
type ArrayTypeBuilder struct {
	node *dst.GenDecl
}

// NewArrayTypeBuilder creates a new ArrayTypeBuilder
func NewArrayTypeBuilder(name string) *ArrayTypeBuilder {
	node := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: &dst.Ident{
					Name: camelize(name),
				},
				Type: &dst.ArrayType{},
			},
		},
	}

	// formatting
	node.Decs.Before = dst.EmptyLine
	node.Decs.After = dst.EmptyLine
	node.Decs.Start.Append("// stride:generate")

	return &ArrayTypeBuilder{
		node: node,
	}
}

// Name returns the type name
func (b *ArrayTypeBuilder) Name() string {
	return b.node.Specs[0].(*dst.TypeSpec).Name.Name
}

// Commentf adds a comment
func (b *ArrayTypeBuilder) Commentf(pattern string, args ...interface{}) {
	commentf(&b.node.Decs.Start, pattern, args...)
}

// Element sets the element
func (b *ArrayTypeBuilder) Element(name string) *ArrayTypeBuilder {
	b.node.Specs[0].(*dst.TypeSpec).Type.(*dst.ArrayType).Elt = &dst.Ident{
		Name: camelize(name),
	}

	return b
}

// Build builds the type
func (b *ArrayTypeBuilder) Build() []dst.Decl {
	// formatting
	b.node.Decs.Before = dst.EmptyLine
	b.node.Decs.After = dst.EmptyLine
	b.node.Decs.Start.Append("// stride:generate")

	return []dst.Decl{b.node}
}

var _ Builder = &MethodTypeBuilder{}

// MethodTypeBuilder builds a method
type MethodTypeBuilder struct {
	node *dst.FuncDecl
}

// NewMethodTypeBuilder creates a new method type builder
func NewMethodTypeBuilder(name string) *MethodTypeBuilder {
	node := &dst.FuncDecl{
		// function receiver
		Recv: &dst.FieldList{
			List: []*dst.Field{},
		},
		// function name
		Name: &dst.Ident{
			Name: camelize(name),
		},
		// function param
		Type: &dst.FuncType{
			Params: &dst.FieldList{
				List: []*dst.Field{},
			},
			Results: &dst.FieldList{
				List: []*dst.Field{},
			},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{},
		},
	}

	node.Decs.Before = dst.EmptyLine
	node.Decs.After = dst.EmptyLine
	node.Decs.Start.Append("// stride:generate")

	return &MethodTypeBuilder{
		node: node,
	}
}

// Name returns the type name
func (b *MethodTypeBuilder) Name() string {
	return b.node.Name.Name
}

// Commentf adds a comment
func (b *MethodTypeBuilder) Commentf(pattern string, args ...interface{}) {
	commentf(&b.node.Decs.Start, pattern, args...)
}

// Receiver creates a return parameter
func (b *MethodTypeBuilder) Receiver(name, kind string) *MethodTypeBuilder {
	field := property(name, kind)
	b.node.Recv.List = append(b.node.Recv.List, field)
	return b
}

// Param creates a parameter
func (b *MethodTypeBuilder) Param(name, kind string) *MethodTypeBuilder {
	field := property(name, kind)
	b.node.Type.Params.List = append(b.node.Type.Params.List, field)
	return b
}

// Return creates a return parameter
func (b *MethodTypeBuilder) Return(kind string) *MethodTypeBuilder {
	field := property("", kind)
	b.node.Type.Results.List = append(b.node.Type.Results.List, field)
	return b
}

// Build builds the method
func (b *MethodTypeBuilder) Build() []dst.Decl {
	return []dst.Decl{b.node}
}

func property(name, kind string) *dst.Field {
	field := &dst.Field{
		Type: &dst.Ident{
			Name: kind,
		},
	}

	if name != "" {
		field.Names = []*dst.Ident{
			&dst.Ident{
				Name: name,
			},
		}
	}

	return field
}

func camelize(text string) string {
	const (
		separator = "-"
		suffix    = "id"
	)

	var (
		parts  = strings.Split(inflect.Dasherize(text), separator)
		buffer = &bytes.Buffer{}
	)

	for index, part := range parts {
		if index > 0 {
			buffer.WriteString(separator)
		}

		if strings.EqualFold(part, suffix) {
			part = strings.ToUpper(part)
		}

		buffer.WriteString(part)
	}

	return inflect.Camelize(buffer.String())
}

func commentf(decorator *dst.Decorations, text string, args ...interface{}) {
	if text == "" {
		return
	}

	var (
		comments = decorator.All()
		index    = len(comments) - 1
	)

	text = fmt.Sprintf(text, args...)
	text = fmt.Sprintf("// %s", text)

	decorator.Replace(comments[:index]...)
	decorator.Append(text)
	decorator.Append(comments[index:]...)
}

func pointer(text string) string {
	return fmt.Sprintf("*%s", text)
}
