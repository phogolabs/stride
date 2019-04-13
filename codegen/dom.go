package codegen

import "github.com/getkin/kin-openapi/openapi3"

// TypeDescriptor represents a type
type TypeDescriptor struct {
	Name        string
	Path        string
	Description string
	IsArray     bool
	IsClass     bool
	IsEnum      bool
	IsPrimitive bool
	Properties  []*PropertyDescriptor
	Parent      *TypeDescriptor
	Ref         *openapi3.SchemaRef
}

// IsNested returns true if the type is nested
func (t *TypeDescriptor) IsNested() bool {
	return !t.IsPrimitive && t.Ref != nil && t.Ref.Ref == ""
}

// PropertyDescriptor definition
type PropertyDescriptor struct {
	Name          string
	Description   string
	Nullable      bool
	Required      bool
	PropertyType  *TypeDescriptor
	ComponentType *TypeDescriptor
	Ref           *openapi3.SchemaRef
}

// OperationDescriptor represents the codegen operation
type OperationDescriptor struct {
	Method      string
	Path        string
	Name        string
	Summary     string
	Description string
}
