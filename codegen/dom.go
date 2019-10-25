package codegen

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/jinzhu/copier"
)

// SpecDescriptor represents a spec
type SpecDescriptor struct {
	Schemas       TypeDescriptorCollection
	Parameters    ParameterDescriptorCollection
	Headers       HeaderDescriptorCollection
	RequestBodies RequestBodyDescriptorCollection
	Responses     ResponseDescriptorCollection
	Operations    OperationDescriptorCollection
}

// TypeDescriptorCollection definition
type TypeDescriptorCollection []*TypeDescriptor

// Len is the number of elements in the collection.
func (t TypeDescriptorCollection) Len() int {
	return len(t)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (t TypeDescriptorCollection) Less(i, j int) bool {
	return t[i].Name < t[j].Name
}

// Swap swaps the elements with indexes i and j.
func (t TypeDescriptorCollection) Swap(i, j int) {
	var x = t[i]
	t[i] = t[j]
	t[j] = x
}

// TypeDeclaration represents a type declaration
type TypeDeclaration struct {
	Name      string
	SchemaRef *openapi3.SchemaRef
}

// NewTypeDeclarationPrimitive creates a type declaration primitive
func NewTypeDeclarationPrimitive(schemaRef *openapi3.SchemaRef) *TypeDeclaration {
	switch schemaRef.Value.Type {
	case "integer":
		switch schemaRef.Value.Format {
		case "int64":
			return &TypeDeclaration{
				Name:      "int64",
				SchemaRef: schemaRef,
			}
		default:
			return &TypeDeclaration{
				Name:      "int32",
				SchemaRef: schemaRef,
			}
		}
	case "number":
		switch schemaRef.Value.Format {
		case "double":
			return &TypeDeclaration{
				Name:      "double",
				SchemaRef: schemaRef,
			}
		default:
			return &TypeDeclaration{
				Name:      "float",
				SchemaRef: schemaRef,
			}
		}
	case "string":
		switch schemaRef.Value.Format {
		case "binary":
			return &TypeDeclaration{
				Name:      "binary",
				SchemaRef: schemaRef,
			}
		case "byte":
			return &TypeDeclaration{
				Name:      "byte",
				SchemaRef: schemaRef,
			}
		case "date":
			return &TypeDeclaration{
				Name:      "date",
				SchemaRef: schemaRef,
			}
		case "date-time":
			return &TypeDeclaration{
				Name:      "date-time",
				SchemaRef: schemaRef,
			}
		case "uuid":
			return &TypeDeclaration{
				Name:      "uuid",
				SchemaRef: schemaRef,
			}
		default:
			return &TypeDeclaration{
				Name:      "string",
				SchemaRef: schemaRef,
			}
		}
	case "boolean":
		return &TypeDeclaration{
			Name:      "boolean",
			SchemaRef: schemaRef,
		}
	default:
		return nil
	}
}

// TypeDescriptor represents a type
type TypeDescriptor struct {
	Name        string
	Description string
	IsArray     bool
	IsClass     bool
	IsEnum      bool
	IsPrimitive bool
	Metadata    map[string]interface{}
	Properties  PropertyDescriptorCollection
}

// NewTypeDescriptor creates a new type descriptor
func NewTypeDescriptor(declaration *TypeDeclaration) *TypeDescriptor {
	descriptor := &TypeDescriptor{
		Name:        declaration.Name,
		Description: declaration.SchemaRef.Value.Description,
		Metadata:    make(map[string]interface{}),
	}
	return descriptor
}

// NewTypeDescriptorPrimitive creates a new primitive type descriptor
func NewTypeDescriptorPrimitive(declaration *TypeDeclaration) *TypeDescriptor {
	descriptor := NewTypeDescriptor(declaration)
	descriptor.IsPrimitive = true
	return descriptor
}

// NewTypeDescriptorEnum creates a new type descriptor for enum
func NewTypeDescriptorEnum(declaration *TypeDeclaration) *TypeDescriptor {
	descriptor := NewTypeDescriptor(declaration)
	descriptor.IsEnum = true
	descriptor.Metadata["enum"] = declaration.SchemaRef.Value.Enum
	return descriptor
}

// Imports returns the required imports
func (t *TypeDescriptor) Imports() []string {
	var imports []string

	for _, property := range t.Properties {
		namespaces := property.PropertyType.Imports()
		imports = append(imports, namespaces...)
	}

	return imports
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

// PropertyDescriptorCollection definition
type PropertyDescriptorCollection []*PropertyDescriptor

// Len is the number of elements in the collection.
func (t PropertyDescriptorCollection) Len() int {
	return len(t)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (t PropertyDescriptorCollection) Less(i, j int) bool {
	return t[i].Name < t[j].Name
}

// Swap swaps the elements with indexes i and j.
func (t PropertyDescriptorCollection) Swap(i, j int) {
	var x = t[i]
	t[i] = t[j]
	t[j] = x
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

// ParameterDescriptorCollection definition
type ParameterDescriptorCollection []*ParameterDescriptor

// Len is the number of elements in the collection.
func (t ParameterDescriptorCollection) Len() int {
	return len(t)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (t ParameterDescriptorCollection) Less(i, j int) bool {
	return t[i].Name < t[j].Name
}

// Swap swaps the elements with indexes i and j.
func (t ParameterDescriptorCollection) Swap(i, j int) {
	var x = t[i]
	t[i] = t[j]
	t[j] = x
}

// HeaderDescriptor definition
type HeaderDescriptor struct {
	Name       string
	HeaderType *TypeDescriptor
}

// HeaderDescriptorCollection definition
type HeaderDescriptorCollection []*HeaderDescriptor

// Len is the number of elements in the collection.
func (t HeaderDescriptorCollection) Len() int {
	return len(t)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (t HeaderDescriptorCollection) Less(i, j int) bool {
	return t[i].Name < t[j].Name
}

// Swap swaps the elements with indexes i and j.
func (t HeaderDescriptorCollection) Swap(i, j int) {
	var x = t[i]
	t[i] = t[j]
	t[j] = x
}

// RequestBodyDescriptor definition
type RequestBodyDescriptor struct {
	Name        string
	Description string
	Required    bool
	Contents    ContentDescriptorCollection
}

// RequestBodyDescriptorCollection definition
type RequestBodyDescriptorCollection []*RequestBodyDescriptor

// Len is the number of elements in the collection.
func (t RequestBodyDescriptorCollection) Len() int {
	return len(t)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (t RequestBodyDescriptorCollection) Less(i, j int) bool {
	return t[i].Name < t[j].Name
}

// Swap swaps the elements with indexes i and j.
func (t RequestBodyDescriptorCollection) Swap(i, j int) {
	var x = t[i]
	t[i] = t[j]
	t[j] = x
}

// ResponseDescriptor definition
type ResponseDescriptor struct {
	Code        int
	Name        string
	Description string
	Headers     HeaderDescriptorCollection
	Contents    ContentDescriptorCollection
}

// ResponseDescriptorCollection definition
type ResponseDescriptorCollection []*ResponseDescriptor

// Len is the number of elements in the collection.
func (t ResponseDescriptorCollection) Len() int {
	return len(t)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (t ResponseDescriptorCollection) Less(i, j int) bool {
	return t[i].Name < t[j].Name
}

// Swap swaps the elements with indexes i and j.
func (t ResponseDescriptorCollection) Swap(i, j int) {
	var x = t[i]
	t[i] = t[j]
	t[j] = x
}

// ContentDescriptor definition
type ContentDescriptor struct {
	Name        string
	ContentType *TypeDescriptor
}

// ContentDescriptorCollection definition
type ContentDescriptorCollection []*ContentDescriptor

// Len is the number of elements in the collection.
func (t ContentDescriptorCollection) Len() int {
	return len(t)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (t ContentDescriptorCollection) Less(i, j int) bool {
	return t[i].Name < t[j].Name
}

// Swap swaps the elements with indexes i and j.
func (t ContentDescriptorCollection) Swap(i, j int) {
	var x = t[i]
	t[i] = t[j]
	t[j] = x
}

// ControllerDescriptor definition
type ControllerDescriptor struct {
	Name        string
	Description string
	Operations  OperationDescriptorCollection
}

// ControllerDescriptorCollection definition
type ControllerDescriptorCollection []*ControllerDescriptor

// Len is the number of elements in the collection.
func (t ControllerDescriptorCollection) Len() int {
	return len(t)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (t ControllerDescriptorCollection) Less(i, j int) bool {
	return t[i].Name < t[j].Name
}

// Swap swaps the elements with indexes i and j.
func (t ControllerDescriptorCollection) Swap(i, j int) {
	var x = t[i]
	t[i] = t[j]
	t[j] = x
}

// OperationDescriptor definition
type OperationDescriptor struct {
	Method      string
	Path        string
	Name        string
	Summary     string
	Description string
	Deprecated  bool
	Tags        []string
	Parameters  ParameterDescriptorCollection
	Responses   ResponseDescriptorCollection
	RequestBody *RequestBodyDescriptor
}

// OperationDescriptorCollection definition
type OperationDescriptorCollection []*OperationDescriptor

// Len is the number of elements in the collection.
func (t OperationDescriptorCollection) Len() int {
	return len(t)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (t OperationDescriptorCollection) Less(i, j int) bool {
	return t[i].Name < t[j].Name
}

// Swap swaps the elements with indexes i and j.
func (t OperationDescriptorCollection) Swap(i, j int) {
	var x = t[i]
	t[i] = t[j]
	t[j] = x
}
