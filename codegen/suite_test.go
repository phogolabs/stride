package codegen_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/phogolabs/stride/codegen"
)

type GetTypeDescriptorFunc func() *codegen.TypeDescriptor

func TestCodegen(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Codegen Suite")
}

func resolve(name string) *codegen.SpecDescriptor {
	var (
		path     = fmt.Sprintf("../fixture/spec/%s", name)
		loader   = openapi3.NewSwaggerLoader()
		resolver = codegen.NewResolver()
	)

	spec, err := loader.LoadSwaggerFromFile(path)
	Expect(err).To(BeNil())
	Expect(spec).NotTo(BeNil())

	return resolver.Resolve(spec)
}

func tmpfile() string {
	tmp, err := ioutil.TempFile("", "stride")
	Expect(err).To(BeNil())

	name := tmp.Name()
	Expect(tmp.Close()).To(Succeed())
	Expect(os.Remove(name)).To(Succeed())
	return name
}

func tmpdir() string {
	dir, err := ioutil.TempDir("", "example")
	Expect(err).To(BeNil())
	Expect(os.Remove(dir)).To(Succeed())
	return dir
}
