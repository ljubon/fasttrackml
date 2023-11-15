//go:build integration

package namespace

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/G-Research/fasttrackml/pkg/api/mlflow/common"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/dao/models"
	"github.com/G-Research/fasttrackml/pkg/ui/admin/request"
	"github.com/G-Research/fasttrackml/tests/integration/golang/helpers"
)

type UpdateNamespaceTestSuite struct {
	helpers.BaseTestSuite
}

func TestUpdateNamespaceTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateNamespaceTestSuite))
}

func (s *UpdateNamespaceTestSuite) Test_Ok() {
	defer func() {
		require.Nil(s.T(), s.NamespaceFixtures.UnloadFixtures())
	}()
	_, err := s.NamespaceFixtures.CreateNamespace(context.Background(), &models.Namespace{
		ID:                  1,
		Code:                "default",
		DefaultExperimentID: common.GetPointer(int32(0)),
	})
	require.Nil(s.T(), err)
	ns, err := s.NamespaceFixtures.CreateNamespace(context.Background(), &models.Namespace{
		ID:                  2,
		Code:                "test2",
		Description:         "test namespace 2 description",
		DefaultExperimentID: common.GetPointer(int32(0)),
	})
	require.Nil(s.T(), err)

	request := request.Namespace{
		Code:        "test2Updated",
		Description: "test namespace 2 description updated",
	}
	require.Nil(
		s.T(),
		s.AdminClient.WithMethod(
			http.MethodPut,
		).WithRequest(
			request,
		).DoRequest("/namespaces/%d", ns.ID),
	)

	namespace, err := s.NamespaceFixtures.GetNamespaceByID(context.Background(), ns.ID)
	require.Nil(s.T(), err)

	assert.Equal(s.T(), namespace.Code, request.Code)
	assert.Equal(s.T(), namespace.Description, request.Description)
}

func (s *UpdateNamespaceTestSuite) Test_Error() {
	defer func() {
		require.Nil(s.T(), s.NamespaceFixtures.UnloadFixtures())
	}()
	_, err := s.NamespaceFixtures.CreateNamespace(context.Background(), &models.Namespace{
		ID:                  1,
		Code:                "default",
		DefaultExperimentID: common.GetPointer(int32(0)),
	})
	require.Nil(s.T(), err)
	_, err = s.NamespaceFixtures.CreateNamespace(context.Background(), &models.Namespace{
		ID:                  2,
		Code:                "test2",
		Description:         "test namespace 2 description",
		DefaultExperimentID: common.GetPointer(int32(0)),
	})
	require.Nil(s.T(), err)
	expectedNamespaces, err := s.NamespaceFixtures.GetNamespaces(context.Background())
	require.Nil(s.T(), err)

	testData := []struct {
		name     string
		ID       string
		request  *request.Namespace
		response map[string]any
	}{
		{
			name: "UpdateNamespaceWithNotFoundID",
			ID:   "10",
			request: &request.Namespace{
				Code:        "testUpdated",
				Description: "test namespace updated",
			},
			response: map[string]any{
				"message": "An unexepected error was encountered: namespace not found by id: 10",
				"status":  "error",
			},
		},
		{
			name: "UpdateNamespaceWithEmptyCode",
			ID:   "2",
			request: &request.Namespace{
				Code:        "",
				Description: "test namespace updated",
			},
			response: map[string]any{
				"message": "The namespace code is invalid.",
				"status":  "error",
			},
		},
		{
			name: "UpdateNamespaceWithDuplicatedCode",
			ID:   "2",
			request: &request.Namespace{
				Code:        "default",
				Description: "test namespace updated",
			},
			response: map[string]any{
				"message": "The namespace code is already in use.",
				"status":  "error",
			},
		},
	}
	for _, tt := range testData {
		s.T().Run(tt.name, func(t *testing.T) {
			var resp any
			require.Nil(
				s.T(),
				s.AdminClient.WithMethod(
					http.MethodPut,
				).WithRequest(
					tt.request,
				).WithResponse(
					&resp,
				).DoRequest(
					"/namespaces/%s", tt.ID,
				),
			)
			assert.Equal(s.T(), resp, tt.response)
		})
		actualNamespaces, err := s.NamespaceFixtures.GetNamespaces(context.Background())
		require.Nil(s.T(), err)
		assert.Equal(s.T(), expectedNamespaces, actualNamespaces)
	}
}