package fluentfga

import (
	"context"

	sdkclient "github.com/openfga/go-sdk/client"
)

func Check[U Object, R Relation[O], O Object](
	user U,
	relation R,
	object O,
	opts ...CheckOption,
) *CheckRequest {
	req := &CheckRequest{
		body: sdkclient.ClientCheckRequest{
			User:     user.String(),
			Relation: relation.Relation(),
			Object:   object.String(),
		},
	}

	defaultOpts := []CheckOption{
		withContextualTuplesFromObjects(user, object),
	}
	opts = append(defaultOpts, opts...)

	for _, opt := range opts {
		opt.applyCheckOption(req)
	}

	return req
}

type CheckOption interface {
	applyCheckOption(*CheckRequest)
}

type CheckRequest struct {
	body    sdkclient.ClientCheckRequest
	options sdkclient.ClientCheckOptions
}

func (r *CheckRequest) getBody() *sdkclient.ClientCheckRequest {
	return &r.body
}

func (r *CheckRequest) getOptions() *sdkclient.ClientCheckOptions {
	return &r.options
}

func (r *CheckRequest) Execute(ctx context.Context, c sdkclient.SdkClient) (bool, error) {
	res, err := c.Check(ctx).
		Body(r.body).
		Options(r.options).
		Execute()

	return res.GetAllowed(), err
}

// TODO: Implement BatchCheck
