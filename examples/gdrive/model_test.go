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
		model.DocumentCommenterRelation{}.Tuple(beth, doc),
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
		model.DocumentOwnerRelation{}.Tuple(anne, doc),
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
}

func Test_OrganizationPermissions(t *testing.T) {
	ctx := context.Background()
	client := NewClient()

	anne := model.User{ID: "anne"}
	beth := model.User{ID: "beth"}
	charles := model.User{ID: "charles"}
	domain := model.Domain{ID: "xyz"}
	doc := model.Document{ID: "2021-budget"}

	err := fluentfga.Write(
		model.DomainMemberRelation{}.Tuple(anne, domain),
	).Execute(ctx, client)
	assert.NoError(t, err)

	err = fluentfga.Write(
		model.DomainMemberRelation{}.Tuple(beth, domain),
	).Execute(ctx, client)
	assert.NoError(t, err)

	err = fluentfga.Write(
		model.DomainMemberRelation{}.Tuple(charles, domain),
	).Execute(ctx, client)
	assert.NoError(t, err)

	err = fluentfga.Write(
		model.DocumentViewerRelation{}.Tuple(model.DomainMemberUserset{domain}, doc),
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
