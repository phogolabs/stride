package golang

import (
	"fmt"
	"go/token"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"
	"github.com/phogolabs/stride/inflect"
)

const (
	// AnnotationGenerateStruct represents the annotation for struct
	AnnotationGenerateStruct = "stride:generate:struct"
	// AnnotationGenerateFunction represents the annotation for function
	AnnotationGenerateFunction = "stride:generate:function"
	// AnnotationGenerateField represents the annotation for fields
	AnnotationGenerateField = "stride:generate:field"
	// AnnotationDefineLiteral represents the annotation for user defined literal
	AnnotationDefineLiteral = "stride:define:literal"
	// AnnotationDefineArray represents the annotation for user defined array
	AnnotationDefineArray = "stride:define:array"
	// AnnotationDefineStruct represents the annotation for user defined struct
	AnnotationDefineStruct = "stride:define:struct"
	// AnnotationDefineField represents the annotation for user defined field
	AnnotationDefineField = "stride:define:field"
	// AnnotationDefineFunction represents the annotation for defined function
	AnnotationDefineFunction = "stride:define:function"
	// AnnotationDefineBlockStart represents the annotation for defined block start
	AnnotationDefineBlockStart = "stride:define:block:start"
	// AnnotationDefineBlockEnd represents the annotation for defined block end
	AnnotationDefineBlockEnd = "stride:define:block:end"
)

// Range represents a range
type Range struct {
	Start int
	End   int
}

// Merger merges files
type Merger struct {
	Target *File
	Source *File
}

// Merge merges the files
func (m *Merger) Merge() error {
	dstutil.Apply(m.Target.node, m.merge, nil)
	dstutil.Apply(m.Source.node, m.append, nil)

	//TODO: sort declarations

	return nil
}

func (m *Merger) merge(cursor *dstutil.Cursor) bool {
	node := cursor.Node()

	if node == nil {
		return true
	}

	// handle stride:generate:struct annotation
	if name, target := m.structType(AnnotationGenerateStruct, node); target != nil {
		if source := m.find(AnnotationGenerateStruct, name, m.Source.node); source != nil {
			var (
				left  = m.structTypeProperties(target)
				right = m.structTypeProperties(source)
			)

			for _, field := range right.List {
				name := inflect.Dasherize(field.Names[0].Name)

				// handle stride:define:field annotation
				if key := m.annotation(AnnotationDefineField, field.Decorations().Start); strings.EqualFold(name, key) {
					left.List = append(left.List, field)
				}
			}
		}
		return false
	}

	// handle stride:generate:function annotation
	if name, target := m.functionType(AnnotationGenerateFunction, node); target != nil {
		if source := m.find(AnnotationGenerateFunction, name, m.Source.node); source != nil {
			m.squash(target, source)
		}

		return false
	}

	return true
}

func (m *Merger) squash(target, source dst.Node) {
	var (
		left       = m.functionTypeBody(target)
		leftRange  = m.functionTypeBodyRange(left)
		right      = m.functionTypeBody(source)
		rightRange = m.functionTypeBodyRange(right)
	)

	if leftRange != nil && rightRange != nil {
		var (
			result = []dst.Stmt{}
			items  = right.List[rightRange.Start : rightRange.End+1]
		)

		// append top block
		for index, item := range left.List {
			if index < leftRange.Start {
				result = append(result, item)
			}
		}

		// append the range block
		result = append(result, items...)

		// append bottom block
		for index, item := range left.List {
			if index > leftRange.End {
				result = append(result, item)
			}
		}

		// sanitize comments
		m.sanitize(result)

		left.List = result
	}
}

func (m *Merger) sanitize(items []dst.Stmt) {
	// const name = "body"

	// for index := 0; index < len(items)-1; index++ {
	// 	node := items[index].Decorations()
	// 	next := items[index+1].Decorations()

	// 	for _, upper := range node.Start.All() {
	// 		for _, lower := range next.Start.All() {
	// 			if upper == lower {
	// 			}
	// 		}

	// 		for _, lower := range next.End.All() {
	// 			if upper == lower {
	// 			}
	// 		}
	// 	}

	// 	for _, upper := range node.End.All() {
	// 		for _, lower := range next.Start.All() {
	// 			if upper == lower {
	// 			}
	// 		}

	// 		for _, lower := range next.End.All() {
	// 			if upper == lower {
	// 			}
	// 		}
	// 	}
	// }
}

func (m *Merger) append(cursor *dstutil.Cursor) bool {
	node := cursor.Node()

	if node == nil {
		return true
	}

	// handle stride:define:literal annotation
	if _, declaration := m.literalType(AnnotationDefineLiteral, node); declaration != nil {
		m.Target.node.Decls = append(m.Target.node.Decls, declaration)
		return false
	}

	// handle stride:define:array annotation
	if _, declaration := m.arrayType(AnnotationDefineArray, node); declaration != nil {
		m.Target.node.Decls = append(m.Target.node.Decls, declaration)
		return false
	}

	// handle stride:define:struct annotation
	if _, declaration := m.structType(AnnotationDefineStruct, node); declaration != nil {
		m.Target.node.Decls = append(m.Target.node.Decls, declaration)
		return false
	}

	// handle stride:define:function annotation
	if _, declaration := m.functionType(AnnotationDefineFunction, node); declaration != nil {
		m.Target.node.Decls = append(m.Target.node.Decls, declaration)
		return false
	}

	return true
}

func (m *Merger) find(prefix, name string, target dst.Node) (tree dst.Node) {
	find := func(cursor *dstutil.Cursor) bool {
		if tree != nil {
			return false
		}

		if node := cursor.Node(); node != nil {
			if key := m.annotation(prefix, node.Decorations().Start); strings.EqualFold(key, name) {
				tree = node
			}
		}

		return tree == nil
	}

	dstutil.Apply(target, find, nil)
	return
}

func (m *Merger) arrayType(annotation string, node dst.Node) (string, *dst.GenDecl) {
	if declaration, ok := node.(*dst.GenDecl); ok {
		if declaration.Tok == token.TYPE {
			if specs := declaration.Specs; len(specs) == 1 {
				if typeSpec, ok := specs[0].(*dst.TypeSpec); ok {
					if _, ok := typeSpec.Type.(*dst.ArrayType); ok {
						name := inflect.Dasherize(typeSpec.Name.Name)
						if key := m.annotation(annotation, node.Decorations().Start); strings.EqualFold(name, key) {
							return name, declaration
						}
					}
				}
			}
		}
	}

	return "", nil
}

func (m *Merger) literalType(annotation string, node dst.Node) (string, *dst.GenDecl) {
	if declaration, ok := node.(*dst.GenDecl); ok {
		if declaration.Tok == token.TYPE {
			if specs := declaration.Specs; len(specs) == 1 {
				if typeSpec, ok := specs[0].(*dst.TypeSpec); ok {
					if _, ok := typeSpec.Type.(*dst.Ident); ok {
						name := inflect.Dasherize(typeSpec.Name.Name)
						if key := m.annotation(annotation, node.Decorations().Start); strings.EqualFold(name, key) {
							return name, declaration
						}
					}
				}
			}
		}
	}

	return "", nil
}

func (m *Merger) structType(annotation string, node dst.Node) (string, *dst.GenDecl) {
	if declaration, ok := node.(*dst.GenDecl); ok {
		if declaration.Tok == token.TYPE {
			if specs := declaration.Specs; len(specs) == 1 {
				if typeSpec, ok := specs[0].(*dst.TypeSpec); ok {
					if _, ok := typeSpec.Type.(*dst.StructType); ok {
						name := inflect.Dasherize(typeSpec.Name.Name)
						if key := m.annotation(annotation, node.Decorations().Start); strings.EqualFold(name, key) {
							return name, declaration
						}
					}
				}
			}
		}
	}

	return "", nil
}

func (m *Merger) structTypeProperties(node dst.Node) *dst.FieldList {
	if declaration, ok := node.(*dst.GenDecl); ok {
		if specs := declaration.Specs; len(specs) == 1 {
			if typeSpec, ok := specs[0].(*dst.TypeSpec); ok {
				if structType, ok := typeSpec.Type.(*dst.StructType); ok {
					return structType.Fields
				}
			}
		}
	}

	return &dst.FieldList{List: []*dst.Field{}}
}

func (m *Merger) functionType(annotation string, node dst.Node) (string, *dst.FuncDecl) {
	if declaration, ok := node.(*dst.FuncDecl); ok {
		name := inflect.Dasherize(declaration.Name.Name)

		if recv := declaration.Recv.List; len(recv) > 0 {
			name = fmt.Sprintf("%s:%s", inflect.Dasherize(kind(recv[0])), name)
		}

		if key := m.annotation(annotation, node.Decorations().Start); strings.EqualFold(name, key) {
			return name, declaration
		}
	}

	return "", nil
}

func (m *Merger) functionTypeBody(node dst.Node) *dst.BlockStmt {
	if declaration, ok := node.(*dst.FuncDecl); ok {
		return declaration.Body
	}
	return &dst.BlockStmt{List: []dst.Stmt{}}
}

func (m *Merger) functionTypeBodyRange(block *dst.BlockStmt) *Range {
	var (
		start *int
		end   *int
	)

	intPtr := func(value int) *int {
		return &value
	}

	annotated := func(key string) bool {
		return strings.EqualFold(key, "body")
	}

	for index, node := range block.List {
		decorations := node.Decorations()

		if start == nil {
			name := AnnotationDefineBlockStart

			if key := m.annotation(name, decorations.Start); annotated(key) {
				start = intPtr(index)
			} else if key := m.annotation(name, decorations.End); annotated(key) {
				start = intPtr(index + 1)
			}
		}

		if end == nil {
			name := AnnotationDefineBlockEnd

			if key := m.annotation(name, decorations.Start); annotated(key) {
				end = intPtr(index - 1)
			} else if key := m.annotation(name, decorations.End); annotated(key) {
				end = intPtr(index)
			}
		}
	}

	if start == nil || end == nil {
		return nil
	}

	return &Range{
		Start: *start,
		End:   *end,
	}
}

func (m *Merger) annotation(prefix string, decorations dst.Decorations) string {
	for _, comment := range decorations.All() {

		comment = strings.TrimPrefix(comment, "//")
		comment = strings.TrimSpace(comment)

		if strings.HasPrefix(comment, prefix) {
			comment = strings.TrimPrefix(comment, prefix)
			comment = strings.TrimSpace(comment)
			return comment
		}
	}

	return ""
}

func (m *Merger) decorations(node dst.Node) dst.Decorations {
	decorations := dst.Decorations{}

	add := func(decor dst.Decorations) {
		decorations.Append(decor.All()...)
	}

	switch node := node.(type) {
	case nil:
		// nothing to do
	case *dst.Field:
		add(node.Decs.Start)
		add(node.Decs.Type)
		add(node.Decs.End)
	case *dst.FieldList:
		add(node.Decs.Start)
		add(node.Decs.Opening)
		add(node.Decs.End)
	case *dst.BadExpr:
		// nothing to do
	case *dst.Ident:
		add(node.Decs.Start)
		add(node.Decs.X)
		add(node.Decs.End)
	case *dst.BasicLit:
		add(node.Decs.Start)
		add(node.Decs.End)
	case *dst.Ellipsis:
		add(node.Decs.Start)
		add(node.Decs.Ellipsis)
		add(node.Decs.End)
	case *dst.FuncLit:
		add(node.Decs.Start)
		add(node.Decs.Type)
		add(node.Decs.End)
	case *dst.CompositeLit:
		add(node.Decs.Start)
		add(node.Decs.Type)
		add(node.Decs.Lbrace)
		add(node.Decs.End)
	case *dst.ParenExpr:
		add(node.Decs.Start)
		add(node.Decs.Lparen)
		add(node.Decs.X)
		add(node.Decs.End)
	case *dst.SelectorExpr:
		add(node.Decs.Start)
		add(node.Decs.X)
		add(node.Decs.End)
	case *dst.IndexExpr:
		add(node.Decs.Start)
		add(node.Decs.X)
		add(node.Decs.Lbrack)
		add(node.Decs.Index)
		add(node.Decs.End)
	case *dst.SliceExpr:
		add(node.Decs.Start)
		add(node.Decs.X)
		add(node.Decs.Lbrack)
		add(node.Decs.Low)
		add(node.Decs.High)
		add(node.Decs.Max)
		add(node.Decs.End)
	case *dst.TypeAssertExpr:
		add(node.Decs.Start)
		add(node.Decs.X)
		add(node.Decs.Lparen)
		add(node.Decs.Type)
		add(node.Decs.End)
	case *dst.CallExpr:
		add(node.Decs.Start)
		add(node.Decs.Fun)
		add(node.Decs.Lparen)
		add(node.Decs.Ellipsis)
		add(node.Decs.End)
	case *dst.StarExpr:
		add(node.Decs.Start)
		add(node.Decs.Star)
		add(node.Decs.End)
	case *dst.UnaryExpr:
		add(node.Decs.Start)
		add(node.Decs.End)
	case *dst.BinaryExpr:
		add(node.Decs.Start)
		add(node.Decs.End)
	case *dst.KeyValueExpr:
		add(node.Decs.Start)
		add(node.Decs.End)
	case *dst.ArrayType:
		add(node.Decs.Start)
		add(node.Decs.End)
	case *dst.StructType:
		add(node.Decs.Start)
		add(node.Decs.End)
	case *dst.FuncType:
		add(node.Decs.Start)
		add(node.Decs.Func)
		add(node.Decs.Params)
		add(node.Decs.End)
	case *dst.InterfaceType:
		add(node.Decs.Start)
		add(node.Decs.Interface)
		add(node.Decs.End)
	case *dst.MapType:
		add(node.Decs.Start)
		add(node.Decs.Map)
		add(node.Decs.Key)
		add(node.Decs.End)
	case *dst.ChanType:
		add(node.Decs.Start)
		add(node.Decs.Begin)
		add(node.Decs.Arrow)
		add(node.Decs.End)
	case *dst.BadStmt:
		// nothing to do
	case *dst.DeclStmt:
		add(node.Decs.Start)
		add(node.Decs.End)
	case *dst.EmptyStmt:
		add(node.Decs.Start)
		add(node.Decs.End)
	case *dst.LabeledStmt:
		add(node.Decs.Start)
		add(node.Decs.Label)
		add(node.Decs.Colon)
		add(node.Decs.End)
	case *dst.ExprStmt:
		add(node.Decs.Start)
		add(node.Decs.End)
	case *dst.SendStmt:
		add(node.Decs.Start)
		add(node.Decs.Chan)
		add(node.Decs.Arrow)
		add(node.Decs.End)
	case *dst.IncDecStmt:
		add(node.Decs.Start)
		add(node.Decs.X)
		add(node.Decs.End)
	case *dst.AssignStmt:
		add(node.Decs.Start)
		add(node.Decs.Tok)
		add(node.Decs.End)
	case *dst.GoStmt:
		add(node.Decs.Start)
		add(node.Decs.Go)
		add(node.Decs.End)
	case *dst.DeferStmt:
		add(node.Decs.Start)
		add(node.Decs.Defer)
		add(node.Decs.End)
	case *dst.ReturnStmt:
		add(node.Decs.Start)
		add(node.Decs.Return)
		add(node.Decs.End)
	case *dst.BranchStmt:
		add(node.Decs.Start)
		add(node.Decs.Tok)
		add(node.Decs.End)
	case *dst.BlockStmt:
		add(node.Decs.Start)
		add(node.Decs.Lbrace)
		add(node.Decs.End)
	case *dst.IfStmt:
		add(node.Decs.Start)
		add(node.Decs.If)
		add(node.Decs.Init)
		add(node.Decs.Cond)
		add(node.Decs.Else)
		add(node.Decs.End)
	case *dst.CaseClause:
		add(node.Decs.Start)
		add(node.Decs.Case)
		add(node.Decs.Colon)
		add(node.Decs.End)
	case *dst.SwitchStmt:
		add(node.Decs.Start)
		add(node.Decs.Switch)
		add(node.Decs.Switch)
		add(node.Decs.Init)
		add(node.Decs.Tag)
		add(node.Decs.End)
	case *dst.TypeSwitchStmt:
		add(node.Decs.Start)
		add(node.Decs.Switch)
		add(node.Decs.Init)
		add(node.Decs.Assign)
		add(node.Decs.End)
	case *dst.CommClause:
		add(node.Decs.Start)
		add(node.Decs.Case)
		add(node.Decs.Comm)
		add(node.Decs.Colon)
		add(node.Decs.End)
	case *dst.SelectStmt:
		add(node.Decs.Start)
		add(node.Decs.Select)
		add(node.Decs.End)
	case *dst.ForStmt:
		add(node.Decs.Start)
		add(node.Decs.For)
		add(node.Decs.Init)
		add(node.Decs.Cond)
		add(node.Decs.Post)
		add(node.Decs.End)
	case *dst.RangeStmt:
		add(node.Decs.Start)
		add(node.Decs.For)
		add(node.Decs.Key)
		add(node.Decs.Value)
		add(node.Decs.Range)
		add(node.Decs.X)
		add(node.Decs.End)
	case *dst.ImportSpec:
		add(node.Decs.Start)
		add(node.Decs.Name)
		add(node.Decs.End)
	case *dst.ValueSpec:
		add(node.Decs.Start)
		add(node.Decs.Assign)
		add(node.Decs.End)
	case *dst.TypeSpec:
		add(node.Decs.Start)
		add(node.Decs.Name)
		add(node.Decs.End)
	case *dst.BadDecl:
		// nothing to do
	case *dst.GenDecl:
		add(node.Decs.Start)
		add(node.Decs.Tok)
		add(node.Decs.Lparen)
		add(node.Decs.End)
	case *dst.FuncDecl:
		add(node.Decs.Start)
		add(node.Decs.Func)
		add(node.Decs.Recv)
		add(node.Decs.Name)
		add(node.Decs.Params)
		add(node.Decs.Results)
		add(node.Decs.End)
	case *dst.File:
		add(node.Decs.Start)
		add(node.Decs.Package)
		add(node.Decs.Name)
		add(node.Decs.End)
	case *dst.Package:
		// nothing to do
	}

	return decorations
}

func body(node dst.Node) *dst.BlockStmt {
	if declaration, ok := node.(*dst.FuncDecl); ok {
		return declaration.Body
	}

	return nil
}
