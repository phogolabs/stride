package codegen

import "github.com/jinzhu/copier"

// Spec represents a spec
type Spec struct {
	Schemas       []*TypeDescriptor
	Parameters    []*ParameterDescriptor
	Headers       []*HeaderDescriptor
	RequestBodies []*RequestBodyDescriptor
	Responses     []*ResponseDescriptor
	Controllers   []*ControllerDescriptor
}

// TypeDescriptor represents a type
type TypeDescriptor struct {
	Name        string
	Description string
	IsArray     bool
	IsClass     bool
	IsEnum      bool
	IsPrimitive bool
	Properties  []*PropertyDescriptor
}

// Clone clones the object
func (t *TypeDescriptor) Clone() *TypeDescriptor {
	descriptor := &TypeDescriptor{}
	copier.Copy(descriptor, t)
	return descriptor
}

// PropertyDescriptor definition
type PropertyDescriptor struct {
	Name         string
	Description  string
	Nullable     bool
	Required     bool
	PropertyType *TypeDescriptor
}

// ParameterDescriptor definition
type ParameterDescriptor struct {
	Name          string
	Path          string
	In            string
	Description   string
	Required      bool
	Deprecated    bool
	ParameterType *TypeDescriptor
}

// HeaderDescriptor definition
type HeaderDescriptor struct {
	Name       string
	HeaderType *TypeDescriptor
}

// RequestBodyDescriptor definition
type RequestBodyDescriptor struct {
	Name        string
	Description string
	Required    bool
	Contents    []*ContentDescriptor
}

// ResponseDescriptor definition
type ResponseDescriptor struct {
	Description string
	Code        int
	Headers     []*HeaderDescriptor
	Contents    []*ContentDescriptor
}

// ContentDescriptor definition
type ContentDescriptor struct {
	Name        string
	ContentType *TypeDescriptor
}

// ControllerDescriptor definition
type ControllerDescriptor struct {
	Name        string
	Description string
	Operations  []*OperationDescriptor
}

// OperationDescriptor definition
type OperationDescriptor struct {
	Method        string
	Path          string
	Name          string
	Summary       string
	Description   string
	Deprecated    bool
	Tags          []string
	Parameters    []*ParameterDescriptor
	Responses     []*ResponseDescriptor
	RequestBodies []*RequestBodyDescriptor
}
