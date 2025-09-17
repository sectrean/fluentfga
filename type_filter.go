package fluentfga

import (
	sdk "github.com/openfga/go-sdk"
)

type UserTypeFilter[U Filterable] struct{}

func (UserTypeFilter[U]) sdkUserTypeFilter() sdk.UserTypeFilter {
	var u U

	filter := sdk.UserTypeFilter{
		Type: u.FgaType(),
	}

	if us, ok := any(u).(Userset); ok {
		rel := us.Relation()
		filter.Relation = &rel
	}

	return filter
}
