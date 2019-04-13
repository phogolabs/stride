package codegen

import (
	"fmt"
	"io/ioutil"

	"github.com/aymerick/raymond"
)

// Generator generates the source code
type Generator struct{}

// Generate generates the source code
func (g *Generator) Generate(spec *SpecDescriptor) error {
	template, err := ioutil.ReadFile("./template/codegen/model.go.mustache")
	if err != nil {
		return err
	}

	result, err := raymond.Render(string(template), spec)
	if err != nil {
		return err
	}

	fmt.Println(result)

	return nil
}
