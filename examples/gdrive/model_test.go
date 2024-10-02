package fga_test

import (
	"context"
	"os"
	"testing"

	fga "github.com/johnrutherford/fluentfga/examples/gdrive"
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
	model := fga.NewAuthorizationModel(NewClient())

	beth := fga.User{UserID: "beth"}
	anne := fga.User{UserID: "anne"}
	doc := fga.Document{DocumentID: "2021-budget"}

	err := model.Document().
		Commenter().
		Write(ctx, beth, doc)
	assert.NoError(t, err)

	allowed, err := model.Document().
		Commenter().
		Check(ctx, beth, doc)
	assert.True(t, allowed)
	assert.NoError(t, err)

	err = model.Document().
		Owner().
		Write(ctx, anne, doc)
	assert.NoError(t, err)

	allowed, err = model.Document().
		Owner().
		Check(ctx, anne, doc)
	assert.True(t, allowed)
	assert.NoError(t, err)

	allowed, err = model.Document().
		Writer().
		Check(ctx, anne, doc)
	assert.True(t, allowed)
	assert.NoError(t, err)
}

func Test_OrganizationPermissions(t *testing.T) {
	ctx := context.Background()
	model := fga.NewAuthorizationModel(NewClient())

	anne := fga.User{UserID: "anne"}
	beth := fga.User{UserID: "beth"}
	charles := fga.User{UserID: "charles"}
	domain := fga.Domain{DomainID: "xyz"}
	doc := fga.Document{DocumentID: "2021-budget"}

	err := model.Domain().
		Member().
		Write(ctx, anne, domain)
	assert.NoError(t, err)

	err = model.Domain().
		Member().
		Write(ctx, beth, domain)
	assert.NoError(t, err)

	err = model.Domain().
		Member().
		Write(ctx, charles, domain)
	assert.NoError(t, err)

	err = model.Document().
		Viewer().
		Write(ctx, domain.MemberUserset(), doc)
	assert.NoError(t, err)

	allowed, err := model.Document().
		Viewer().
		Check(ctx, charles, doc)
	assert.True(t, allowed)
	assert.NoError(t, err)
}
