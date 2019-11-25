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
	// AnnotationNote represents the annotation for note
	AnnotationNote Annotation = "NOTE:"
)

const (
	bodyStart   = "body:start"
	bodyMessage = "Not Implemented"
	bodyEnd     = "body:end"
)

// Annotation represents an annotation
type Annotation string

// Format formats the annotation
func (n Annotation) Format(text ...string) string {
	buffer := &bytes.Buffer{}

	for _, part := range text {
		part = inflect.Dasherize(part)

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
