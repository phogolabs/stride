package codegen

import (
	"bytes"
	"fmt"
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

// CollectFromParameters collect the type descriptors from parameters collection
func (m TypeDescriptorMap) CollectFromParameters(descriptors ParameterDescriptorCollection) {
	for _, descriptor := range descriptors {
		m.add(descriptor.ParameterType)
	}
}

// CollectFromRequests collect the type descriptors from requests collection
func (m TypeDescriptorMap) CollectFromRequests(descriptors RequestDescriptorCollection) {
	for _, descriptor := range descriptors {
		m.CollectFromParameters(descriptor.Parameters)
		m.add(descriptor.RequestType)
	}
}

// CollectFromResponses collect the type descriptors from responses collection
func (m TypeDescriptorMap) CollectFromResponses(descriptors ResponseDescriptorCollection) {
	for _, descriptor := range descriptors {
		m.CollectFromParameters(descriptor.Parameters)
		m.add(descriptor.ResponseType)
	}
}

// CollectFromSchemas collect the type descriptors from types collection
func (m TypeDescriptorMap) CollectFromSchemas(descriptors TypeDescriptorCollection) {
	for _, descriptor := range descriptors {
		m.add(descriptor)
	}
}

// CollectFromControllers collect the type descriptors from controllers collection
func (m TypeDescriptorMap) CollectFromControllers(descriptors ControllerDescriptorCollection) {
	for _, controller := range descriptors {
		for _, operation := range controller.Operations {
			m.CollectFromRequests(operation.Requests)
			m.CollectFromResponses(operation.Responses)
		}
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

// Tags returns the associated tagss
func (d *TypeDescriptor) Tags(required bool) []*Tag {
	var (
		tags []*Tag
		tag  *Tag
	)

	oneof := func(values []interface{}) string {
		buffer := &bytes.Buffer{}

		for index, value := range values {
			if index > 0 {
				fmt.Fprint(buffer, " ")
			}

			fmt.Fprintf(buffer, "%v", value)
		}

		return buffer.String()
	}

	// validation
	tag = &Tag{
		Key:  "validate",
		Name: "-",
	}

	if required {
		tag.Options = append(tag.Options, "required")
	}

	tags = append(tags, tag)

	for k, v := range d.Metadata {
		switch k {
		case "unique":
			if unique, ok := v.(bool); ok {
				if unique {
					tag.Options = append(tag.Options, "unique")
				}
			}
		case "pattern":
			if value, ok := v.(string); ok {
				if value != "" {
					// TODO: add support for regex
					// options = append(options, fmt.Sprintf("regex=%v", value))
				}
			}
		case "multiple_of":
			if value, ok := v.(*float64); ok {
				if value != nil {
					// TODO: add support for multileof
					// options = append(options, fmt.Sprintf("multipleof=%v", value))
				}
			}
		case "min":
			if value, ok := v.(*float64); ok {
				if value != nil {
					if exclusive, ok := d.Metadata["min_exclusive"].(bool); ok {
						if exclusive {
							tag.Options = append(tag.Options, fmt.Sprintf("gt=%v", *value))
						} else {
							tag.Options = append(tag.Options, fmt.Sprintf("gte=%v", *value))
						}
					}
				}
			}
		case "max":
			if value, ok := v.(*float64); ok {
				if value != nil {
					if exclusive, ok := d.Metadata["max_exclusive"].(bool); ok {
						if exclusive {
							tag.Options = append(tag.Options, fmt.Sprintf("lt=%v", *value))
						} else {
							tag.Options = append(tag.Options, fmt.Sprintf("lte=%v", *value))
						}
					}
				}
			}
		case "values":
			if values, ok := v.([]interface{}); ok {
				if len(values) > 0 {
					tag.Options = append(tag.Options, fmt.Sprintf("oneof=%v", oneof(values)))
				}
			}
		}
	}

	if len(tag.Options) > 0 {
		tag.Name = tag.Options[0]
		tag.Options = tag.Options[1:]
	}

	// default
	if value := d.Default; value != nil {
		tag = &Tag{
			Key:  "default",
			Name: fmt.Sprintf("%v", value),
		}

		tags = append(tags, tag)
	}

	return tags
}

// Kind returns the golang kind
func (d *TypeDescriptor) Kind() string {
	name := strings.ToLower(d.Name)

	switch {
	case name == "date-time":
		return "time.Time"
	case name == "date":
		return "time.Time"
	case name == "uuid":
		return "schema.UUID"
	case d.IsPrimitive:
		return name
	}

	name = camelize(name)

	if item := element(d); item.IsNullable {
		name = fmt.Sprintf("*%s", name)
	}

	return name
}

// HasProperties returns true if the type has properties
func (d *TypeDescriptor) HasProperties() bool {
	return len(d.Properties) > 0
}

// PropertyDescriptor definition
type PropertyDescriptor struct {
	Name         string
	Description  string
	Required     bool
	ReadOnly     bool
	WriteOnly    bool
	PropertyType *TypeDescriptor
}

// Tags returns the underlying tags
func (p *PropertyDescriptor) Tags() []*Tag {
	var (
		tags = []*Tag{}
		tag  *Tag
	)

	omitempty := func() []string {
		options := []string{}

		if !p.Required {
			options = append(options, "omitempty")
		}

		return options
	}

	// json marshalling
	tag = &Tag{
		Key:     "json",
		Name:    p.Name,
		Options: omitempty(),
	}

	tags = append(tags, tag)

	// validation
	tags = append(tags, p.PropertyType.Tags(p.Required)...)

	return tags
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

	return ni < nj
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
	Explode       bool
	Required      bool
	Deprecated    bool
	ParameterType *TypeDescriptor
}

// Tags returns the tags for this parameter
func (p *ParameterDescriptor) Tags() []*Tag {
	var (
		tags = []*Tag{}
		tag  *Tag
	)

	// style
	tag = &Tag{
		Key:  p.In,
		Name: p.Name,
	}

	if style := strings.ToLower(p.Style); style != "" {
		tag.Options = append(tag.Options, style)
	}

	if exlode := p.Explode; exlode {
		tag.Options = append(tag.Options, "explode")
	}

	tags = append(tags, tag)

	// validation
	tags = append(tags, p.ParameterType.Tags(p.Required)...)

	return tags
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

// RequestDescriptor definition
type RequestDescriptor struct {
	ContentType string
	Description string
	Required    bool
	RequestType *TypeDescriptor
	Parameters  ParameterDescriptorCollection
}

// RequestDescriptorCollection definition
type RequestDescriptorCollection []*RequestDescriptor

// Len is the number of elements in the collection.
func (t RequestDescriptorCollection) Len() int {
	return len(t)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (t RequestDescriptorCollection) Less(i, j int) bool {
	return t[i].ContentType < t[j].ContentType
}

// Swap swaps the elements with indexes i and j.
func (t RequestDescriptorCollection) Swap(i, j int) {
	var x = t[i]
	t[i] = t[j]
	t[j] = x
}

// ResponseDescriptor definition
type ResponseDescriptor struct {
	Code         int
	Description  string
	ContentType  string
	ResponseType *TypeDescriptor
	Parameters   ParameterDescriptorCollection
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
	var (
		x = t[i]
		y = t[j]
	)

	if x.ContentType == y.ContentType {
		return x.Code < y.Code
	}

	return x.ContentType < y.ContentType
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

// OperationDescriptor definition
type OperationDescriptor struct {
	Method      string
	Path        string
	Name        string
	Summary     string
	Description string
	Deprecated  bool
	Tags        []string
	Requests    RequestDescriptorCollection
	Responses   ResponseDescriptorCollection
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

func element(descriptor *TypeDescriptor) *TypeDescriptor {
	element := descriptor

	for element.IsAlias {
		element = element.Element
	}

	return element
}
