package codegen

import (
	"sort"
	"strings"
)

// Metadata of the TypeDescriptor
type Metadata map[string]interface{}

// SpecDescriptor represents a spec
type SpecDescriptor struct {
	Types       TypeDescriptorCollection
	Controllers ControllerDescriptorCollection
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
	IsNullable  bool
	Element     *TypeDescriptor
	Default     interface{}
	Metadata    Metadata
	Properties  PropertyDescriptorCollection
}

// PropertyDescriptor definition
type PropertyDescriptor struct {
	Name         string
	Description  string
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
	var (
		ti = t[i].PropertyType.Name
		tj = t[j].PropertyType.Name
		ni = t[i].Name
		nj = t[j].Name
	)

	isPrimaryKey := func(value string) bool {
		return strings.EqualFold(value, "id")
	}

	isForeignKey := func(value string) bool {
		return strings.HasSuffix(value, "_id")
	}

	if isPrimaryKey(ni) || isForeignKey(ni) {
		return true
	}

	if isPrimaryKey(nj) || isForeignKey(nj) {
		return false
	}

	if ti == tj {
		return ni < nj
	}

	return ti < tj
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
	In            string
	Description   string
	Style         string
	Explode       *bool
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
	Description     string
	Required        bool
	RequestBodyType *TypeDescriptor
}

// ResponseDescriptor definition
type ResponseDescriptor struct {
	Name         string
	Description  string
	ResponseType *TypeDescriptor
	Headers      HeaderDescriptorCollection
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

// ControllerDescriptor definition
type ControllerDescriptor struct {
	Name        string
	Description string
	Operations  OperationDescriptorCollection
}

// ControllerDescriptorMap definition
type ControllerDescriptorMap map[string]*ControllerDescriptor

// Get returns a descriptor
func (m ControllerDescriptorMap) Get(keys []string) *ControllerDescriptor {
	if len(keys) == 0 {
		keys = []string{"default"}
	}

	var (
		name           = keys[0]
		controller, ok = m[name]
	)

	if !ok {
		controller = &ControllerDescriptor{
			Name: name,
		}
	}

	m[name] = controller
	return controller
}

// Collection returns the descriptors as collection
func (m ControllerDescriptorMap) Collection() ControllerDescriptorCollection {
	descriptors := ControllerDescriptorCollection{}

	for _, descriptor := range m {
		descriptors = append(descriptors, descriptor)
	}

	return descriptors
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
