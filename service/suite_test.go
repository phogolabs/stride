package service_test

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

func path(path string) string {
	source, err := os.Open(path)
	Expect(err).To(BeNil())

	target, err := ioutil.TempFile("", "ginkgo")
	Expect(err).To(BeNil())

	_, err = io.Copy(target, source)
	Expect(err).To(BeNil())

	return target.Name()
}
