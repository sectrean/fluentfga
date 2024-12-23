package fluentfga

// TODO: Conditions on contextual tuples

type CheckListOption interface {
	CheckOption
	ListObjectsOption
	ListUsersOption
}

func WithContextualTuples(tuples ...Tuple) CheckListOption {
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

func WithContext(context map[string]any) CheckListOption {
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
