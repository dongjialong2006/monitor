package monitor

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMonitor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Monitor Suite")
}

var _ = Describe("Monitor", func() {
	Specify("Monitor test", func() {
		monitor := New(context.Background())
		Expect(monitor).ShouldNot(BeNil())

		go func() {
			time.Sleep(time.Second * 5)
			monitor.Stop()
		}()

		monitor.watch()
	})
})
