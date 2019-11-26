package torrent_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTorrent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Torrent Suite")
}
