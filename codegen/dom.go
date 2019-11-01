package codegen

import (
	"fmt"
	"io"
	"sort"
)

// Metadata of the TypeDescriptor
type Metadata map[string]interface{}

// SpecDescriptor represents a spec
type SpecDescriptor struct {
	Types      TypeDescriptorCollection
	Operations OperationDescriptorCollection
}

// TypeDescriptorMap definition
type TypeDescriptorMap map[string]*TypeDescriptor

// CollectFrom the collection
func (m TypeDescriptorMap) CollectFrom(descriptors TypeDescriptorCollection) {
	for _, descriptor := range descriptors {
		m.add(descriptor)
	}
}

// Collection return the map as collection
func (m TypeDescriptorMap) Collection() TypeDescriptorCollection {
	descriptors := TypeDescriptorCollection{}

	for _, descriptor := range m {
		descriptors = append(descriptors, descriptor)
	}

	// sort the descriptors
	sort.Sort(descriptors)

	return descriptors
}

func (m TypeDescriptorMap) add(descriptor *TypeDescriptor) {
	if descriptor.IsPrimitive {
		return
	}

	key := descriptor.Name

	if _, ok := m[key]; ok {
		return
	}

	m[key] = descriptor

	if element := descriptor.Element; element != nil {
		m.add(element)
	}

	for _, property := range descriptor.Properties {
		m.add(property.PropertyType)
	}
}

// TypeDescriptorCollection definition
type TypeDescriptorCollection []*TypeDescriptor

// Render renders the collection
func (t TypeDescriptorCollection) Render(w io.Writer) {
	end := len(t) - 1

	for index, descriptor := range t {
		if index > 0 {
			fmt.Fprintln(w)
		}

		descriptor.Render(w)

		if index < end {
			fmt.Fprintln(w)
		}
	}
}

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
	IsAlias     bool
	Element     *TypeDescriptor
	Metadata    Metadata
	Properties  PropertyDescriptorCollection
}

// String returns the string representation
func (t *TypeDescriptor) String() string {
	if t.IsClass {
		return fmt.Sprintf("*%s", t.Name)
	}

	return t.Name
}

// Render renders the type
func (t *TypeDescriptor) Render(w io.Writer) {
	fmt.Fprintf(w, "// %s is auto-generated type", t.Name)
	fmt.Fprintln(w)

	if t.Description != "" {
		fmt.Fprintf(w, "// %s ", t.Description)
		fmt.Fprintln(w)
	}

	switch {
	case t.IsAlias:
		fmt.Fprintf(w, "type %s %s", t.Name, t.Element.Name)
	case t.IsArray:
		fmt.Fprintf(w, "type %s []%v", t.Name, t.Element)
	case t.IsClass:
		fmt.Fprintf(w, "type %s struct {", t.Name)
		t.Properties.Render(w)
		fmt.Fprintf(w, "}")
	case t.IsEnum:
		fmt.Fprintf(w, "type %s %v", t.Name, t.Element.Name)
	}
}

// PropertyDescriptor definition
type PropertyDescriptor struct {
	Name         string
	Description  string
	Nullable     bool
	Required     bool
	PropertyType *TypeDescriptor
}

// Render renders the type
func (p *PropertyDescriptor) Render(w io.Writer) {
	if p.Description != "" {
		fmt.Fprintf(w, "// %s ", p.Description)
		fmt.Fprintln(w)
	}

	fmt.Fprintf(w, "%s %v %s", camelize(p.Name), p.PropertyType, p.Tag())
}

// Tag returns the tag
func (p *PropertyDescriptor) Tag() string {
	return ""
}

// PropertyDescriptorCollection definition
type PropertyDescriptorCollection []*PropertyDescriptor

// Render renders the collection
func (t PropertyDescriptorCollection) Render(w io.Writer) {
	end := len(t) - 1

	for index, property := range t {
		fmt.Fprintln(w)

		property.Render(w)

		if index == end {
			fmt.Fprintln(w)
		}
	}
}

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

// RequestDescriptor represents a request descriptor
type RequestDescriptor struct {
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
