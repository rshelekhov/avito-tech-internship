package api_tests

import (
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
)

func newTestAPI(t *testing.T) *httpexpect.Expect {
	apiHost := os.Getenv("TEST_API_HOST")
	if apiHost == "" {
		apiHost = "http://0.0.0.0:8080"
	}

	return httpexpect.Default(t, apiHost)
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, true, 10)
}
