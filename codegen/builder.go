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
	Build() []dst.Decl
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
	for _, builder := range b.builders {
		if value, ok := builder.(*StructTypeBuilder); ok {
			if strings.EqualFold(value.Name, name) {
				return value
			}
		}
	}

	builder := &StructTypeBuilder{
		Name: name,
	}

	b.builders = append(b.builders, builder)
	return builder
}

// Literal returns a literal type
func (b *FileBuilder) Literal(name string) *LiteralTypeBuilder {
	for _, builder := range b.builders {
		if value, ok := builder.(*LiteralTypeBuilder); ok {
			if strings.EqualFold(value.Name, name) {
				return value
			}
		}
	}

	builder := &LiteralTypeBuilder{
		Name: name,
	}

	b.builders = append(b.builders, builder)
	return builder
}

// Array returns a array type
func (b *FileBuilder) Array(name string) *ArrayTypeBuilder {
	for _, builder := range b.builders {
		if value, ok := builder.(*ArrayTypeBuilder); ok {
			if strings.EqualFold(value.Name, name) {
				return value
			}
		}
	}

	builder := &ArrayTypeBuilder{
		Name: name,
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
	Name string

	// private
	comments []string
	fields   []*Field
	builders []Builder
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
	field := &Field{
		Name: name,
		Type: kind,
		Tags: tags,
	}

	b.fields = append(b.fields, field)
}

// Method returns a struct method
func (b *StructTypeBuilder) Method(name string) *MethodTypeBuilder {
	for _, builder := range b.builders {
		if value, ok := builder.(*MethodTypeBuilder); ok {
			if strings.EqualFold(value.Name, name) {
				return value
			}
		}
	}

	builder := &MethodTypeBuilder{
		Name: name,
		receiver: &Param{
			Name: "x",
			Type: "*" + b.Name,
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
					Name: b.Name,
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
	Name string
	// private
	element  string
	comments []string
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
	b.element = name
}

// Build builds the type
func (b *LiteralTypeBuilder) Build() []dst.Decl {
	node := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: &dst.Ident{
					Name: b.Name,
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

var _ Builder = &ArrayTypeBuilder{}

// ArrayTypeBuilder builds an array type
type ArrayTypeBuilder struct {
	Name string
	// private
	element  string
	comments []string
}

// Element sets the element
func (b *ArrayTypeBuilder) Element(name string) {
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
					Name: b.Name,
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

// Param represents the method parameters
type Param struct {
	Name string
	Type string
}

var _ Builder = &MethodTypeBuilder{}

// MethodTypeBuilder builds a method
type MethodTypeBuilder struct {
	Name string
	// private
	receiver   *Param
	parameters []*Param
	comments   []string
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

// Build builds the method
func (b *MethodTypeBuilder) Build() []dst.Decl {
	node := &dst.FuncDecl{
		// function receiver
		Recv: &dst.FieldList{
			List: []*dst.Field{},
		},
		// function name
		Name: &dst.Ident{
			Name: b.Name,
		},
		// function param
		Type: &dst.FuncType{
			Params: &dst.FieldList{
				List: []*dst.Field{},
			},
		},
		Body: &dst.BlockStmt{},
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

// Method represents a method
type Method struct {
	Name       string
	Receiver   string
	Parameters []string
}

// BlockBuilder build blocks
type BlockBuilder struct {
	stmt []dst.Stmt
}

// Call calls a method
func (b *BlockBuilder) Call(m *Method) {
	var (
		fun  dst.Expr
		args []dst.Expr
	)

	if receiver := m.Receiver; receiver == "" {
		fun = &dst.Ident{
			Name: m.Name,
		}
	} else {
		fun = &dst.SelectorExpr{
			X: &dst.Ident{
				Name: receiver,
			},
			Sel: &dst.Ident{
				Name: m.Name,
			},
		}
	}

	for _, param := range m.Parameters {
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

	b.stmt = append(b.stmt, stmt)
}

// Build builds the block
func (b *BlockBuilder) Build() *dst.BlockStmt {
	return &dst.BlockStmt{
		List: b.stmt,
	}
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
