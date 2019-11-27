package service_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/stride/fake"
	"github.com/phogolabs/stride/service"
)

var _ = Describe("Validator", func() {
	var validator *service.Validator

	BeforeEach(func() {
		reporter := &fake.Reporter{}
		reporter.WithReturns(reporter)

		validator = &service.Validator{
			Path:     path("../fixture/spec/schemas-array.yaml"),
			Reporter: reporter,
		}
	})

	It("validates the spec successfully", func() {
		Expect(validator.Validate()).To(Succeed())
	})

	Context("when the file does not exists", func() {
		BeforeEach(func() {
			validator.Path = "./i-do-not-exist.yaml"
		})

		It("returns an error", func() {
			Expect(validator.Validate()).To(MatchError("open ./i-do-not-exist.yaml: no such file or directory"))
		})
	})

	Context("when the validation fails", func() {
		BeforeEach(func() {
			validator.Path = "../fixture/spec/schemas-object.yaml"
		})

		It("returns an error", func() {
			Expect(validator.Validate()).To(MatchError("message: Please check the error log for more details"))
		})
	})
})
