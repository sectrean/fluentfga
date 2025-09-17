package fluentfga

import (
	"context"

	sdk "github.com/openfga/go-sdk"
	sdkclient "github.com/openfga/go-sdk/client"
)

func ListUsers[O Object, R Relation[O], U Filterable](
	object O,
	relation R,
	filter UserTypeFilter[U],
	opts ...ListUsersOption,
) *ListUsersRequest[U] {
	req := &ListUsersRequest[U]{
		provider: object.Provider(),
		body: sdkclient.ClientListUsersRequest{
			Object: sdk.FgaObject{
				Type: object.FgaType(),
				Id:   object.Identifier(),
			},
			Relation: relation.Relation(),
			UserFilters: []sdk.UserTypeFilter{
				filter.sdkUserTypeFilter(),
			},
		},
	}

	for _, opt := range opts {
		opt.applyListUsersOption(req)
	}

	return req
}

type ListUsersOption interface {
	applyListUsersOption(listUsersRequestInterface)
}

type listUsersRequestInterface interface {
	getBody() *sdkclient.ClientListUsersRequest
	getOptions() *sdkclient.ClientListUsersOptions
}

type ListUsersRequest[U Filterable] struct {
	provider ObjectProvider
	body     sdkclient.ClientListUsersRequest
	options  sdkclient.ClientListUsersOptions
}

func (r *ListUsersRequest[U]) getBody() *sdkclient.ClientListUsersRequest {
	return &r.body
}

func (r *ListUsersRequest[U]) getOptions() *sdkclient.ClientListUsersOptions {
	return &r.options
}

func (r *ListUsersRequest[U]) Execute(ctx context.Context, c sdkclient.SdkClient) ([]U, error) {
	res, err := c.ListUsers(ctx).
		Body(r.body).
		Options(r.options).
		Execute()

	if err != nil {
		return nil, err
	}

	return NewUsers[U](res.Users, r.provider)
}
