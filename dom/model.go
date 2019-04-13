package dom

import "github.com/getkin/kin-openapi/openapi3"

// Operation represents the codegen operation
type Operation struct {
	Method      string
	Path        string
	Name        string
	Summary     string
	Description string
	Requests    []*Request
	Responses   []*Response
	Parameters  []*Parameter
}

// TypeDescriptor represents a type
type TypeDescriptor struct {
	Name        string
	Path        string
	Description string
	IsArray     bool
	IsClass     bool
	IsEnum      bool
	EnumValues  []interface{}
	Properties  []*PropertyDescriptor
	Parent      *TypeDescriptor
	Ref         *openapi3.SchemaRef
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

// Request for this operation
type Request struct {
	Description string
	Required    bool
	Type        string
}

// Response of the operation
type Response struct {
	Headers     []*Parameter
	Description string
	ContentType string
	Required    bool
	Type        string
}

// Parameter represents a codegen parameter
type Parameter struct {
	// Field *Field
	In string
}
