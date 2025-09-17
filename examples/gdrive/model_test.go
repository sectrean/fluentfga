package model_test

import (
	"context"
	"testing"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/sectrean/fluentfga"
	model "github.com/sectrean/fluentfga/examples/gdrive"
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

func (s *ModelSuite) Test_IndividualPermissions() {
	ctx := context.Background()
	client := s.NewStore(ctx, s.T().Name(), s.authzModel)

	beth := model.User{ID: "beth"}
	anne := model.User{ID: "anne"}
	doc := model.Document{ID: "2021-budget"}

	err := fluentfga.Write(
		model.DocumentCommenterRelation{}.NewTuple(beth, doc),
	).Execute(ctx, client)

	s.NoError(err)

	allowed, err := fluentfga.Check(
		beth,
		model.DocumentCommenterRelation{},
		doc,
	).Execute(ctx, client)

	s.True(allowed)
	s.NoError(err)

	err = fluentfga.Write(
		model.DocumentOwnerRelation{}.NewTuple(anne, doc),
	).Execute(ctx, client)

	s.NoError(err)

	allowed, err = fluentfga.Check(
		anne,
		model.DocumentOwnerRelation{},
		doc,
	).Execute(ctx, client)

	s.True(allowed)
	s.NoError(err)

	allowed, err = fluentfga.Check(
		anne,
		model.DocumentWriterRelation{},
		doc,
	).Execute(ctx, client)

	s.True(allowed)
	s.NoError(err)

	users, err := fluentfga.ListUsers(
		doc,
		model.DocumentOwnerRelation{},
		fluentfga.UserTypeFilter[model.User]{},
	).Execute(ctx, client)

	// TODO: Check this assertion
	s.ElementsMatch([]model.User{anne}, users)
	s.NoError(err)
}

func (s *ModelSuite) Test_OrganizationPermissions() {
	ctx := context.Background()
	client := s.NewStore(ctx, s.T().Name(), s.authzModel)

	anne := model.User{ID: "anne"}
	beth := model.User{ID: "beth"}
	charles := model.User{ID: "charles"}
	domainMember := model.DomainMemberRelation{}
	domain := model.Domain{ID: "xyz"}
	doc := model.Document{ID: "2021-budget"}
	documentViewer := model.DocumentViewerRelation{}

	err := fluentfga.Write(
		fluentfga.NewTuple(anne, domainMember, domain),
	).Execute(ctx, client)
	s.NoError(err)

	err = fluentfga.Write(
		fluentfga.NewTuple(beth, domainMember, domain),
	).Execute(ctx, client)
	s.NoError(err)

	err = fluentfga.Write(
		domainMember.NewTuple(charles, domain),
	).Execute(ctx, client)
	s.NoError(err)

	err = fluentfga.Write(
		documentViewer.NewTuple(model.DomainMemberUserset{domain}, doc),
		// Alternate way to create the tuple:
		//
		// fluentfga.NewTuple(
		// 	model.DomainMemberUserset{domain},
		// 	documentViewer,
		// 	doc,
		// ),
	).Execute(ctx, client)
	s.NoError(err)

	allowed, err := fluentfga.Check(
		charles,
		model.DocumentViewerRelation{},
		doc,
	).Execute(ctx, client)

	s.True(allowed)
	s.NoError(err)
}

func (s *ModelSuite) Test_ContextualTuples() {
	ctx := context.Background()
	client := s.NewStore(ctx, s.T().Name(), s.authzModel)

	doc := model.Document{
		ID: "1234",
		Parent: &model.Document{
			ID: "5678",
		},
	}

	fluentfga.Check(
		model.User{ID: "john"},
		model.DocumentViewerRelation{},
		doc,
	).Execute(ctx, client)
}
