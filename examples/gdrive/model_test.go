package model_test

import (
	"context"
	"os"
	"testing"

	"github.com/johnrutherford/fluentfga"
	model "github.com/johnrutherford/fluentfga/examples/gdrive"
	sdkclient "github.com/openfga/go-sdk/client"
	"github.com/stretchr/testify/assert"
)

func NewClient() sdkclient.SdkClient {
	client, err := sdkclient.NewSdkClient(&sdkclient.ClientConfiguration{
		ApiUrl:               os.Getenv("FGA_API_URL"),
		StoreId:              os.Getenv("FGA_STORE_ID"),
		AuthorizationModelId: os.Getenv("FGA_MODEL_ID"),
	})
	if err != nil {
		panic(err)
	}

	return client
}

func Test_IndividualPermissions(t *testing.T) {
	ctx := context.Background()
	client := NewClient()

	beth := model.User{ID: "beth"}
	anne := model.User{ID: "anne"}
	doc := model.Document{ID: "2021-budget"}

	err := fluentfga.Write(
		model.DocumentCommenterRelation{}.NewTuple(beth, doc),
	).Execute(ctx, client)

	assert.NoError(t, err)

	allowed, err := fluentfga.Check(
		beth,
		model.DocumentCommenterRelation{},
		doc,
	).Execute(ctx, client)

	assert.True(t, allowed)
	assert.NoError(t, err)

	err = fluentfga.Write(
		model.DocumentOwnerRelation{}.NewTuple(anne, doc),
	).Execute(ctx, client)

	assert.NoError(t, err)

	allowed, err = fluentfga.Check(
		anne,
		model.DocumentOwnerRelation{},
		doc,
	).Execute(ctx, client)

	assert.True(t, allowed)
	assert.NoError(t, err)

	allowed, err = fluentfga.Check(
		anne,
		model.DocumentWriterRelation{},
		doc,
	).Execute(ctx, client)

	assert.True(t, allowed)
	assert.NoError(t, err)

	users, err := fluentfga.ListUsers(
		doc,
		model.DocumentOwnerRelation{},
		fluentfga.UserTypeFilter[model.User]{},
	).Execute(ctx, client)

	// TODO: Check this assertion
	assert.ElementsMatch(t, []model.User{anne}, users)
	assert.NoError(t, err)
}

func Test_OrganizationPermissions(t *testing.T) {
	ctx := context.Background()
	client := NewClient()

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
	assert.NoError(t, err)

	err = fluentfga.Write(
		fluentfga.NewTuple(beth, domainMember, domain),
	).Execute(ctx, client)
	assert.NoError(t, err)

	err = fluentfga.Write(
		domainMember.NewTuple(charles, domain),
	).Execute(ctx, client)
	assert.NoError(t, err)

	err = fluentfga.Write(
		// fluentfga.NewTuple(
		// 	model.DomainMemberUserset{domain},
		// 	documentViewer,
		// 	doc,
		// ),
		documentViewer.NewTuple(model.DomainMemberUserset{domain}, doc),
	).Execute(ctx, client)
	assert.NoError(t, err)

	allowed, err := fluentfga.Check(
		charles,
		model.DocumentViewerRelation{},
		doc,
	).Execute(ctx, client)
	assert.True(t, allowed)
	assert.NoError(t, err)
}
