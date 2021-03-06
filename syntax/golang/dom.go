package golang

import (
	"bytes"
	"fmt"
	"go/token"
	"io"
	"os"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/phogolabs/stride/codedom"
	"github.com/phogolabs/stride/inflect"
	"golang.org/x/tools/imports"
)

const (
	docConst = "// %s is a %q constant auto-generated from OpenAPI spec"
	docType  = "// %s is a type auto-generated from OpenAPI spec"
)

//go:generate counterfeiter -fake-name Writer -o ../../fake/writer.go . Writer

// Writer represents a writer
type Writer io.Writer

// Node returns the node
type Node interface {
	Node() *dst.GenDecl
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
			Decls: []dst.Decl{},
		},
	}
}

// ReadFile reads the file from reader
func ReadFile(name string, reader io.Reader) (*File, error) {
	node, err := decorator.Parse(reader)
	if err != nil {
		return nil, err
	}

	return &File{
		name: name,
		node: node,
	}, nil
}

// OpenFile opens a file
func OpenFile(name string) (*File, error) {
	reader, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return ReadFile(reader.Name(), reader)
}

// Name returns the name of the file
func (f *File) Name() string {
	return f.name
}

// Node returns the node
func (f *File) Node() *dst.File {
	return f.node
}

// AddImport adds an import
func (f *File) AddImport(name string) {
	if name == "" {
		return
	}

	name = fmt.Sprintf("%q", name)
	container := f.container()

	for _, spec := range container.Specs {
		if pkg, ok := spec.(*dst.ImportSpec); ok {
			if strings.EqualFold(pkg.Path.Value, name) {
				return
			}
		}
	}

	spec := &dst.ImportSpec{
		Path: &dst.BasicLit{
			Kind:  token.STRING,
			Value: name,
		},
	}

	container.Specs = append(container.Specs, spec)
}

// AddFunction adds a function
func (f *File) AddFunction(data string) error {
	buffer := &bytes.Buffer{}

	fmt.Fprintln(buffer, "package body")
	fmt.Fprintln(buffer, data)

	file, err := decorator.Parse(buffer.String())
	if err != nil {
		return err
	}

	for _, decl := range file.Decls {
		if node, ok := decl.(*dst.FuncDecl); ok {
			f.node.Decls = append(f.node.Decls, node)
		}
	}

	return nil
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
	buffer := &bytes.Buffer{}

	if err := decorator.Fprint(buffer, f.node); err != nil {
		return 0, err
	}

	data, err := imports.Process(f.name, buffer.Bytes(), nil)
	if err != nil {
		return 0, err
	}

	n, err := w.Write(data)
	return int64(n), err
}

// ReadFrom reads the content from a file
func (f *File) ReadFrom(r io.Reader) (int64, error) {
	node, err := decorator.Parse(r)
	if err != nil {
		return 0, err
	}

	f.node = node
	return 0, nil
}

// AddNode adds a node to the file
func (f *File) AddNode(node Node) {
	f.node.Decls = append(f.node.Decls, node.Node())
}

func (f *File) container() *dst.GenDecl {
	if count := len(f.node.Decls); count > 0 {
		if declaration, ok := f.node.Decls[0].(*dst.GenDecl); ok {
			if declaration.Tok == token.IMPORT {
				return declaration
			}
		}
	}

	container := &dst.GenDecl{
		Tok:   token.IMPORT,
		Specs: []dst.Spec{},
	}

	f.node.Decls = append([]dst.Decl{container}, f.node.Decls...)
	return container
}

// StructType builds a struct
type StructType struct {
	node *dst.GenDecl
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
					Name: inflect.Camelize(name),
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
	node.Decs.Start.Append(fmt.Sprintf(docType, inflect.Camelize(name)))
	node.Decs.Start.Append(AnnotationGenerate.Key(name))

	return &StructType{
		node: node,
	}
}

// Node returns the node
func (b *StructType) Node() *dst.GenDecl {
	return b.node
}

// Name returns the type name
func (b *StructType) Name() string {
	return b.node.Specs[0].(*dst.TypeSpec).Name.Name
}

// Commentf adds a comment
func (b *StructType) Commentf(pattern string, args ...interface{}) {
	commentf(&b.node.Decs.Start, pattern, args...)
}

// HasFields return true if the struct has fields
func (b *StructType) HasFields() bool {
	spec := b.node.Specs[0].(*dst.TypeSpec).Type.(*dst.StructType)
	return len(spec.Fields.List) > 0
}

// AddField defines a field
func (b *StructType) AddField(name, kind string, tags ...*codedom.TagDescriptor) {
	field := property(inflect.Camelize(name), kind)
	field.Decs.Before = dst.NewLine
	field.Decs.After = dst.NewLine
	field.Decs.Start.Append(AnnotationGenerate.Key(name))

	if tag := codedom.TagDescriptorCollection(tags).String(); tag != "" {
		field.Tag = &dst.BasicLit{
			Kind:  token.STRING,
			Value: tag,
		}
	}

	spec := b.node.Specs[0].(*dst.TypeSpec).Type.(*dst.StructType)
	spec.Fields.List = append(spec.Fields.List, field)
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
					Name: inflect.Camelize(name),
				},
			},
		},
	}

	// formatting
	node.Decs.Before = dst.EmptyLine
	node.Decs.After = dst.EmptyLine
	// comments
	node.Decs.Start.Append(fmt.Sprintf(docType, inflect.Camelize(name)))
	node.Decs.Start.Append(AnnotationGenerate.Key(name))

	return &LiteralType{
		node: node,
	}
}

// Name returns the type name
func (b *LiteralType) Name() string {
	return b.node.Specs[0].(*dst.TypeSpec).Name.Name
}

// Node returns the node
func (b *LiteralType) Node() *dst.GenDecl {
	return b.node
}

// Commentf adds a comment
func (b *LiteralType) Commentf(pattern string, args ...interface{}) {
	commentf(&b.node.Decs.Start, pattern, args...)
}

// Element sets the element
func (b *LiteralType) Element(name string) *LiteralType {
	b.node.Specs[0].(*dst.TypeSpec).Type = &dst.Ident{
		Name: name,
	}

	return b
}

// ConstBlockType builds a var block type
type ConstBlockType struct {
	node *dst.GenDecl
}

// NewConstBlockType creates a new ConstBlcokType
func NewConstBlockType() *ConstBlockType {
	node := &dst.GenDecl{
		Tok:   token.CONST,
		Specs: []dst.Spec{},
	}

	// formatting
	node.Decs.Before = dst.EmptyLine
	node.Decs.After = dst.EmptyLine

	return &ConstBlockType{
		node: node,
	}
}

// Node returns the node
func (b *ConstBlockType) Node() *dst.GenDecl {
	return b.node
}

// Commentf adds a comment
func (b *ConstBlockType) Commentf(pattern string, args ...interface{}) {
	commentf(&b.node.Decs.Start, pattern, args...)
}

// AddConst defines a var
func (b *ConstBlockType) AddConst(name, kind, value string) {
	field := property(name, kind)

	spec := &dst.ValueSpec{
		Names: field.Names,
		Type:  field.Type,
	}

	spec.Decs.Before = dst.NewLine
	spec.Decs.After = dst.EmptyLine
	spec.Decs.Start.Append(fmt.Sprintf(docConst, inflect.Camelize(name), value))
	spec.Decs.Start.Append(AnnotationGenerate.Key(name))

	if value != "" {
		spec.Values = []dst.Expr{
			&dst.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("%q", value),
			},
		}
	}

	b.node.Specs = append(b.node.Specs, spec)
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
					Name: inflect.Camelize(name),
				},
				Type: &dst.ArrayType{},
			},
		},
	}

	// formatting
	node.Decs.Before = dst.EmptyLine
	node.Decs.After = dst.EmptyLine
	// comments
	node.Decs.Start.Append(fmt.Sprintf(docType, inflect.Camelize(name)))
	node.Decs.Start.Append(AnnotationGenerate.Key(name))

	return &ArrayType{
		node: node,
	}
}

// Node returns the node
func (b *ArrayType) Node() *dst.GenDecl {
	return b.node
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
		Name: name,
	}

	return b
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
			{
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
