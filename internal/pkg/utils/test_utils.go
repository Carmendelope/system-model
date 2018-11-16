package utils

import "os"

func RunIntegrationTests() bool {
	var runIntegration = os.Getenv("RUN_INTEGRATION_TEST")
	return runIntegration == "true"
}