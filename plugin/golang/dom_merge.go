package golang

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"
)

const (
	// AnnotationGenerate represents the annotation for tgenerated code
	AnnotationGenerate Annotation = "stride:generate"
	// AnnotationDefine represents the annotation for user-defined code
	AnnotationDefine Annotation = "stride:define"
	// AnnotationDefineBlockStart represents the annotation for defined block start
	AnnotationDefineBlockStart Annotation = "stride:define:block:start"
	// AnnotationDefineBlockEnd represents the annotation for defined block end
	AnnotationDefineBlockEnd Annotation = "stride:define:block:end"
)

// Annotation represents an annotation
type Annotation string

// Format formats the annotation
func (n Annotation) Format(text ...string) string {
	buffer := &bytes.Buffer{}

	for _, part := range text {
		if part = strings.TrimSpace(part); part == "" {
			continue
		}

		if buffer.Len() > 0 {
			fmt.Fprint(buffer, ":")
		}

		fmt.Fprint(buffer, part)
	}

	return fmt.Sprintf("// %s %s", n, buffer.String())
}

// Find returns the name of the annotation of exists in the decorations
func (n Annotation) Find(decorations dst.Decorations) (string, bool) {
	prefix := string(n)

	for _, comment := range decorations.All() {

		comment = strings.TrimPrefix(comment, "//")
		comment = strings.TrimSpace(comment)

		if strings.HasPrefix(comment, prefix) {
			comment = strings.TrimPrefix(comment, prefix)
			comment = strings.TrimSpace(comment)
			return comment, true
		}
	}

	return "", false
}

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
	var (
		node       = cursor.Node()
		parent     = cursor.Parent()
		annotation = AnnotationGenerate
	)

	if node == nil {
		return true
	}

	if _, ok := node.(*dst.File); ok {
		return true
	}

	if _, ok := parent.(*dst.File); !ok {
		return false
	}

	if name, ok := m.findAnnotation(annotation, node); ok {
		if source := m.findNode(annotation, name, m.Source.node); source != nil {
			// merge the nodes if they are a struct
			m.mergeStruct(node, source)
			// merge the node if they are a func
			m.mergeFunc(node, source)
		}
	}

	return false
}

func (m *Merger) mergeStruct(target, source dst.Node) {
	var (
		left  = m.structTypeProperties(target)
		right = m.structTypeProperties(source)
	)

	for _, field := range right.List {
		if m.hasAnnotation(AnnotationDefine, field) {
			left.List = append(left.List, field)
		}
	}

	//TODO: sort fields by name
}

func (m *Merger) mergeFunc(target, source dst.Node) {
}

// func (m *Merger) squash(target, source dst.Node) {
// 	var (
// 		left       = m.functionTypeBody(target)
// 		leftRange  = m.functionTypeBodyRange(left)
// 		right      = m.functionTypeBody(source)
// 		rightRange = m.functionTypeBodyRange(right)
// 	)

// 	if leftRange != nil && rightRange != nil {
// 		var (
// 			result = []dst.Stmt{}
// 			items  = right.List[rightRange.Start : rightRange.End+1]
// 		)

// 		// append top block
// 		for index, item := range left.List {
// 			if index < leftRange.Start {
// 				result = append(result, item)
// 			}
// 		}

// 		// append the range block
// 		result = append(result, items...)

// 		// append bottom block
// 		for index, item := range left.List {
// 			if index > leftRange.End {
// 				result = append(result, item)
// 			}
// 		}

// 		// sanitize comments
// 		m.sanitize(result)

// 		left.List = result
// 	}
// }

// func (m *Merger) sanitize(items []dst.Stmt) {
// 	// const name = "body"

// 	// for index := 0; index < len(items)-1; index++ {
// 	// 	node := items[index].Decorations()
// 	// 	next := items[index+1].Decorations()

// 	// 	for _, upper := range node.Start.All() {
// 	// 		for _, lower := range next.Start.All() {
// 	// 			if upper == lower {
// 	// 			}
// 	// 		}

// 	// 		for _, lower := range next.End.All() {
// 	// 			if upper == lower {
// 	// 			}
// 	// 		}
// 	// 	}

// 	// 	for _, upper := range node.End.All() {
// 	// 		for _, lower := range next.Start.All() {
// 	// 			if upper == lower {
// 	// 			}
// 	// 		}

// 	// 		for _, lower := range next.End.All() {
// 	// 			if upper == lower {
// 	// 			}
// 	// 		}
// 	// 	}
// 	// }
// }

func (m *Merger) append(cursor *dstutil.Cursor) bool {
	var (
		node       = cursor.Node()
		parent     = cursor.Parent()
		annotation = AnnotationDefine
	)

	if node == nil {
		return true
	}

	if _, ok := node.(*dst.File); ok {
		return true
	}

	if _, ok := parent.(*dst.File); !ok {
		return false
	}

	// handle stride:define annotation
	if m.hasAnnotation(annotation, node) {
		if declaration, ok := node.(dst.Decl); ok {
			m.Target.node.Decls = append(m.Target.node.Decls, declaration)
		}
	}

	return false
}

// func (m *Merger) findNodeByName(annotation Annotation, name string, target dst.Node) (tree dst.Node) {
// find := func(cursor *dstutil.Cursor) bool {
// 	if tree != nil {
// 		return false
// 	}

// 	if node := cursor.Node(); node != nil {
// 		if key := m.annotation(prefix, node.Decorations().Start); strings.EqualFold(key, name) {
// 			tree = node
// 		}

// 	if name, ok := annotation.Find(node.Decorations().Start); ok {
// 			tree = node
// 	}
// 	}

// 	return tree == nil
// }

// dstutil.Apply(target, find, nil)
// return
// }

func (m *Merger) hasAnnotation(annotation Annotation, node dst.Node) bool {
	_, ok := annotation.Find(node.Decorations().Start)
	return ok
}

func (m *Merger) findAnnotation(annotation Annotation, node dst.Node) (string, bool) {
	return annotation.Find(node.Decorations().Start)
}

func (m *Merger) findNode(annotation Annotation, key string, node dst.Node) (tree dst.Node) {
	find := func(cursor *dstutil.Cursor) bool {
		var (
			node       = cursor.Node()
			annotation = AnnotationGenerate
		)

		if node == nil {
			return true
		}

		if name, ok := m.findAnnotation(annotation, node); ok {
			if strings.EqualFold(name, key) {
				tree = node
			}
		}

		return tree == nil
	}

	dstutil.Apply(node, find, nil)
	return
}

// func (m *Merger) find(annotation Annotation, node dst.Node) (string, dst.Node) {
// 	if name, ok := annotation.Find(node.Decorations().Start); ok {
// 		return name, declaration
// 	}

// 	return "", nil
// }

// func (m *Merger) findByName(annotation Annotation, name string, node dst.Node) dst.Node {
// 	if _, ok := annotation.Find(node.Decorations().Start); ok {
// 		return node
// 	}

// 	return nil
// }

// func (m *Merger) arrayType(annotation string, node dst.Node) (string, *dst.GenDecl) {
// 	if declaration, ok := node.(*dst.GenDecl); ok {
// 		if declaration.Tok == token.TYPE {
// 			if specs := declaration.Specs; len(specs) == 1 {
// 				if typeSpec, ok := specs[0].(*dst.TypeSpec); ok {
// 					if _, ok := typeSpec.Type.(*dst.ArrayType); ok {
// 						name := inflect.Dasherize(typeSpec.Name.Name)
// 						if key := m.annotation(annotation, node.Decorations().Start); strings.EqualFold(name, key) {
// 							return name, declaration
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}

// 	return "", nil
// }

// func (m *Merger) literalType(annotation string, node dst.Node) (string, *dst.GenDecl) {
// 	if declaration, ok := node.(*dst.GenDecl); ok {
// 		if declaration.Tok == token.TYPE {
// 			if specs := declaration.Specs; len(specs) == 1 {
// 				if typeSpec, ok := specs[0].(*dst.TypeSpec); ok {
// 					if _, ok := typeSpec.Type.(*dst.Ident); ok {
// 						name := inflect.Dasherize(typeSpec.Name.Name)
// 						if key := m.annotation(annotation, node.Decorations().Start); strings.EqualFold(name, key) {
// 							return name, declaration
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}

// 	return "", nil
// }

// func (m *Merger) structType(annotation string, node dst.Node) (string, *dst.GenDecl) {
// 	if declaration, ok := node.(*dst.GenDecl); ok {
// 		if declaration.Tok == token.TYPE {
// 			if specs := declaration.Specs; len(specs) == 1 {
// 				if typeSpec, ok := specs[0].(*dst.TypeSpec); ok {
// 					if _, ok := typeSpec.Type.(*dst.StructType); ok {
// 						name := inflect.Dasherize(typeSpec.Name.Name)
// 						if key := m.annotation(annotation, node.Decorations().Start); strings.EqualFold(name, key) {
// 							return name, declaration
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}

// 	return "", nil
// }

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

// func (m *Merger) functionType(annotation Annotation, node dst.Node) (string, *dst.FuncDecl) {
// 	if declaration, ok := node.(*dst.FuncDecl); ok {
// 		name := inflect.Dasherize(declaration.Name.Name)

// 		if recv := declaration.Recv.List; len(recv) > 0 {
// 			name = fmt.Sprintf("%s:%s", inflect.Dasherize(kind(recv[0])), name)
// 		}

// 		if key := m.annotation(annotation, node.Decorations().Start); strings.EqualFold(name, key) {
// 			return name, declaration
// 		}
// 	}

// 	return "", nil
// }

// func (m *Merger) functionTypeBody(node dst.Node) *dst.BlockStmt {
// 	if declaration, ok := node.(*dst.FuncDecl); ok {
// 		return declaration.Body
// 	}
// 	return &dst.BlockStmt{List: []dst.Stmt{}}
// }

// func (m *Merger) functionTypeBodyRange(block *dst.BlockStmt) *Range {
// 	var (
// 		start *int
// 		end   *int
// 	)

// 	intPtr := func(value int) *int {
// 		return &value
// 	}

// 	annotated := func(key string) bool {
// 		return strings.EqualFold(key, "body")
// 	}

// 	for index, node := range block.List {
// 		decorations := node.Decorations()

// 		if start == nil {
// 			name := AnnotationDefineBlockStart

// 			if key := m.annotation(name, decorations.Start); annotated(key) {
// 				start = intPtr(index)
// 			} else if key := m.annotation(name, decorations.End); annotated(key) {
// 				start = intPtr(index + 1)
// 			}
// 		}

// 		if end == nil {
// 			name := AnnotationDefineBlockEnd

// 			if key := m.annotation(name, decorations.Start); annotated(key) {
// 				end = intPtr(index - 1)
// 			} else if key := m.annotation(name, decorations.End); annotated(key) {
// 				end = intPtr(index)
// 			}
// 		}
// 	}

// 	if start == nil || end == nil {
// 		return nil
// 	}

// 	return &Range{
// 		Start: *start,
// 		End:   *end,
// 	}
// }

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
