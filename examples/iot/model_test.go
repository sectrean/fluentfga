package model_test

import (
	"context"
	"testing"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/sectrean/fluentfga"
	model "github.com/sectrean/fluentfga/examples/iot"
	"github.com/sectrean/fluentfga/fgatest"
	authzmodel "github.com/sectrean/fluentfga/model"
	"github.com/stretchr/testify/suite"
)

type ModelSuite struct {
	fgatest.Suite
	authzModel *openfgav1.AuthorizationModel
}

func TestModelSuite(t *testing.T) {
	suite.Run(t, new(ModelSuite))
}

func (s *ModelSuite) SetupSuite() {
	s.Suite.SetupSuite()

	authzModel, err := authzmodel.ReadModelFromFile("model.fga")
	s.Require().NoError(err)

	s.authzModel = authzModel
}

func (s *ModelSuite) Test_ContextualTuples() {
	ctx := context.Background()
	client := s.NewStore(ctx, s.authzModel)

	anne := model.User{ID: "anne"}
	device := model.Device{
		ID: "1",

		// This will be passed in as a contextual tuple
		SecurityGuard: anne,
	}

	allowed, err := fluentfga.Check(
		anne,
		model.DeviceLiveVideoViewerRelation{},
		device,
	).Execute(ctx, client)

	s.True(allowed)
	s.NoError(err)
}
