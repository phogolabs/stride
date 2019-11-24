package service_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/stride/service"
)

var _ = Describe("Editor", func() {
	var (
		server *http.Server
		config *service.EditorConfig
	)

	BeforeEach(func() {
		config = &service.EditorConfig{
			Addr: ":8080",
			Path: path("../fixture/spec/schemas-array.yaml"),
		}

	})

	JustBeforeEach(func() {
		server = service.NewEditor(config)
		go server.ListenAndServe()
	})

	AfterEach(func() {
		Expect(server.Shutdown(context.TODO())).To(Succeed())
	})

	Context("POST /swagger.spec", func() {
		It("saves the spec successfully", func() {
			response, err := http.Post("http://127.0.0.1:8080/swagger.spec", "text/plain", bytes.NewBufferString("hello"))
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))

			data, err := ioutil.ReadFile(config.Path)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte("hello")))
		})

		Context("when the file cannot be created", func() {
			BeforeEach(func() {
				config.Path = "./file-directory/i-do-not-exist.yaml"
			})

			It("returns an error", func() {
				response, err := http.Post("http://127.0.0.1:8080/swagger.spec", "text/plain", bytes.NewBufferString("hello"))
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(500))
			})
		})
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
