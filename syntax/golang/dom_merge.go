package golang

import (
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"
)

// Range represents the range
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
		left  = m.fieldList(target)
		right = m.fieldList(source)
	)

	for _, field := range right.List {
		if m.hasAnnotation(AnnotationDefine, field) {
			left.List = append(left.List, field)
		}
	}

	//TODO: sort fields by name
}

func (m *Merger) mergeFunc(target, source dst.Node) {
	var (
		// blocks
		left  = m.blockStmt(target)
		right = m.blockStmt(source)
		// ranges
		leftRange  = m.blockStmtRange(left)
		rightRange = m.blockStmtRange(right)
	)

	if leftRange == nil || rightRange == nil {
		return
	}

	var (
		result = []dst.Stmt{}
		items  = right.List[rightRange.Start : rightRange.End+1]
	)

	if len(items) == 0 {
		return
	}

	// append top block
	for index, item := range left.List {
		if index < leftRange.Start {
			result = append(result, item)
		}
	}

	for _, item := range items {
		result = append(result, item)
	}

	// append bottom block
	for index, item := range left.List {
		if index > leftRange.End {
			result = append(result, item)
		}
	}

	m.squash(result)

	left.List = result
}

func (m *Merger) blockStmtRange(block *dst.BlockStmt) *Range {
	var (
		start      *int
		end        *int
		annotation = AnnotationDefine
	)

	intPtr := func(v int) *int {
		return &v
	}

	for index, node := range block.List {
		decorations := node.Decorations()

		if start == nil {
			if annotation.In(decorations.Start, bodyStartKey) {
				start = intPtr(index)
			} else if annotation.In(decorations.End, bodyStartKey) {
				start = intPtr(index + 1)
			}
		}

		if end == nil {
			if annotation.In(decorations.Start, bodyEndKey) {
				end = intPtr(index - 1)
			} else if annotation.In(decorations.End, bodyEndKey) {
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

func (m *Merger) squash(items []dst.Stmt) {
	kv := map[string]bool{}

	remove := func(kind string, node *dst.NodeDecs) {
		var (
			comments    = []string{}
			decorations *dst.Decorations
		)

		switch kind {
		case "start":
			decorations = &node.Start
		case "end":
			decorations = &node.End
		}

		for _, comment := range decorations.All() {
			comment = strings.TrimSpace(comment)

			if len(comment) == 0 {
				comments = append(comments, newline)
				continue
			}

			if strings.EqualFold(comment, bodyInfo) {
				continue
			}

			if _, ok := kv[comment]; ok {
				continue
			}

			switch {
			case strings.EqualFold(comment, bodyStart):
				switch kind {
				case "start":
					node.Before = dst.EmptyLine
				case "end":
					//TODO: handle this case
				}
			case strings.EqualFold(comment, bodyEnd):
				switch kind {
				case "start":
					node.Before = dst.NewLine
				case "end":
					//TODO: handle this case
				}
			}

			comments = append(comments, comment)
			// mark as processed
			kv[comment] = true
		}

		decorations.Replace(comments...)
	}

	for _, item := range items {
		decorations := item.Decorations()

		remove("start", decorations)
		remove("end", decorations)
	}
}

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

func (m *Merger) findAnnotation(annotation Annotation, node dst.Node) (string, bool) {
	return annotation.Find(node.Decorations().Start)
}

func (m *Merger) hasAnnotation(annotation Annotation, node dst.Node) bool {
	_, ok := annotation.Find(node.Decorations().Start)
	return ok
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

func (m *Merger) fieldList(node dst.Node) *dst.FieldList {
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

func (m *Merger) blockStmt(node dst.Node) *dst.BlockStmt {
	if declaration, ok := node.(*dst.FuncDecl); ok {
		return declaration.Body
	}

	return &dst.BlockStmt{List: []dst.Stmt{}}
}
