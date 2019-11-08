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

	// create the file
	file := builder.Build()

	// generate contract
	for _, descriptor := range g.Collection {
		var builder Builder

		switch {
		case descriptor.IsAlias:
			builder = &LiteralTypeBuilder{
				Name:    descriptor.Name,
				Element: descriptor.Element.Name,
			}
		case descriptor.IsArray:
			builder = &ArrayTypeBuilder{
				Name:    descriptor.Name,
				Element: descriptor.Element.Name,
			}
		case descriptor.IsClass:
			builder = &StructTypeBuilder{
				Name:   descriptor.Name,
				Fields: descriptor.Fields(),
			}
		case descriptor.IsEnum:
			//TODO: implement enum builder
			continue
		}

		builder.Commentf("%s is a struct type auto-generated from OpenAPI spec", descriptor.Name)
		builder.Commentf(descriptor.Description)
		file.Decls = append(file.Decls, builder.Build())
	}

	return &File{
		Name:    filepath.Join(g.Path, "contract.go"),
		Content: file,
	}
}
