package codegen

import (
	"fmt"
	"path/filepath"
)

// ContractGenerator generates a contract
type ContractGenerator struct {
	Path       string
	Collection TypeDescriptorCollection
}

// Generate generates the file
func (g *ContractGenerator) Generate() *File {
	root := &FileBuilder{
		Package: "service",
	}

	// generate contract
	for _, descriptor := range g.Collection {
		var parent Builder

		fmt.Println(descriptor.Name)

		switch {
		case descriptor.IsAlias:
			builder := root.Literal(descriptor.Name)
			builder.Element(descriptor.Element.Name)
			parent = builder
		case descriptor.IsArray:
			builder := root.Array(descriptor.Name)
			builder.Element(descriptor.Element.Name)
			parent = builder
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

	return &File{
		Name:    filepath.Join(g.Path, "contract.go"),
		Content: root.Build(),
	}
}
