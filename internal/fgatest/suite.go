package fgatest

import (
	"context"
	"encoding/json"

	proto "github.com/openfga/api/proto/openfga/v1"
	sdkclient "github.com/openfga/go-sdk/client"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	ContainerName  = "fluentfga-test"
	ContainerImage = "openfga/openfga:latest"
)

type Suite struct {
	suite.Suite

	container testcontainers.Container
	apiUrl    string
}

func (s *Suite) SetupSuite() {
	ctx := context.Background()

	var err error
	s.container, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         ContainerName,
			Image:        ContainerImage,
			ExposedPorts: []string{"8080/tcp"},
		},
		Started: true,
		Reuse:   true,
	})
	s.Require().NoError(err)

	s.apiUrl, err = s.container.PortEndpoint(ctx, "8080", "http")
	s.Require().NoError(err)
}

// NewStore creates a new openfga store with the given model and returns an API client.
// It sets the store ID and authorization model ID on the client for further operations.
func (s *Suite) NewStore(ctx context.Context, name string, model *proto.AuthorizationModel) sdkclient.SdkClient {
	modelRequest, err := modelToClientRequest(model)
	s.Require().NoError(err)

	client, err := sdkclient.NewSdkClient(&sdkclient.ClientConfiguration{
		ApiUrl: s.apiUrl,
	})
	s.Require().NoError(err)

	store, err := client.CreateStore(ctx).
		Body(sdkclient.ClientCreateStoreRequest{
			Name: name,
		}).
		Execute()
	s.Require().NoError(err)

	err = client.SetStoreId(store.Id)
	s.Require().NoError(err)

	authzModel, err := client.WriteAuthorizationModel(ctx).
		Body(modelRequest).
		Execute()
	s.Require().NoError(err)

	err = client.SetAuthorizationModelId(authzModel.AuthorizationModelId)
	s.Require().NoError(err)

	return client
}

func (s *Suite) TearDownSuite() {
	ctx := context.Background()

	err := s.container.Terminate(ctx)
	s.Require().NoError(err)
}

func modelToClientRequest(model *proto.AuthorizationModel) (sdkclient.ClientWriteAuthorizationModelRequest, error) {
	jsonBytes, err := protojson.Marshal(model)
	if err != nil {
		return sdkclient.ClientWriteAuthorizationModelRequest{}, err
	}

	var request sdkclient.ClientWriteAuthorizationModelRequest
	err = json.Unmarshal(jsonBytes, &request)
	if err != nil {
		return sdkclient.ClientWriteAuthorizationModelRequest{}, err
	}

	return request, nil
}
