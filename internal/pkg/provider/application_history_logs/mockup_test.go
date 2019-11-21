package application_history_logs

import "github.com/onsi/ginkgo"

var _ = ginkgo.Describe("Mockup Application Network provider", func() {
	provider := NewMockupApplicationHistoryLogsProvider()
	RunTest(provider)
})