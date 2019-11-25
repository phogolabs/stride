package codegen

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/fatih/structtag"
	"github.com/phogolabs/stride/inflect"
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

// Add adds a type descriptor
func (m TypeDescriptorMap) Add(descriptor *TypeDescriptor) {
	m[descriptor.Name] = descriptor
}

// Get returns the TypeDescriptor for given name
func (m TypeDescriptorMap) Get(name string) *TypeDescriptor {
	if descriptor, ok := m[name]; ok {
		return descriptor
	}

	return nil
}

// Clear clears the map
func (m TypeDescriptorMap) Clear() {
	for k := range m {
		delete(m, k)
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
	reflect.Swapper(t)(i, j)
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
func (d *TypeDescriptor) Tags(required bool) TagDescriptorCollection {
	var (
		tags = TagDescriptorCollection{}
		tag  *TagDescriptor
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
	tag = &TagDescriptor{
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
		tag = &TagDescriptor{
			Key:  "default",
			Name: fmt.Sprintf("%v", value),
		}

		tags = append(tags, tag)
	}

	return tags
}

// Namespace returns the namespace
func (d *TypeDescriptor) Namespace() string {
	switch d.Kind() {
	case "time.Time":
		return "time"
	case "schema.UUID":
		return "github.com/phogolabs/schema"
	default:
		return ""
	}
}

// Kind returns the golang kind
func (d *TypeDescriptor) Kind() string {
	name := strings.ToLower(d.Name)

	switch {
	case name == "date-time":
		name = "time.Time"
	case name == "date":
		name = "time.Time"
	case name == "uuid":
		name = "schema.UUID"
	default:
		if !d.IsPrimitive {
			name = inflect.Camelize(d.Name)
		}
	}

	if item := element(d); item.IsNullable {
		name = inflect.Pointer(name)
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
func (p *PropertyDescriptor) Tags() TagDescriptorCollection {
	var (
		tags = TagDescriptorCollection{}
		tag  *TagDescriptor
	)

	omitempty := func() []string {
		options := []string{}

		if !p.Required {
			options = append(options, "omitempty")
		}

		return options
	}

	// json marshalling
	tag = &TagDescriptor{
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
	reflect.Swapper(t)(i, j)
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
func (p *ParameterDescriptor) Tags() TagDescriptorCollection {
	var (
		tags = TagDescriptorCollection{}
		tag  *TagDescriptor
	)

	// style
	tag = &TagDescriptor{
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
	reflect.Swapper(t)(i, j)
}

// RequestDescriptor definition
type RequestDescriptor struct {
	ContentType string
	Description string
	Required    bool
	RequestType *TypeDescriptor
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
	reflect.Swapper(t)(i, j)
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
	reflect.Swapper(t)(i, j)
}

// ControllerDescriptor definition
type ControllerDescriptor struct {
	Name        string
	Description string
	Operations  OperationDescriptorCollection
}

// ControllerDescriptorMap definition
type ControllerDescriptorMap map[string]*ControllerDescriptor

// Add adds a descriptor to the map
func (m ControllerDescriptorMap) Add(descriptor *ControllerDescriptor) {
	if _, ok := m[descriptor.Name]; ok {
		return
	}

	m[descriptor.Name] = descriptor
}

// Get returns a descriptor
func (m ControllerDescriptorMap) Get(key string) *ControllerDescriptor {
	descriptor, ok := m[key]

	if !ok {
		descriptor = &ControllerDescriptor{
			Name: key,
		}

		m[key] = descriptor
	}

	return descriptor
}

// Clear clears the map
func (m ControllerDescriptorMap) Clear() {
	for k := range m {
		delete(m, k)
	}
}

// Collection returns the descriptors as collection
func (m ControllerDescriptorMap) Collection() ControllerDescriptorCollection {
	descriptors := ControllerDescriptorCollection{}

	for _, descriptor := range m {
		// sort the operations
		sort.Sort(descriptor.Operations)

		descriptors = append(descriptors, descriptor)
	}

	// sort the descriptors
	sort.Sort(descriptors)

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
	reflect.Swapper(t)(i, j)
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
	reflect.Swapper(t)(i, j)
}

// TagDescriptor represents a tag
type TagDescriptor struct {
	Key     string
	Name    string
	Options []string
}

// TagDescriptorCollection represents a field tag list
type TagDescriptorCollection []*TagDescriptor

func (tags TagDescriptorCollection) String() string {
	builder := &structtag.Tags{}

	for _, descriptor := range tags {
		builder.Set(&structtag.Tag{
			Key:     descriptor.Key,
			Name:    descriptor.Name,
			Options: descriptor.Options,
		})
	}

	if value := builder.String(); value != "" {
		return fmt.Sprintf("`%s`", value)
	}

	return ""
}

func element(descriptor *TypeDescriptor) *TypeDescriptor {
	element := descriptor

	for element.IsAlias {
		element = element.Element
	}

	return element
}
