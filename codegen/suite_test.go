package codegen_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestCodegen(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Codegen Suite")
}

func load(path string) *openapi3.Swagger {
	var loader = openapi3.NewSwaggerLoader()

	spec, err := loader.LoadSwaggerFromFile(path)
	Expect(err).To(BeNil())
	Expect(spec).NotTo(BeNil())
	return spec
}
