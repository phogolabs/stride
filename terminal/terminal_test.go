package terminal_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gbytes"
	"github.com/phogolabs/stride/contract"
	"github.com/phogolabs/stride/terminal"
)

var _ = Describe("Terminal", func() {
	var (
		buffer   *gbytes.Buffer
		reporter *terminal.Reporter
	)

	BeforeEach(func() {
		buffer = gbytes.NewBuffer()
		reporter = &terminal.Reporter{
			Writer: buffer,
		}
	})

	ItWritesAMessage := func(severity contract.Severity) {
		Describe("Notice", func() {
			It("writes a message", func() {
				reporter.With(severity).Notice("hello")
				Expect(buffer).To(gbytes.Say("hello"))
			})
		})

		Describe("Info", func() {
			It("writes a message", func() {
				reporter.With(severity).Info("hello")
				Expect(buffer).To(gbytes.Say("hello"))
			})
		})

		Describe("Success", func() {
			It("writes a message", func() {
				reporter.With(severity).Success("hello")
				Expect(buffer).To(gbytes.Say("hello"))
			})
		})

		Describe("Warn", func() {
			It("writes a message", func() {
				reporter.With(severity).Warn("hello")
				Expect(buffer).To(gbytes.Say("hello"))
			})
		})

		Describe("Error", func() {
			It("writes a message", func() {
				reporter.With(severity).Error("hello")
				Expect(buffer).To(gbytes.Say("hello"))
			})
		})
	}

	ItWritesAMessage(contract.SeverityLow)
	ItWritesAMessage(contract.SeverityNormal)
	ItWritesAMessage(contract.SeverityHigh)
	ItWritesAMessage(contract.SeverityVeryHigh)
})
