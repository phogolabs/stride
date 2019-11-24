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
	// AnnotationNote represents the annotation for note
	AnnotationNote Annotation = "NOTE:"
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

// Has returns true if the annotation with given name exists
func (n Annotation) Has(decorations dst.Decorations, name string) bool {
	prefix := string(n)
	for _, comment := range decorations.All() {
		comment = strings.TrimPrefix(comment, "//")
		comment = strings.TrimSpace(comment)

		if strings.HasPrefix(comment, prefix) {
			comment = strings.TrimPrefix(comment, prefix)
			comment = strings.TrimSpace(comment)

			if strings.EqualFold(name, comment) {
				return true
			}
		}
	}

	return false
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
	const (
		keyStart = "body:start"
		keyEnd   = "body:end"
	)

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
			if annotation.Has(decorations.Start, keyStart) {
				start = intPtr(index)
			} else if annotation.Has(decorations.End, keyStart) {
				start = intPtr(index + 1)
			}
		}

		if end == nil {
			if annotation.Has(decorations.Start, keyEnd) {
				end = intPtr(index - 1)
			} else if annotation.Has(decorations.End, keyEnd) {
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
	var (
		kv    = map[string]bool{}
		help  = AnnotationNote.Format("write your code here")
		start = AnnotationDefine.Format("body:start")
		end   = AnnotationDefine.Format("body:end")
	)

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

			if strings.EqualFold(comment, help) {
				continue
			}

			if _, ok := kv[comment]; ok {
				continue
			}

			if strings.EqualFold(comment, start) {
				switch kind {
				case "start":
					node.Before = dst.EmptyLine
					node.After = dst.NewLine
				case "end":
					node.Before = dst.NewLine
					node.After = dst.NewLine
				}
			}

			if strings.EqualFold(comment, end) {
				switch kind {
				case "start":
					node.Before = dst.NewLine
					node.After = dst.NewLine
				case "end":
					node.Before = dst.NewLine
					node.After = dst.EmptyLine
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
