package codegen

import "path/filepath"

// ContractGenerator generates a contract
type ContractGenerator struct {
	Path       string
	Collection TypeDescriptorCollection
}

// Generate generates the file
func (g *ContractGenerator) Generate() *File {
	root := NewFile(filepath.Join(g.Path, "contract.go"))

	// generate contract
	for _, descriptor := range g.Collection {
		switch {
		case descriptor.IsAlias:
			root.
				Literal(descriptor.Name).
				Element(descriptor.Element.Name).
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
				builder.Field(property.Name, kind, tags...)
			}
		case descriptor.IsEnum:
			//TODO: implement enum builder
			continue
		}
	}

	return root
}
