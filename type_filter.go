package fluentfga

import (
	sdk "github.com/openfga/go-sdk"
)

type UserTypeFilter[U FilterType] struct{}

func (UserTypeFilter[U]) UserTypeFilter() sdk.UserTypeFilter {
	var u U

	filter := sdk.UserTypeFilter{
		Type: u.FgaType(),
	}

	if r, ok := any(u).(relation); ok {
		rel := r.Relation()
		filter.Relation = &rel
	}

	return filter
}
