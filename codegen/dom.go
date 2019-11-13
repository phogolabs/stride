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
	"github.com/go-openapi/inflect"
)

// Builder builds a node from the code
type Builder interface {
	Name() string
	Commentf(string, ...interface{})
}

// File represents the generated file
type File struct {
	name string
	node *dst.File
}

// NewFile creates a new file
func NewFile(name string) *File {
	return &File{
		name: name,
		node: &dst.File{
			Name: &dst.Ident{
				Name: "service",
			},
		},
	}
}

// Open opens a file
func Open(name string) (*File, error) {
	reader, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	node, err := decorator.Parse(reader)
	if err != nil {
		return nil, err
	}

	return &File{
		name: name,
		node: node,
	}, nil
}

// Name returns the name of the file
func (f *File) Name() string {
	return f.name
}

// Merge merges the files
func (f *File) Merge(source *File) error {
	merger := &Merger{
		Target: f,
		Source: source,
	}

	return merger.Merge()
}

// Sync syncs the content to the file system
func (f *File) Sync() error {
	w, err := os.Create(f.name)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = f.WriteTo(w)
	return err
}

// WriteTo writes to a file
func (f *File) WriteTo(w io.Writer) (int64, error) {
	if err := decorator.Fprint(w, f.node); err != nil {
		return 0, err
	}

	return 0, nil
}

// Struct returns a struct type
func (f *File) Struct(name string) *StructType {
	builder := NewStructType(name)
	builder.file = f.node

	f.node.Decls = append(f.node.Decls, builder.node)
	return builder
}

// Literal returns a literal type
func (f *File) Literal(name string) *LiteralType {
	builder := NewLiteralType(name)

	f.node.Decls = append(f.node.Decls, builder.node)
	return builder
}

// Array returns a array type
func (f *File) Array(name string) *ArrayType {
	builder := NewArrayType(name)

	f.node.Decls = append(f.node.Decls, builder.node)
	return builder
}

// StructType builds a struct
type StructType struct {
	node *dst.GenDecl
	file *dst.File
}

// NewStructType creates a new struct type builder
func NewStructType(name string) *StructType {
	fields := &dst.FieldList{}
	fields.Closing = true

	fields.Decs.Before = dst.NewLine
	fields.Decs.After = dst.NewLine

	node := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: &dst.Ident{
					Name: camelize(name),
				},
				Type: &dst.StructType{
					Fields: fields,
				},
			},
		},
	}

	// formatting
	node.Decs.Before = dst.EmptyLine
	node.Decs.After = dst.EmptyLine
	// comments
	node.Decs.Start.Append(fmt.Sprintf("// %s is a struct type auto-generated from OpenAPI spec", camelize(name)))
	node.Decs.Start.Append(fmt.Sprintf("// stride:generate:struct %s", dasherize(name)))

	return &StructType{
		node: node,
	}
}

// Name returns the type name
func (b *StructType) Name() string {
	return b.node.Specs[0].(*dst.TypeSpec).Name.Name
}

// Commentf adds a comment
func (b *StructType) Commentf(pattern string, args ...interface{}) {
	commentf(&b.node.Decs.Start, pattern, args...)
}

// AddField defines a field
func (b *StructType) AddField(name, kind string, tags ...*TagDescriptor) {
	field := property(camelize(name), kind)
	field.Decs.Before = dst.NewLine
	field.Decs.After = dst.NewLine
	field.Decs.Start.Append("// stride:generate:field " + dasherize(name))

	if tag := TagDescriptorCollection(tags).String(); tag != "" {
		field.Tag = &dst.BasicLit{
			Kind:  token.STRING,
			Value: tag,
		}
	}

	spec := b.node.Specs[0].(*dst.TypeSpec).Type.(*dst.StructType)
	spec.Fields.List = append(spec.Fields.List, field)
}

// Function returns a struct method
func (b *StructType) Function(name string) *FunctionType {
	builder := NewFunctionType(name).AddReceiver("x", pointer(b.Name()))
	b.file.Decls = append(b.file.Decls, builder.node)
	return builder
}

// LiteralType builds a literal type
type LiteralType struct {
	node *dst.GenDecl
}

// NewLiteralType creates a new literal type builder
func NewLiteralType(name string) *LiteralType {
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
	// comments
	node.Decs.Start.Append(fmt.Sprintf("// %s is a literal type auto-generated from OpenAPI spec", camelize(name)))
	node.Decs.Start.Append(fmt.Sprintf("// stride:generate:literal %s", dasherize(name)))

	return &LiteralType{
		node: node,
	}
}

// Name returns the type name
func (b *LiteralType) Name() string {
	return b.node.Specs[0].(*dst.TypeSpec).Name.Name
}

// Commentf adds a comment
func (b *LiteralType) Commentf(pattern string, args ...interface{}) {
	commentf(&b.node.Decs.Start, pattern, args...)
}

// Element sets the element
func (b *LiteralType) Element(name string) *LiteralType {
	b.node.Specs[0].(*dst.TypeSpec).Type = &dst.Ident{
		Name: camelize(name),
	}

	return b
}

// ArrayType builds an array type
type ArrayType struct {
	node *dst.GenDecl
}

// NewArrayType creates a new ArrayType
func NewArrayType(name string) *ArrayType {
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
	// comments
	node.Decs.Start.Append(fmt.Sprintf("// %s is a array type auto-generated from OpenAPI spec", camelize(name)))
	node.Decs.Start.Append(fmt.Sprintf("// stride:generate:array %s", dasherize(name)))

	return &ArrayType{
		node: node,
	}
}

// Name returns the type name
func (b *ArrayType) Name() string {
	return b.node.Specs[0].(*dst.TypeSpec).Name.Name
}

// Commentf adds a comment
func (b *ArrayType) Commentf(pattern string, args ...interface{}) {
	commentf(&b.node.Decs.Start, pattern, args...)
}

// Element sets the element
func (b *ArrayType) Element(name string) *ArrayType {
	b.node.Specs[0].(*dst.TypeSpec).Type.(*dst.ArrayType).Elt = &dst.Ident{
		Name: camelize(name),
	}

	return b
}

var _ Builder = &FunctionType{}

// FunctionType builds a method
type FunctionType struct {
	node *dst.FuncDecl
}

// NewFunctionType creates a new method type builder
func NewFunctionType(name string) *FunctionType {
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

	// formatting
	node.Decs.Before = dst.EmptyLine
	node.Decs.After = dst.EmptyLine
	// comments
	node.Decs.Start.Append(fmt.Sprintf("// stride:generate:function %s", dasherize(name)))

	return &FunctionType{
		node: node,
	}
}

// Name returns the type name
func (b *FunctionType) Name() string {
	return b.node.Name.Name
}

// Commentf adds a comment
func (b *FunctionType) Commentf(pattern string, args ...interface{}) {
	commentf(&b.node.Decs.Start, pattern, args...)
}

// AddReceiver creates a return parameter
func (b *FunctionType) AddReceiver(name, kind string) *FunctionType {
	field := property(name, kind)
	b.node.Recv.List = append(b.node.Recv.List, field)

	key := fmt.Sprintf("%s:%s", dasherize(kind), dasherize(b.Name()))

	var (
		comments = b.node.Decs.Start.All()
		index    = len(comments) - 1
	)

	comments[index] = fmt.Sprintf("// stride:generate:function %s", key)
	b.node.Decs.Start.Replace(comments...)

	return b
}

// AddParam creates a parameter
func (b *FunctionType) AddParam(name, kind string) *FunctionType {
	field := property(name, kind)
	b.node.Type.Params.List = append(b.node.Type.Params.List, field)
	return b
}

// Body sets the body
func (b *FunctionType) Body() *BlockType {
	block := &BlockType{
		node:   b.node.Body,
		buffer: &bytes.Buffer{},
	}
	return block
}

// AddReturn creates a return parameter
func (b *FunctionType) AddReturn(kind string) *FunctionType {
	field := property("", kind)
	b.node.Type.Results.List = append(b.node.Type.Results.List, field)
	return b
}

// BlockType represents a block type
type BlockType struct {
	node   *dst.BlockStmt
	buffer *bytes.Buffer
}

// Write the block
func (b *BlockType) Write(content string, args ...interface{}) {
	fmt.Fprintf(b.buffer, content, args...)
	fmt.Fprintln(b.buffer)
}

// Build builds the block
func (b *BlockType) Build() error {
	const newline = "\n"

	var (
		content = b.buffer.String()
		buffer  = &bytes.Buffer{}
	)

	fmt.Fprintln(buffer, "package body")
	fmt.Fprintln(buffer, "func body() {")
	fmt.Fprintln(buffer, strings.TrimSuffix(content, newline))
	fmt.Fprintln(buffer, "}")

	file, err := decorator.Parse(buffer.String())
	if err != nil {
		return err
	}

	if node, ok := file.Decls[0].(*dst.FuncDecl); ok {
		b.node.List = append(b.node.List, node.Body.List...)
	}

	b.buffer.Reset()
	return nil
}

// WriteComment writes the body block comment
func (b *BlockType) WriteComment() {
	fmt.Fprintln(b.buffer, "// stride:define:block:start body")
	fmt.Fprintln(b.buffer, "// NOTE: You can your code within the comment block")
	fmt.Fprintln(b.buffer, "// stride:define:block:end body")
}

func kind(field *dst.Field) string {
	kind := field.Type

	if starExpr, ok := kind.(*dst.StarExpr); ok {
		kind = starExpr.X
	}

	if selectorExpr, ok := kind.(*dst.SelectorExpr); ok {
		if ident, ok := selectorExpr.X.(*dst.Ident); ok {
			return fmt.Sprintf("%s.%s", ident.Name, selectorExpr.Sel.Name)
		}
	}

	if ident, ok := kind.(*dst.Ident); ok {
		return ident.Name
	}

	return ""
}

func property(name, kind string) *dst.Field {
	const star = "*"

	field := &dst.Field{}

	// set field name
	if name != "" {
		field.Names = []*dst.Ident{
			&dst.Ident{
				Name: name,
			},
		}
	}

	getType := func() dst.Expr {
		if parts := strings.Split(kind, "."); len(parts) == 2 {
			return &dst.SelectorExpr{
				X: &dst.Ident{
					Name: parts[0],
				},
				Sel: &dst.Ident{
					Name: parts[1],
				},
			}
		}

		return &dst.Ident{
			Name: kind,
		}
	}

	if strings.HasPrefix(kind, star) {
		kind = strings.TrimPrefix(kind, star)

		field.Type = &dst.StarExpr{
			X: getType(),
		}

		return field
	}

	field.Type = getType()

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

func dasherize(text string) string {
	text = strings.TrimPrefix(text, "*")
	return inflect.Dasherize(text)
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

	decorator.Clear()
	decorator.Append(comments[:index]...)
	decorator.Append(text)
	decorator.Append(comments[index:]...)
}

func pointer(text string) string {
	return fmt.Sprintf("*%s", text)
}
