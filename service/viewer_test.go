package service_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/stride/service"
)

var _ = Describe("Viewer", func() {
	var (
		server *http.Server
		config *service.ViewerConfig
	)

	BeforeEach(func() {
		config = &service.ViewerConfig{
			Addr: ":8080",
			Path: path("../fixture/spec/schemas-array.yaml"),
		}

		server = service.NewViewer(config)
		go server.ListenAndServe()
		wait(config.Addr)
	})

	AfterEach(func() {
		Expect(server.Shutdown(context.TODO())).To(Succeed())
	})

	Context("GET /swagger.spec", func() {
		It("returns the spec successfully", func() {
			response, err := http.Get("http://127.0.0.1:8080/swagger.spec")
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
		})
	})

	Context("GET /*", func() {
		It("returns the assets", func() {
			response, err := http.Get("http://127.0.0.1:8080/")
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
		})
	})
})
