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
	"github.com/dave/dst/dstutil"
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

// Finder finds attributes
type Finder struct {
	node dst.Node
}

func (f *Finder) extract(text, prefix string) string {
	text = strings.TrimPrefix(text, "//")
	text = strings.TrimSpace(text)

	if strings.HasPrefix(text, prefix) {
		text = strings.TrimPrefix(text, prefix)
		text = strings.TrimSpace(text)

		return text
	}

	return ""
}

// Find finds the special comment
func (f *Finder) Find(prefix string) string {
	for _, line := range f.node.Decorations().Start.All() {
		if name := f.extract(line, prefix); name != "" {
			return name
		}
	}

	for _, line := range f.node.Decorations().End.All() {
		if name := f.extract(line, prefix); name != "" {
			return name
		}
	}

	return ""
}

// Merge merges the files
func (f *File) Merge(source *File) error {
	merge := func(cursor *dstutil.Cursor) bool {
		if node := cursor.Node(); node != nil {
			finder := Finder{node: node}

			if kind := finder.Find("stride:struct"); kind != "" {
				fmt.Println("STRUCT", cursor.Name(), kind)
				return true
			}

			if kind := finder.Find("stride:field"); kind != "" {
				fmt.Println("FIELD", cursor.Name(), kind)
				return false
			}

			if kind := finder.Find("stride:function"); kind != "" {
				fmt.Println("FUNC", cursor.Name(), kind)
				return true
			}

			if kind := finder.Find("stride:block"); kind != "" {
				fmt.Println("BLOCK", cursor.Name(), kind)
				return false
			}
		}

		return true
	}

	dstutil.Apply(f.node, merge, nil)
	return nil
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
	fields.Decs.Opening.Append("// stride:field open")
	fields.Decs.Opening.Append("// TODO: Please add your implementation here")
	fields.Decs.Opening.Append("// stride:field close")

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
	node.Decs.Start.Append(fmt.Sprintf("// stride:struct %s", dasherize(name)))

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
	node.Decs.Start.Append(fmt.Sprintf("// stride:literal %s", dasherize(name)))

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
	node.Decs.Start.Append(fmt.Sprintf("// stride:array %s", dasherize(name)))

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
	node.Decs.Start.Append(fmt.Sprintf("// stride:function %s", dasherize(name)))

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
	return b
}

// AddParam creates a parameter
func (b *FunctionType) AddParam(name, kind string) *FunctionType {
	field := property(name, kind)
	b.node.Type.Params.List = append(b.node.Type.Params.List, field)
	return b
}

// Block sets the block
func (b *FunctionType) Block(content string, args ...interface{}) *FunctionType {
	content = fmt.Sprintf(content, args...)
	buffer := &bytes.Buffer{}

	fmt.Fprintln(buffer, "package block")
	fmt.Fprintln(buffer, "func block() {")
	fmt.Fprintln(buffer, content)
	fmt.Fprintln(buffer, "}")

	file, err := decorator.Parse(buffer.String())
	if err != nil {
		panic(err)
	}

	if node, ok := file.Decls[0].(*dst.FuncDecl); ok {
		b.node.Body = node.Body
	}

	return b
}

// AddReturn creates a return parameter
func (b *FunctionType) AddReturn(kind string) *FunctionType {
	field := property("", kind)
	b.node.Type.Results.List = append(b.node.Type.Results.List, field)
	return b
}

// BlockWriter writes the block
type BlockWriter struct {
	buffer *bytes.Buffer
}

// NewBlockWriter creates a new block writera
func NewBlockWriter() *BlockWriter {
	return &BlockWriter{
		buffer: &bytes.Buffer{},
	}
}

// Write the block
func (b *BlockWriter) Write(content string, args ...interface{}) {
	if b.buffer.Len() > 0 {
		fmt.Fprintln(b.buffer)
	}

	fmt.Fprintf(b.buffer, content, args...)
}

// String returns the block as string
func (b *BlockWriter) String() string {
	return b.buffer.String()
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
