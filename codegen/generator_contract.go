package codegen

import "path/filepath"

// ContractGenerator generates a contract
type ContractGenerator struct {
	Path       string
	Collection TypeDescriptorCollection
}

// Generate generates the file
func (g *ContractGenerator) Generate() *File {
	builder := &FileBuilder{
		Package: "service",
	}

	// generate contract
	for _, descriptor := range g.Collection {
		switch {
		case descriptor.IsAlias:
			builder.
				Literal(descriptor.Name).
				Element(descriptor.Element.Name)
		case descriptor.IsArray:
			builder.
				Array(descriptor.Name).
				Element(descriptor.Element.Name)
		case descriptor.IsClass:
			builder.Type(descriptor.Name)
			// builder = &StructTypeBuilder{
			// 	Name:   descriptor.Name,
			// 	Fields: descriptor.Fields(),
			// }
		case descriptor.IsEnum:
			//TODO: implement enum builder
			continue
		}

		// builder.Commentf("%s is a struct type auto-generated from OpenAPI spec", descriptor.Name)
		// builder.Commentf(descriptor.Description)

		// file.Decls = append(file.Decls, builder.Build())
	}

	return &File{
		Name:    filepath.Join(g.Path, "contract.go"),
		Content: builder.Build(),
	}
}
