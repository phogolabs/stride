package codegen_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/stride/codegen"
)

var _ = FDescribe("Merge", func() {
	var merger *codegen.Merger

	BeforeEach(func() {
		merger = &codegen.Merger{}
	})

	Context("when the range is in the beginning", func() {
		BeforeEach(func() {
			target, err := codegen.Open("../fixture/merge_top_target.go.fixture")
			Expect(err).To(BeNil())

			source, err := codegen.Open("../fixture/merge_top_source.go.fixture")
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
			target, err := codegen.Open("../fixture/merge_bottom_target.go.fixture")
			Expect(err).To(BeNil())

			source, err := codegen.Open("../fixture/merge_bottom_source.go.fixture")
			Expect(err).To(BeNil())

			merger.Target = target
			merger.Source = source
		})

		It("merges the body successfully", func() {
			Expect(merger.Merge()).To(Succeed())
			merger.Target.WriteTo(os.Stdout)
		})
	})

	FContext("when the range is in the middle", func() {
		BeforeEach(func() {
			target, err := codegen.Open("../fixture/merge_middle_target.go.fixture")
			Expect(err).To(BeNil())

			source, err := codegen.Open("../fixture/merge_middle_source.go.fixture")
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
