package torrent_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/stride/torrent"
)

var _ = Describe("Torrent", func() {
	Describe("GetAsync", func() {
		It("downloads the file successfully", func() {
			task, err := torrent.GetAsync("../fixture/spec/headers-array.yaml")
			Expect(err).To(BeNil())
			Expect(task.Wait()).To(Succeed())
			Expect(task.Data()).To(BeAnExistingFile())
		})

		Context("when the file does not exist", func() {
			It("returns an error", func() {
				task, err := torrent.GetAsync("./i-dont-exist.yaml")
				Expect(err).To(BeNil())
				Expect(task.Wait()).To(HaveOccurred())
				Expect(task.Data()).NotTo(BeAnExistingFile())
			})
		})
	})
})
