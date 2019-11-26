package async_test

import (
	"context"
	"sync"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAsync(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Async Suite")
}

type Runner struct {
	Error error
	Count int
	Mutex sync.RWMutex
}

func (r *Runner) Run(ctx context.Context) error {
	r.Mutex.Lock()
	r.Count++
	r.Mutex.Unlock()

	Expect(ctx).NotTo(BeNil())
	return r.Error
}

func (r *Runner) Execution() int {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	return r.Count
}
