package service_test

import (
	"io"
	"io/ioutil"
	"net"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	_ "github.com/phogolabs/stride/template"
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

func wait(addr string) {
	Eventually(func() error {
		_, err := net.Dial("tcp", addr)
		return err
	}).Should(Succeed())
}
