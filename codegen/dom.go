package codegen

import "github.com/jinzhu/copier"

// SpecDescriptor represents a spec
type SpecDescriptor struct {
	Schemas       TypeDescriptorCollection
	Parameters    ParameterDescriptorCollection
	Headers       HeaderDescriptorCollection
	RequestBodies RequestBodyDescriptorCollection
	Responses     ResponseDescriptorCollection
	Controllers   ControllerDescriptorCollection
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

// TypeDescriptor represents a type
type TypeDescriptor struct {
	Name        string
	Description string
	IsArray     bool
	IsClass     bool
	IsEnum      bool
	IsPrimitive bool
	Properties  PropertyDescriptorCollection
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
