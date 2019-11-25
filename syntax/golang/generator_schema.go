package golang

import (
	"fmt"
	"path/filepath"

	"github.com/phogolabs/stride/codedom"
	"github.com/phogolabs/stride/inflect"
)

// SchemaGenerator generates a contract
type SchemaGenerator struct {
	Path       string
	Collection codedom.TypeDescriptorCollection
}

// Generate generates the file
func (g *SchemaGenerator) Generate() *File {
	var (
		filename = filepath.Join(g.Path, "schema.go")
		root     = NewFile(filename)
	)

	// generate contract
	for _, descriptor := range g.Collection {
		switch {
		case descriptor.IsAlias:
			spec := NewLiteralType(descriptor.Name).Element(inflect.Unpointer(descriptor.Element.Kind()))
			spec.Commentf(descriptor.Description)
			// add the spec the file
			root.AddNode(spec)
		case descriptor.IsArray:
			spec := NewArrayType(descriptor.Name).Element(descriptor.Element.Kind())
			spec.Commentf(descriptor.Description)
			// add the spec the file
			root.AddNode(spec)
		case descriptor.IsClass:
			spec := NewStructType(descriptor.Name)
			spec.Commentf(descriptor.Description)
			// add the spec the file
			root.AddNode(spec)

			// add fields
			for _, property := range descriptor.Properties {
				var (
					tags = property.Tags()
					kind = property.PropertyType.Kind()
				)

				// add a import if needed
				root.AddImport(property.PropertyType.Namespace())
				// add the field
				spec.AddField(property.Name, kind, tags...)
			}
		case descriptor.IsEnum:
			spec := NewLiteralType(descriptor.Name).Element("string")
			spec.Commentf(descriptor.Description)
			// add the spec the file
			root.AddNode(spec)

			block := NewConstBlockType()
			// add the spec the file
			root.AddNode(block)

			if values, ok := descriptor.Metadata["values"].([]interface{}); ok {
				for _, item := range values {
					var (
						value = fmt.Sprintf("%v", item)
						name  = inflect.Camelize(spec.Name(), value)
					)

					block.AddConst(name, spec.Name(), value)
				}
			}
			continue
		}
	}

	return root
}
