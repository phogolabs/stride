package golang_test

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGolang(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Golang Suite")
}

func tmpfile() string {
	tmp, err := ioutil.TempFile("", "stride")
	Expect(err).To(BeNil())

	name := tmp.Name()
	Expect(tmp.Close()).To(Succeed())
	Expect(os.Remove(name)).To(Succeed())
	return name
}

func tmpdir() string {
	dir, err := ioutil.TempDir("", "example")
	Expect(err).To(BeNil())
	Expect(os.Remove(dir)).To(Succeed())
	return dir
}
