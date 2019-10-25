package codegen_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/phogolabs/stride/codegen"
)

func TestCodegen(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Codegen Suite")
}

func resolve(name string) *codegen.SpecDescriptor {
	var (
		path     = fmt.Sprintf("../fixture/%s", name)
		loader   = openapi3.NewSwaggerLoader()
		resolver = &codegen.Resolver{
			Cache: make(map[string]*codegen.TypeDescriptor),
		}
	)

	spec, err := loader.LoadSwaggerFromFile(path)
	Expect(err).To(BeNil())
	Expect(spec).NotTo(BeNil())

	return resolver.Resolve(spec)
}
