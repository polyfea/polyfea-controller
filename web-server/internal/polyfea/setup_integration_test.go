package polyfea

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const TestServerUrlName = "TEST_SERVER_URL"

type Test struct {
	Name string
	Func func(*testing.T)
}

type IntegrationTestSuite struct {
	TestRouter http.Handler
	TestSet    []Test
}

func TestPolyfeaIntegration(t *testing.T) {
	integrationTestSuites := collectIntegrationTestSuites()

	for _, integrationTestSuite := range integrationTestSuites {
		// Test server
		var testServer = httptest.NewServer(integrationTestSuite.TestRouter)
		defer testServer.Close()

		t.Setenv(TestServerUrlName, testServer.URL)
		for _, test := range integrationTestSuite.TestSet {
			t.Run(test.Name, test.Func)
		}
	}
}

func collectIntegrationTestSuites() []IntegrationTestSuite {
	return []IntegrationTestSuite{
		apiPolyfeaTestSuite,
		middlewaresTestSuite,
		proxyTestSuite,
	}
}
