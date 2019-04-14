package codegen

import (
	"fmt"
	"path/filepath"

	"github.com/go-openapi/inflect"
)

// Mapper maps the descriptor to a model
type Mapper struct {
	Types map[string]*TypeModel
}

// Map does the mapping
func (m *Mapper) Map(spec *SpecDescriptor) {
	for _, schema := range spec.Schemas {
		m.mapType("schema", schema)
	}

	for _, header := range spec.Headers {
		m.mapType("header", header.HeaderType)
	}

	for _, parameter := range spec.Parameters {
		m.mapType("parameter", parameter.ParameterType)
	}

	for _, request := range spec.RequestBodies {
		for _, content := range request.Contents {
			// NOTE: Looks like the spec returns only one content
			// It's safe to use the request's name as a parent rather than
			// request's name + content-type
			name := request.Name + "_request"
			m.mapType(name, content.ContentType)
		}
	}

	for _, response := range spec.Responses {
		for _, content := range response.Contents {
			// NOTE: Looks like the spec returns only one content
			// It's safe to use the response's name as a parent rather than
			// response's name + content-type
			name := response.Name + "_response"
			m.mapType(name, content.ContentType)
		}
	}
}

func (m *Mapper) mapType(path string, schema *TypeDescriptor) {
	if schema == nil {
		fmt.Println("skip", path)
		return
	}

	if schema.IsPrimitive {
		return
	}

	if schema.Name == "" {
		schema.Name = path
	}

	//TODO: take care for types where we added suffix
	schema.Name = inflect.Camelize(filepath.Base(schema.Name))

	if _, ok := m.Types[schema.Name]; ok {
		return
	}

	m.Types[schema.Name] = nil

	for _, property := range schema.Properties {
		name := inflect.Camelize(schema.Name + "_" + property.Name)
		m.mapType(name, property.PropertyType)
	}
}
