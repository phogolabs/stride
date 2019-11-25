package golang

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/dave/dst"
	"github.com/phogolabs/stride/inflect"
)

const (
	// AnnotationGenerate represents the annotation for tgenerated code
	AnnotationGenerate Annotation = "stride:generate"
	// AnnotationDefine represents the annotation for user-defined code
	AnnotationDefine Annotation = "stride:define"
	// AnnotationNotice represents the annotation for note
	AnnotationNotice Annotation = "NOTE:"
)

var (
	bodyInfo  = AnnotationNotice.Comment("not implemented")
	bodyStart = AnnotationDefine.Key("body", "start")
	bodyEnd   = AnnotationDefine.Key("body", "end")
)

const (
	bodyStartKey = "body:start"
	bodyEndKey   = "body:end"
	newline      = "\n"
)

// Annotation represents an annotation
type Annotation string

// Key formats the annotation as key
func (n Annotation) Key(text ...string) string {
	buffer := &bytes.Buffer{}

	for _, part := range text {
		if part = strings.TrimSpace(part); part == "" {
			continue
		}

		part = inflect.Dasherize(part)

		if buffer.Len() > 0 {
			fmt.Fprint(buffer, ":")
		}

		fmt.Fprint(buffer, part)
	}

	return n.Comment(buffer.String())
}

// Comment formats the annotation as comment
func (n Annotation) Comment(text string) string {
	return fmt.Sprintf("// %s %s", n, text)
}

// Find returns the name of the annotation of exists in the decorations
func (n Annotation) Find(decorations dst.Decorations) (string, bool) {
	var (
		prefix = string(n)
		name   string
	)

	for _, comment := range decorations.All() {
		comment = n.uncomment(comment)

		if strings.HasPrefix(comment, prefix) {
			name = strings.TrimPrefix(comment, prefix)
			name = strings.TrimSpace(name)

			return name, true
		}
	}

	return name, false
}

// In returns true if the annotation with given name exists in the decorations list
func (n Annotation) In(decorations dst.Decorations, term string) bool {
	var (
		prefix = string(n)
		name   string
	)

	for _, comment := range decorations.All() {
		comment = n.uncomment(comment)

		if strings.HasPrefix(comment, prefix) {
			name = strings.TrimPrefix(comment, prefix)
			name = strings.TrimSpace(name)

			if strings.EqualFold(name, term) {
				return true
			}
		}
	}

	return false
}

func (n Annotation) uncomment(comment string) string {
	comment = strings.TrimPrefix(comment, "//")
	comment = strings.TrimSpace(comment)
	return comment
}
