package polyfea

import (
	_ "embed"
	"testing"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MiddlewaresTestSuite struct {
	suite.Suite
	mfcRepository repository.Repository[*v1alpha1.MicroFrontendClass]
}

func TestMiddlewaresTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewaresTestSuite))
}

func (suite *MiddlewaresTestSuite) SetupTest() {
	suite.SetupMfcRepository()
}

// SetupMfcRepository initializes the in-memory repository with test data.
func (suite *MiddlewaresTestSuite) SetupMfcRepository() {
	testData := []struct {
		name      string
		namespace string
		baseUri   string
	}{
		{"default", "default", "/"},
		{"fea", "feas", "/fea/"},
		{"feature", "features", "/feature/"},
	}

	suite.mfcRepository = repository.NewInMemoryRepository[*v1alpha1.MicroFrontendClass]()

	for _, data := range testData {
		suite.mfcRepository.Store(&v1alpha1.MicroFrontendClass{
			ObjectMeta: metav1.ObjectMeta{
				Name:      data.name,
				Namespace: data.namespace,
			},
			Spec: v1alpha1.MicroFrontendClassSpec{
				BaseUri: &data.baseUri,
			},
		})
	}
}

// TestGetMicrofrontendAndBase tests the retrieval of microfrontend class and base URI.
func (suite *MiddlewaresTestSuite) TestGetMicrofrontendAndBase() {
	testParams := []struct {
		requestPath     string
		expectBasePath  string
		expectClassName string
	}{
		// Test case: Nonexistent path should default to "default" class
		{requestPath: "/nonexistent", expectBasePath: "/", expectClassName: "default"},
		// Test case: Path matching "fea" should return "fea" class
		{requestPath: "/fea/asset", expectBasePath: "/fea/", expectClassName: "fea"},
		{requestPath: "/fea", expectBasePath: "/fea/", expectClassName: "fea"},
		// Test case: Path matching "feature" should return "feature" class
		{requestPath: "/feature", expectBasePath: "/feature/", expectClassName: "feature"},
		{requestPath: "/feature/asset", expectBasePath: "/feature/", expectClassName: "feature"},
		// Test case: Path not matching any class should default to "default" class
		{requestPath: "/fea-nix", expectBasePath: "/", expectClassName: "default"},
		{requestPath: "/other/qweqwesop", expectBasePath: "/", expectClassName: "default"},
	}

	for _, params := range testParams {
		suite.Run(params.requestPath, func() {
			// Arrange
			// Act
			basePath, microfrontend, err := getMicrofrontendClassAndBase(params.requestPath, suite.mfcRepository)

			// Assert
			suite.Nil(err, "Expected no error")
			suite.Equal(params.expectBasePath, basePath, "Expected base path to match")
			suite.Equal(params.expectClassName, microfrontend.Name, "Expected class name to match")
		})
	}
}

// Added tests for edge cases and error handling
func (suite *MiddlewaresTestSuite) TestGetMicrofrontendAndBaseEdgeCases() {
	testParams := []struct {
		requestPath     string
		expectBasePath  string
		expectClassName string
		expectError     bool
	}{
		// Test case: Empty requestPath
		{requestPath: "", expectBasePath: "/", expectClassName: "default", expectError: false},
		// Test case: Special characters in requestPath
		{requestPath: "/!@#$%^&*()", expectBasePath: "/", expectClassName: "default", expectError: false},
		// Test case: Path longer than any base URI
		{requestPath: "/this/path/is/way/too/long/to/match/anything", expectBasePath: "/", expectClassName: "default", expectError: false},
		// Test case: Simulate repository error
		{requestPath: "/feature", expectBasePath: "", expectClassName: "", expectError: true},
	}

	for _, params := range testParams {
		suite.Run(params.requestPath, func() {
			// Arrange
			if params.expectError {
				suite.mfcRepository = nil // Simulate repository error
			}

			// Act
			basePath, microfrontend, err := getMicrofrontendClassAndBase(params.requestPath, suite.mfcRepository)

			// Assert
			if params.expectError {
				suite.NotNil(err, "Expected an error")
			} else {
				suite.Nil(err, "Expected no error")
				suite.Equal(params.expectBasePath, basePath, "Expected base path to match")
				suite.Equal(params.expectClassName, microfrontend.Name, "Expected class name to match")
			}
		})
	}
}
