package fluentfga

import (
	sdk "github.com/openfga/go-sdk"
)

// NewTuple creates a new Tuple with the given user, relation, and object.
//
// The user is not type-constrained, so it is up to the caller to ensure that the user is of a valid type
// for the relation.
//
// It's recommended to use [DirectRelation.NewTuple] to ensure type safety.
func NewTuple[U, O Object](user U, relation Relation[O], object O) TupleWithoutCondition {
	return tuple{
		user:     user.String(),
		relation: relation.Relation(),
		object:   object.String(),
	}
}

type tuple struct {
	user      string
	relation  string
	object    string
	condition *sdk.RelationshipCondition
}

func (t tuple) WithCondition(c Condition) Tuple {
	cond := c.SdkRelationshipCondition()
	t.condition = &cond

	return t
}

func (t tuple) SdkTupleKey() sdk.TupleKey {
	return sdk.TupleKey{
		User:      t.user,
		Relation:  t.relation,
		Object:    t.object,
		Condition: t.condition,
	}
}

func (t tuple) SdkTupleKeyWithoutCondition() sdk.TupleKeyWithoutCondition {
	return sdk.TupleKeyWithoutCondition{
		User:     t.user,
		Relation: t.relation,
		Object:   t.object,
	}
}
