package golang_test

import (
	"bytes"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/stride/plugin/golang"
)

var _ = Describe("Merge", func() {
	var merger *golang.Merger

	BeforeEach(func() {
		merger = &golang.Merger{}
	})

	FDescribe("Struct", func() {
		Context("when the struct is auto-generated", func() {
			BeforeEach(func() {
				target, err := golang.OpenFile("../../fixture/code/struct_generated_target.go.fixture")
				Expect(err).To(BeNil())

				source, err := golang.OpenFile("../../fixture/code/struct_generated_source.go.fixture")
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

		Context("when the struct is user-defined", func() {
			BeforeEach(func() {
				target, err := golang.OpenFile("../../fixture/code/struct_defined_target.go.fixture")
				Expect(err).To(BeNil())

				source, err := golang.OpenFile("../../fixture/code/struct_defined_source.go.fixture")
				Expect(err).To(BeNil())

				merger.Target = target
				merger.Source = source
			})

			It("appends the struct to the end of the file", func() {
				Expect(merger.Merge()).To(Succeed())

				target := &bytes.Buffer{}
				merger.Target.WriteTo(target)

				merged, err := ioutil.ReadFile("../../fixture/code/struct_defined_merged.go.fixture")
				Expect(err).To(BeNil())

				// fmt.Println(target.String())
				// fmt.Println("----------")
				// fmt.Println(string(merged))

				Expect(target.String()).To(Equal(string(merged)))
			})
		})
	})

	Describe("Function", func() {
		Context("when the struct is initialized by the user", func() {
			It("preservs the initialization", func() {
			})
		})
	})

	Context("when the range is in the beginning", func() {
		BeforeEach(func() {
			target, err := golang.OpenFile("../fixture/merge_top_target.go.fixture")
			Expect(err).To(BeNil())

			source, err := golang.OpenFile("../fixture/merge_top_source.go.fixture")
			Expect(err).To(BeNil())

			merger.Target = target
			merger.Source = source
		})

		It("merges the body successfully", func() {
			Expect(merger.Merge()).To(Succeed())
			merger.Target.WriteTo(os.Stdout)
		})
	})

	Context("when the range is in the end", func() {
		BeforeEach(func() {
			target, err := golang.OpenFile("../fixture/merge_bottom_target.go.fixture")
			Expect(err).To(BeNil())

			source, err := golang.OpenFile("../fixture/merge_bottom_source.go.fixture")
			Expect(err).To(BeNil())

			merger.Target = target
			merger.Source = source
		})

		It("merges the body successfully", func() {
			Expect(merger.Merge()).To(Succeed())
			merger.Target.WriteTo(os.Stdout)
		})
	})

	Context("when the range is in the middle", func() {
		BeforeEach(func() {
			target, err := golang.OpenFile("../fixture/merge_middle_target.go.fixture")
			Expect(err).To(BeNil())

			source, err := golang.OpenFile("../fixture/merge_middle_source.go.fixture")
			Expect(err).To(BeNil())

			merger.Target = target
			merger.Source = source
		})

		It("merges the body successfully", func() {
			Expect(merger.Merge()).To(Succeed())
			merger.Target.WriteTo(os.Stdout)
		})
	})
})
