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
			root.
				Literal(descriptor.Name).
				Element(inflect.Unpointer(descriptor.Element.Kind())).
				Commentf(descriptor.Description)
		case descriptor.IsArray:
			root.
				Array(descriptor.Name).
				Element(descriptor.Element.Kind()).
				Commentf(descriptor.Description)
		case descriptor.IsClass:
			builder := root.Struct(descriptor.Name)

			// add fields
			for _, property := range descriptor.Properties {
				var (
					tags = property.Tags()
					kind = property.PropertyType.Kind()
				)

				// add a import if needed
				root.AddImport(property.PropertyType.Namespace())
				// add the field
				builder.AddField(property.Name, kind, tags...)
			}
		case descriptor.IsEnum:
			builder := root.
				Literal(descriptor.Name).
				Element("string")

			builder.Commentf(descriptor.Description)

			block := root.Const()

			if values, ok := descriptor.Metadata["values"].([]interface{}); ok {
				for _, item := range values {
					var (
						value = fmt.Sprintf("%v", item)
						name  = inflect.Camelize(builder.Name(), value)
					)

					block.AddConst(name, builder.Name(), value)
				}
			}
			continue
		}
	}

	return root
}
