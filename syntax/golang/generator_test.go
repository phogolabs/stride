package golang_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/stride/codegen"
	"github.com/phogolabs/stride/syntax/golang"
)

var _ = Describe("Generator", func() {
	var generator *golang.Generator

	BeforeEach(func() {
		generator = &golang.Generator{
			Path: tmpdir(),
		}
	})

	It("generates the package successfully", func() {
		descriptor := &codegen.ControllerDescriptor{
			Name: "user",
			Operations: codegen.OperationDescriptorCollection{
				&codegen.OperationDescriptor{
					Method: "GET",
					Path:   "/accounts",
					Name:   "get-accounts",
				},
			},
		}

		spec := &codegen.SpecDescriptor{}
		spec.Controllers = append(spec.Controllers, descriptor)

		Expect(generator.Generate(spec)).To(Succeed())
		Expect(generator.Path).To(BeADirectory())
		Expect(generator.Path + "/service").To(BeADirectory())
	})

	Context("when cannot create the directory", func() {
		BeforeEach(func() {
			generator.Path = "/my-dir"
		})

		It("returns the error", func() {
			spec := &codegen.SpecDescriptor{}
			Expect(generator.Generate(spec)).To(MatchError("mkdir /my-dir: read-only file system"))
			Expect(generator.Path).NotTo(BeADirectory())
		})
	})
})
