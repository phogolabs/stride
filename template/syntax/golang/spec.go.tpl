package service_test

import (
	"github.com/go-chi/chi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"{{ .project }}"
)

var _ = Describe("{{ .receiver }}", func() {
	var (
		router chi.Router
		// TODO: Uncomment if you are going to test your code
		// request  *http.Request
		// recorder *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		router = chi.NewRouter()

		controller := &service.CustomerAPI{}
		controller.Mount(router)

		Expect(router.Routes()).NotTo(BeEmpty())

		// TODO: Uncomment if you are going to test your code
		// recorder = httptest.NewRecorder()
	})

	{{ range .operations }}

	Describe("{{ .Method | uppercase }} {{ .Path }}", func() {
		// TODO: Implement the test cases for {{ .Name | camelize }} operation
	})
	{{ end }}
})
