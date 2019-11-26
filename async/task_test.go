package async_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/stride/async"
)

var _ = Describe("Task", func() {
	Describe("NewTask", func() {
		It("creates a task successfully", func() {
			runner := &Runner{}

			task := async.NewTask(runner.Run, "root")
			Expect(task).NotTo(BeNil())
		})
	})

	Describe("Run", func() {
		It("starts a task successfully", func() {
			runner := &Runner{}

			task := async.NewTask(runner.Run, "root")
			task.Run()

			Eventually(runner.Execution).Should(Equal(1))
		})
	})

	Describe("Data", func() {
		It("returns the associated data", func() {
			runner := &Runner{}

			task := async.NewTask(runner.Run, "root")
			Expect(task.Data()).To(Equal("root"))
		})

		Context("when no data is provided", func() {
			It("returns nil", func() {
				runner := &Runner{}

				task := async.NewTask(runner.Run)
				Expect(task.Data()).To(BeNil())
			})
		})

		Context("when more than one item is provided", func() {
			It("returns the associated data", func() {
				runner := &Runner{}

				task := async.NewTask(runner.Run, "root", "guest")
				Expect(task.Data()).To(HaveLen(2))
				Expect(task.Data()).To(ContainElement("root"))
				Expect(task.Data()).To(ContainElement("guest"))
			})
		})
	})

	Describe("Wait", func() {
		It("waits for the function successfully", func() {
			runner := &Runner{}

			task := async.NewTask(runner.Run, "root")
			task.Run()

			Expect(task.Wait()).To(Succeed())
			Expect(runner.Count).To(Equal(1))
		})

		Context("when the runner fails", func() {
			It("returns an error", func() {
				runner := &Runner{
					Error: fmt.Errorf("oh no"),
				}

				task := async.NewTask(runner.Run, "root")
				task.Run()

				Expect(task.Wait()).To(MatchError("oh no"))
				Expect(runner.Count).To(Equal(1))
			})
		})
	})

	Describe("Stop", func() {
		It("stops the function successfully", func() {
			runner := &Runner{}

			task := async.NewTask(runner.Run, "root")
			task.Run()

			Expect(task.Stop()).To(Succeed())
			Expect(runner.Count).To(Equal(1))
		})
	})

	Context("when the runner fails", func() {
		It("returns an error", func() {
			runner := &Runner{
				Error: fmt.Errorf("oh no"),
			}

			task := async.NewTask(runner.Run, "root")
			task.Run()

			Expect(task.Stop()).To(MatchError("oh no"))
			Expect(runner.Count).To(Equal(1))
		})
	})
})
