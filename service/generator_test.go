package service_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/stride/codedom"
	"github.com/phogolabs/stride/fake"
	"github.com/phogolabs/stride/service"
)

var _ = Describe("Generator", func() {
	var (
		generator *service.Generator
		resolver  *fake.SpecResolver
		coder     *fake.CodeGenerator
	)

	BeforeEach(func() {
		resolver = &fake.SpecResolver{}
		resolver.ResolveReturns(&codedom.SpecDescriptor{}, nil)

		coder = &fake.CodeGenerator{}

		generator = &service.Generator{
			Path:      path("../fixture/spec/schemas-array.yaml"),
			Generator: coder,
			Resolver:  resolver,
		}
	})

	It("generates the project successfully", func() {
		Expect(generator.Generate()).To(Succeed())
		Expect(resolver.ResolveCallCount()).To(Equal(1))
		Expect(coder.GenerateCallCount()).To(Equal(1))
	})

	Context("when the code generator fails", func() {
		BeforeEach(func() {
			coder.GenerateReturns(fmt.Errorf("oh no"))
		})

		It("returns an error", func() {
			Expect(generator.Generate()).To(MatchError("oh no"))
		})
	})

	Context("when the file does not exists", func() {
		BeforeEach(func() {
			generator.Path = "./i-do-not-exist.yaml"
		})

		It("returns an error", func() {
			Expect(generator.Generate()).To(MatchError("open ./i-do-not-exist.yaml: no such file or directory"))
			Expect(resolver.ResolveCallCount()).To(BeZero())
			Expect(coder.GenerateCallCount()).To(BeZero())
		})
	})
})
