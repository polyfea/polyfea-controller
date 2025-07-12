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
func (suite *MiddlewaresTestSuite) SetupMfcRepository() {
	suite.mfcRepository = repository.NewInMemoryRepository[*v1alpha1.MicroFrontendClass]()
	defaultBase := "/"
	suite.mfcRepository.StoreItem(&v1alpha1.MicroFrontendClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "default",
		},
		Spec: v1alpha1.MicroFrontendClassSpec{
			BaseUri: &defaultBase,
		},
	})

	feaBase := "/fea/"
	suite.mfcRepository.StoreItem(&v1alpha1.MicroFrontendClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fea",
			Namespace: "feas",
		},
		Spec: v1alpha1.MicroFrontendClassSpec{
			BaseUri: &feaBase,
		},
	})

	featureBase := "/feature/"
	suite.mfcRepository.StoreItem(&v1alpha1.MicroFrontendClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "feature",
			Namespace: "features",
		},
		Spec: v1alpha1.MicroFrontendClassSpec{
			BaseUri: &featureBase,
		},
	})
}

func (suite *MiddlewaresTestSuite) TestGetMicrofrontendAndBase() {
	testParams := []struct {
		requestPath     string
		expectBasePath  string
		expectClassName string
	}{
		{"/nonexistent", "/", "default"},
		{"/fea/asset", "/fea/", "fea"},
		{"/fea", "/fea/", "fea"},
		{"/feature", "/feature/", "feature"},
		{"/feature/asset", "/feature/", "feature"},
		{"/fea-nix", "/", "default"},
		{"/other/qweqwesop", "/", "default"},
	}

	for _, params := range testParams {
		suite.Run(params.requestPath, func() {
			// Arrange
			// Act
			basePath, microfrontend, err := getMicrofrontendClassAndBase(params.requestPath, suite.mfcRepository)

			// Assert
			suite.Nil(err)
			suite.Equal(params.expectBasePath, basePath)
			suite.Equal(params.expectClassName, microfrontend.Name)
		})
	}

}
