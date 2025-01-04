package fluentfga

import (
	"context"

	sdkclient "github.com/openfga/go-sdk/client"
)

func Write(
	tuple Tuple,
	opts ...WriteOption,
) *WriteRequest {
	req := &WriteRequest{
		body: sdkclient.ClientWriteRequest{
			Writes: []sdkclient.ClientTupleKey{
				tuple.SdkTupleKey(),
			},
		},
	}

	for _, opt := range opts {
		opt.applyWriteOption(req)
	}

	return req
}

func Delete(
	tuple TupleWithoutCondition,
	opts ...WriteOption,
) *WriteRequest {
	req := &WriteRequest{
		body: sdkclient.ClientWriteRequest{
			Deletes: []sdkclient.ClientTupleKeyWithoutCondition{
				tuple.SdkTupleKeyWithoutCondition(),
			},
		},
	}

	for _, opt := range opts {
		opt.applyWriteOption(req)
	}

	return req
}

func WriteMany(
	writes []Tuple,
	deletes []TupleWithoutCondition,
	opts ...WriteOption,
) *WriteRequest {
	req := &WriteRequest{}

	for _, w := range writes {
		req.body.Writes = append(req.body.Writes, w.SdkTupleKey())
	}
	for _, d := range deletes {
		req.body.Deletes = append(req.body.Deletes, d.SdkTupleKeyWithoutCondition())
	}

	for _, opt := range opts {
		opt.applyWriteOption(req)
	}

	return req
}

type WriteOption interface {
	applyWriteOption(*WriteRequest)
}

type WriteRequest struct {
	body    sdkclient.ClientWriteRequest
	options sdkclient.ClientWriteOptions
}

func (r *WriteRequest) getBody() *sdkclient.ClientWriteRequest {
	return &r.body
}

func (r *WriteRequest) getOptions() *sdkclient.ClientWriteOptions {
	return &r.options
}

func (r *WriteRequest) Execute(ctx context.Context, c sdkclient.SdkClient) error {
	_, err := c.Write(ctx).
		Body(r.body).
		Options(r.options).
		Execute()

	return err
}

type writeOption func(*WriteRequest)

func (o writeOption) applyWriteOption(req *WriteRequest) {
	o(req)
}

func WithWrites(tuples ...Tuple) WriteOption {
	return writeOption(func(req *WriteRequest) {
		body := req.getBody()
		for _, t := range tuples {
			body.Writes = append(body.Writes, t.SdkTupleKey())
		}
	})
}

func WithDeletes(tuples ...TupleWithoutCondition) WriteOption {
	return writeOption(func(req *WriteRequest) {
		body := req.getBody()
		for _, t := range tuples {
			body.Deletes = append(body.Deletes, t.SdkTupleKeyWithoutCondition())
		}
	})
}
