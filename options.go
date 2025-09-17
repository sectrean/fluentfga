package fluentfga

import (
	sdk "github.com/openfga/go-sdk"
	sdkclient "github.com/openfga/go-sdk/client"
)

type QueryOption interface {
	CheckOption
	ListObjectsOption
	ListUsersOption
}

func WithContextualTuples(tuples ...Tuple) QueryOption {
	return contextualTuplesOption{tuples}
}

func withContextualTuplesFromObjects(objects ...Object) QueryOption {
	var tuples []Tuple
	for _, obj := range objects {
		if ct, ok := any(obj).(HasContextualTuples); ok {
			tuples = append(tuples, ct.ContextualTuples()...)
		}
	}

	return contextualTuplesOption{tuples}
}

type contextualTuplesOption struct {
	Tuples []Tuple
}

func (o contextualTuplesOption) applyCheckOption(req *CheckRequest) {
	body := req.getBody()
	for _, t := range o.Tuples {
		body.ContextualTuples = append(body.ContextualTuples, t.SdkTupleKey())
	}
}

func (o contextualTuplesOption) applyListObjectsOption(req listObjectsRequestInterface) {
	body := req.getBody()
	for _, t := range o.Tuples {
		body.ContextualTuples = append(body.ContextualTuples, t.SdkTupleKey())
	}
}

func (o contextualTuplesOption) applyListUsersOption(req listUsersRequestInterface) {
	body := req.getBody()
	for _, t := range o.Tuples {
		body.ContextualTuples = append(body.ContextualTuples, t.SdkTupleKey())
	}
}

func WithContext(context map[string]any) QueryOption {
	return checkContextOption{context}
}

type checkContextOption struct {
	Context map[string]any
}

func (o checkContextOption) applyCheckOption(req *CheckRequest) {
	req.getBody().Context = &o.Context
}

func (o checkContextOption) applyListObjectsOption(req listObjectsRequestInterface) {
	req.getBody().Context = &o.Context
}

func (o checkContextOption) applyListUsersOption(req listUsersRequestInterface) {
	req.getBody().Context = &o.Context
}

func WithAuthorizationModelID(id string) AuthorizationModelIDOption {
	return authorizationModelIDOption{id}
}

type AuthorizationModelIDOption interface {
	QueryOption
	WriteOption
}

type authorizationModelIDOption struct {
	AuthorizationModelID string
}

func (o authorizationModelIDOption) applyWriteOption(req *WriteRequest) {
	options := req.getOptions()
	options.AuthorizationModelId = &o.AuthorizationModelID
}

func (o authorizationModelIDOption) applyCheckOption(req *CheckRequest) {
	options := req.getOptions()
	options.AuthorizationModelId = &o.AuthorizationModelID
}

func (o authorizationModelIDOption) applyListObjectsOption(req listObjectsRequestInterface) {
	options := req.getOptions()
	options.AuthorizationModelId = &o.AuthorizationModelID
}

func (o authorizationModelIDOption) applyListUsersOption(req listUsersRequestInterface) {
	options := req.getOptions()
	options.AuthorizationModelId = &o.AuthorizationModelID
}

func WithStoreID(id string) StoreIDOption {
	return storeIDOption{id}
}

type StoreIDOption interface {
	QueryOption
	WriteOption
}

type storeIDOption struct {
	StoreID string
}

func (o storeIDOption) applyWriteOption(req *WriteRequest) {
	options := req.getOptions()
	options.StoreId = &o.StoreID
}

func (o storeIDOption) applyCheckOption(req *CheckRequest) {
	options := req.getOptions()
	options.StoreId = &o.StoreID
}

func (o storeIDOption) applyListObjectsOption(req listObjectsRequestInterface) {
	options := req.getOptions()
	options.StoreId = &o.StoreID
}

func (o storeIDOption) applyListUsersOption(req listUsersRequestInterface) {
	options := req.getOptions()
	options.StoreId = &o.StoreID
}

func WithTransaction(opts sdkclient.TransactionOptions) WriteOption {
	return transactionOption{opts}
}

type transactionOption struct {
	Transaction sdkclient.TransactionOptions
}

func (t transactionOption) applyWriteOption(req *WriteRequest) {
	req.getOptions().Transaction = &t.Transaction
}

func WithConsistency(preference sdk.ConsistencyPreference) QueryOption {
	return consistencyOption{preference}
}

type consistencyOption struct {
	Consistency sdk.ConsistencyPreference
}

func (c consistencyOption) applyListObjectsOption(req listObjectsRequestInterface) {
	req.getOptions().Consistency = &c.Consistency
}

func (c consistencyOption) applyListUsersOption(req listUsersRequestInterface) {
	req.getOptions().Consistency = &c.Consistency
}

func (c consistencyOption) applyCheckOption(req *CheckRequest) {
	req.getOptions().Consistency = &c.Consistency
}
