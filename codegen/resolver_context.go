package codegen

import (
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/inflect"
)

var emptyCtx = &ResolverContext{}

// ResolverContext is the current resolver context
type ResolverContext struct {
	Name   string
	Schema *openapi3.SchemaRef
	Parent *ResolverContext
}

// IsRoot returns true if it's root
func (r *ResolverContext) IsRoot() bool {
	return r == emptyCtx
}

// Child returns the child context
func (r *ResolverContext) Child(name string, schema *openapi3.SchemaRef) *ResolverContext {
	ctx := &ResolverContext{
		Name:   r.name(name),
		Schema: schema,
		Parent: r,
	}

	return ctx
}

// Dereference returns the dereferenced context
func (r *ResolverContext) Dereference() *ResolverContext {
	ctx := &ResolverContext{
		Name:   filepath.Base(r.Schema.Ref),
		Schema: &openapi3.SchemaRef{Value: r.Schema.Value},
		Parent: r,
	}

	return ctx
}

// Array returns the array context
func (r *ResolverContext) Array() *ResolverContext {
	ctx := &ResolverContext{
		Name:   inflect.Singularize(r.Name),
		Schema: r.Schema.Value.Items,
		Parent: r,
	}

	return ctx
}

func (r *ResolverContext) name(text string) string {
	items := []string{}

	if r.Name != "" {
		items = append(items, r.Name)
	}

	if text != "" {
		items = append(items, text)
	}

	text = strings.Join(items, "_")
	text = inflect.Camelize(text)

	return text
}
