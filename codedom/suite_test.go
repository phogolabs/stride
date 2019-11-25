package codedom_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/phogolabs/stride/codedom"
)

type GetTypeDescriptorFunc func() *codedom.TypeDescriptor

func TestCodegen(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Codedom Suite")
}

func resolve(name string) *codedom.SpecDescriptor {
	var (
		path     = fmt.Sprintf("../fixture/spec/%s", name)
		loader   = openapi3.NewSwaggerLoader()
		resolver = codedom.NewResolver()
	)

	spec, err := loader.LoadSwaggerFromFile(path)
	Expect(err).To(BeNil())
	Expect(spec).NotTo(BeNil())

	return resolver.Resolve(spec)
}
