package model_test

import (
	"context"
	"os"
	"testing"

	"github.com/johnrutherford/fluentfga"
	model "github.com/johnrutherford/fluentfga/examples/iot"
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
	client := NewClient()

	anne := model.User{ID: "anne"}
	bob := model.User{ID: "bob"}
	device := model.Device{ID: "1"}
	securityGuard := model.DeviceSecurityGuardRelation{}
	liveVideoViewer := model.DeviceLiveVideoViewerRelation{}

	err := fluentfga.Write(
		securityGuard.NewTuple(anne, device),
	).Execute(ctx, client)
	assert.NoError(t, err)

	allowed, err := fluentfga.Check(
		anne,
		liveVideoViewer,
		device,
		fluentfga.WithContextualTuples(
			fluentfga.NewTuple(bob, securityGuard, device),
		),
		fluentfga.WithContext(map[string]any{}),
	).Execute(ctx, client)

	assert.True(t, allowed)
	assert.NoError(t, err)
}
