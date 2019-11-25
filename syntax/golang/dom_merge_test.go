package syntax_test

import (
	"bytes"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/stride/plugin/syntax"
)

var _ = Describe("Merge", func() {
	var merger *syntax.Merger

	BeforeEach(func() {
		merger = &syntax.Merger{}
	})

	Context("when the struct is auto-generated", func() {
		BeforeEach(func() {
			target, err := syntax.OpenFile("../../fixture/code/struct_generated_target.go.fixture")
			Expect(err).To(BeNil())

			source, err := syntax.OpenFile("../../fixture/code/struct_generated_source.go.fixture")
			Expect(err).To(BeNil())

			merger.Target = target
			merger.Source = source
		})

		It("appends the user-defined fields", func() {
			Expect(merger.Merge()).To(Succeed())

			target := &bytes.Buffer{}
			merger.Target.WriteTo(target)

			merged, err := ioutil.ReadFile("../../fixture/code/struct_generated_merged.go.fixture")
			Expect(err).To(BeNil())

			Expect(target.String()).To(Equal(string(merged)))
		})
	})

	Context("when the user defined declarations", func() {
		BeforeEach(func() {
			target, err := syntax.OpenFile("../../fixture/code/user_defined_target.go.fixture")
			Expect(err).To(BeNil())

			source, err := syntax.OpenFile("../../fixture/code/user_defined_source.go.fixture")
			Expect(err).To(BeNil())

			merger.Target = target
			merger.Source = source
		})

		It("preserves the user definitions", func() {
			Expect(merger.Merge()).To(Succeed())

			target := &bytes.Buffer{}
			merger.Target.WriteTo(target)

			merged, err := ioutil.ReadFile("../../fixture/code/user_defined_merged.go.fixture")
			Expect(err).To(BeNil())

			Expect(target.String()).To(Equal(string(merged)))
		})
	})

	Context("when the function has user-defined body", func() {
		BeforeEach(func() {
			target, err := syntax.OpenFile("../../fixture/code/function_body_target.go.fixture")
			Expect(err).To(BeNil())

			source, err := syntax.OpenFile("../../fixture/code/function_body_source.go.fixture")
			Expect(err).To(BeNil())

			merger.Target = target
			merger.Source = source
		})

		It("preserves the body", func() {
			Expect(merger.Merge()).To(Succeed())

			target := &bytes.Buffer{}
			merger.Target.WriteTo(target)

			merged, err := ioutil.ReadFile("../../fixture/code/function_body_merged.go.fixture")
			Expect(err).To(BeNil())

			Expect(target.String()).To(Equal(string(merged)))
		})
	})
})
