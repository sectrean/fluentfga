package fga_test

import (
	"context"
	"os"
	"testing"

	fga "github.com/johnrutherford/fluentfga/examples/iot"
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

func Test(t *testing.T) {
	ctx := context.Background()
	model := fga.NewAuthorizationModel(NewClient())

	anne := fga.User{UserID: "anne"}
	device := fga.Device{DeviceID: "1"}

	err := model.Device().
		SecurityGuard().
		Write(ctx, anne, device)
	assert.NoError(t, err)

	allowed, err := model.Device().
		LiveVideoViewer().
		Check(ctx, anne, device)
	assert.True(t, allowed)
	assert.NoError(t, err)
}
