package golang_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/stride/plugin/golang"
)

var _ = PDescribe("Merge", func() {
	var merger *golang.Merger

	BeforeEach(func() {
		merger = &golang.Merger{}
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

	FContext("when the range is in the middle", func() {
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
