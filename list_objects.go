package fluentfga

import (
	"context"

	sdkclient "github.com/openfga/go-sdk/client"
)

func ListObjects[U, O Object, R Relation[O]](
	user U,
	relation R,
	opts ...ListObjectsOption,
) *ListObjectsRequest[O] {
	req := &ListObjectsRequest[O]{
		provider: user.Provider(),
		body: sdkclient.ClientListObjectsRequest{
			User:     user.String(),
			Relation: relation.Relation(),
			Type:     relation.FgaType(),
		},
	}

	for _, opt := range opts {
		opt.applyListObjectsOption(req)
	}

	return req
}

type ListObjectsOption interface {
	applyListObjectsOption(listObjectsRequestInterface)
}

type listObjectsRequestInterface interface {
	getBody() *sdkclient.ClientListObjectsRequest
}

type ListObjectsRequest[O Object] struct {
	provider ObjectProvider
	body     sdkclient.ClientListObjectsRequest
	options  sdkclient.ClientListObjectsOptions
}

func (r *ListObjectsRequest[O]) getBody() *sdkclient.ClientListObjectsRequest {
	return &r.body
}

func (r *ListObjectsRequest[O]) Execute(ctx context.Context, c sdkclient.SdkClient) ([]O, error) {
	res, err := c.ListObjects(ctx).
		Body(r.body).
		Options(r.options).
		Execute()

	if err != nil {
		return nil, err
	}

	return ParseObjects[O](res.Objects, r.provider)
}
