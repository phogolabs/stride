package codegen

import (
	"path/filepath"
)

// ContractGenerator generates a contract
type ContractGenerator struct {
	Path       string
	Collection TypeDescriptorCollection
}

// Generate generates the file
func (g *ContractGenerator) Generate() *File {
	root := NewFileBuilder("service")

	// generate contract
	for _, descriptor := range g.Collection {
		var parent Builder

		switch {
		case descriptor.IsAlias:
			parent = root.
				Literal(descriptor.Name).
				Element(descriptor.Element.Name)
		case descriptor.IsArray:
			parent = root.
				Array(descriptor.Name).
				Element(descriptor.Element.Kind())
		case descriptor.IsClass:
			builder := root.Type(descriptor.Name)
			parent = builder

			// add fields
			for _, property := range descriptor.Properties {
				var (
					tags = property.Tags()
					kind = property.PropertyType.Kind()
				)
				builder.Field(property.Name, kind, tags...)
			}
		case descriptor.IsEnum:
			//TODO: implement enum builder
			continue
		}

		parent.Commentf("%s is a struct type auto-generated from OpenAPI spec", parent.Name())
		parent.Commentf(descriptor.Description)
	}

	return root.Build(filepath.Join(g.Path, "contract.go"))
}
