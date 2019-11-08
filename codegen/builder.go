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
	Package string
	// private
	builders []Builder
}

// Build builds the file
func (b *FileBuilder) Build() *dst.File {
	root := &dst.File{
		Name: &dst.Ident{
			Name: b.Package,
		},
	}

	for _, builder := range b.builders {
		root.Decls = append(root.Decls, builder.Build()...)
	}

	return root
}

// Type returns a struct type
func (b *FileBuilder) Type(name string) *StructTypeBuilder {
	name = camelize(name)

	for _, builder := range b.builders {
		if value, ok := builder.(*StructTypeBuilder); ok {
			if strings.EqualFold(value.name, name) {
				return value
			}
		}
	}

	builder := &StructTypeBuilder{
		name: name,
	}

	b.builders = append(b.builders, builder)
	return builder
}

// Literal returns a literal type
func (b *FileBuilder) Literal(name string) *LiteralTypeBuilder {
	name = camelize(name)

	for _, builder := range b.builders {
		if value, ok := builder.(*LiteralTypeBuilder); ok {
			if strings.EqualFold(value.name, name) {
				return value
			}
		}
	}

	builder := &LiteralTypeBuilder{
		name: name,
	}

	b.builders = append(b.builders, builder)
	return builder
}

// Array returns a array type
func (b *FileBuilder) Array(name string) *ArrayTypeBuilder {
	name = camelize(name)

	for _, builder := range b.builders {
		if value, ok := builder.(*ArrayTypeBuilder); ok {
			if strings.EqualFold(value.name, name) {
				return value
			}
		}
	}

	builder := &ArrayTypeBuilder{
		name: name,
	}

	b.builders = append(b.builders, builder)
	return builder
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

var _ Builder = &StructTypeBuilder{}

// StructTypeBuilder builds a struct
type StructTypeBuilder struct {
	name     string
	comments []string
	fields   []*Field
	builders []Builder
}

// Name returns the type name
func (b *StructTypeBuilder) Name() string {
	return b.name
}

// Commentf adds a comment
func (b *StructTypeBuilder) Commentf(pattern string, args ...interface{}) {
	if pattern == "" {
		return
	}

	b.comments = append(b.comments, commentf(pattern, args...))
}

// Field defines a field
func (b *StructTypeBuilder) Field(name, kind string, tags ...*Tag) {
	name = camelize(name)

	field := &Field{
		Name: name,
		Type: kind,
		Tags: tags,
	}

	b.fields = append(b.fields, field)
}

// Method returns a struct method
func (b *StructTypeBuilder) Method(name string) *MethodTypeBuilder {
	name = camelize(name)

	for _, builder := range b.builders {
		if value, ok := builder.(*MethodTypeBuilder); ok {
			if strings.EqualFold(value.name, name) {
				return value
			}
		}
	}

	builder := &MethodTypeBuilder{
		name: name,
		receiver: &Param{
			Name: "x",
			Type: "*" + b.name,
		},
	}

	b.builders = append(b.builders, builder)
	return builder
}

// Build builds the type
func (b *StructTypeBuilder) Build() []dst.Decl {
	var expr dst.Expr

	if spec := b.spec(); spec != nil {
		expr = spec
	}

	node := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: &dst.Ident{
					Name: b.name,
				},
				Type: expr,
			},
		},
	}

	node.Decs.Before = dst.EmptyLine
	node.Decs.After = dst.EmptyLine

	for _, text := range b.comments {
		node.Decs.Start.Append(text)
	}

	node.Decs.Start.Append(commentf("stride:generate"))

	tree := []dst.Decl{}
	tree = append(tree, node)

	for _, builder := range b.builders {
		tree = append(tree, builder.Build()...)
	}

	return tree
}

func (b *StructTypeBuilder) spec() dst.Expr {
	spec := &dst.StructType{
		Fields: &dst.FieldList{},
	}

	for _, descriptor := range b.fields {
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

		tags := &structtag.Tags{}

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

var _ Builder = &LiteralTypeBuilder{}

// LiteralTypeBuilder builds a literal type
type LiteralTypeBuilder struct {
	name     string
	element  string
	comments []string
}

// Name returns the type name
func (b *LiteralTypeBuilder) Name() string {
	return b.name
}

// Commentf adds a comment
func (b *LiteralTypeBuilder) Commentf(pattern string, args ...interface{}) {
	if pattern == "" {
		return
	}

	b.comments = append(b.comments, commentf(pattern, args...))
}

// Element sets the element
func (b *LiteralTypeBuilder) Element(name string) {
	name = camelize(name)
	b.element = name
}

// Build builds the type
func (b *LiteralTypeBuilder) Build() []dst.Decl {
	node := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: &dst.Ident{
					Name: b.name,
				},
				Type: &dst.Ident{
					Name: b.element,
				},
			},
		},
	}

	node.Decs.Before = dst.EmptyLine
	node.Decs.After = dst.EmptyLine

	for _, text := range b.comments {
		node.Decs.Start.Append(text)
	}

	node.Decs.Start.Append(commentf("stride:generate"))

	return []dst.Decl{node}
}

var _ Builder = &ArrayTypeBuilder{}

// ArrayTypeBuilder builds an array type
type ArrayTypeBuilder struct {
	name     string
	element  string
	comments []string
}

// Name returns the type name
func (b *ArrayTypeBuilder) Name() string {
	return b.name
}

// Element sets the element
func (b *ArrayTypeBuilder) Element(name string) {
	name = camelize(name)
	b.element = name
}

// Commentf adds a comment
func (b *ArrayTypeBuilder) Commentf(pattern string, args ...interface{}) {
	if pattern == "" {
		return
	}

	b.comments = append(b.comments, commentf(pattern, args...))
}

// Build builds the type
func (b *ArrayTypeBuilder) Build() []dst.Decl {
	node := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: &dst.Ident{
					Name: b.name,
				},
				Type: &dst.ArrayType{
					Elt: &dst.Ident{
						Name: b.element,
					},
				},
			},
		},
	}

	node.Decs.Before = dst.EmptyLine
	node.Decs.After = dst.EmptyLine

	for _, text := range b.comments {
		node.Decs.Start.Append(text)
	}

	node.Decs.Start.Append(commentf("stride:generate"))

	return []dst.Decl{node}
}

// Param represents the method parameters
type Param struct {
	Name string
	Type string
}

var _ Builder = &MethodTypeBuilder{}

// MethodTypeBuilder builds a method
type MethodTypeBuilder struct {
	name       string
	receiver   *Param
	parameters []*Param
	comments   []string
	block      *dst.BlockStmt
}

// Name returns the type name
func (b *MethodTypeBuilder) Name() string {
	return b.name
}

// Commentf adds a comment
func (b *MethodTypeBuilder) Commentf(pattern string, args ...interface{}) {
	if pattern == "" {
		return
	}

	b.comments = append(b.comments, commentf(pattern, args...))
}

// Param creates a parameter
func (b *MethodTypeBuilder) Param(name, kind string) {
	param := &Param{
		Name: name,
		Type: kind,
	}

	b.parameters = append(b.parameters, param)
}

// Block returns the block type builder
func (b *MethodTypeBuilder) Block() *BlockTypeBuilder {
	b.init()

	return &BlockTypeBuilder{
		block: b.block,
	}
}

// Build builds the method
func (b *MethodTypeBuilder) Build() []dst.Decl {
	b.init()

	node := &dst.FuncDecl{
		// function receiver
		Recv: &dst.FieldList{
			List: []*dst.Field{},
		},
		// function name
		Name: &dst.Ident{
			Name: b.name,
		},
		// function param
		Type: &dst.FuncType{
			Params: &dst.FieldList{
				List: []*dst.Field{},
			},
		},
		Body: b.block,
	}

	// receiver param
	if receiver := b.receiver; receiver != nil {
		field := b.field(receiver)
		node.Recv.List = append(node.Recv.List, field)
	}

	// function param
	for _, param := range b.parameters {
		field := b.field(param)
		node.Type.Params.List = append(node.Type.Params.List, field)
	}

	node.Decs.Before = dst.EmptyLine
	node.Decs.After = dst.EmptyLine

	for _, text := range b.comments {
		node.Decs.Start.Append(text)
	}

	node.Decs.Start.Append(commentf("stride:generate"))

	return []dst.Decl{node}
}

func (b *MethodTypeBuilder) init() {
	if b.block == nil {
		b.block = &dst.BlockStmt{}
	}
}

func (b *MethodTypeBuilder) field(param *Param) *dst.Field {
	field := &dst.Field{
		Names: []*dst.Ident{
			&dst.Ident{
				Name: param.Name,
			},
		},
		Type: &dst.Ident{
			Name: param.Type,
		},
	}

	return field
}

// BlockTypeBuilder build blocks
type BlockTypeBuilder struct {
	receiver string
	block    *dst.BlockStmt
}

// Call calls a method
func (b *BlockTypeBuilder) Call(name string, params ...string) {
	b.CallWithReceiver("", name, params...)
}

// CallWithReceiver calls a method
func (b *BlockTypeBuilder) CallWithReceiver(receiver, name string, params ...string) {
	var (
		fun  dst.Expr
		args []dst.Expr
	)

	if receiver == "" {
		fun = &dst.Ident{
			Name: name,
		}
	} else {
		fun = &dst.SelectorExpr{
			X: &dst.Ident{
				Name: receiver,
			},
			Sel: &dst.Ident{
				Name: name,
			},
		}
	}

	for _, param := range params {
		arg := &dst.Ident{
			Name: param,
		}

		args = append(args, arg)
	}

	stmt := &dst.ExprStmt{
		X: &dst.CallExpr{
			Fun:  fun,
			Args: args,
		},
	}

	b.block.List = append(b.block.List, stmt)
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

func commentf(text string, args ...interface{}) string {
	text = fmt.Sprintf(text, args...)
	return fmt.Sprintf("// %s", text)
}

func pointer(text string) string {
	return fmt.Sprintf("*%s", text)
}
