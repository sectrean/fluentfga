package fluentfga

import (
	sdk "github.com/openfga/go-sdk"
)

type Object interface {
	FgaType() string
	Identifier() string
	String() string

	Provider() ObjectProvider
}

type ObjectProvider interface {
	NewObject(typ, id, rel string) (Object, error)
}

type FilterType interface {
	Object

	IsFilterable()
}

type Relation[O Object] interface {
	FgaType() string
	Relation() string
	String() string

	// ObjectType does nothing.
	//
	// Deprecated: Don't use. It exists only to enforce type constraints.
	ObjectType(O)
}

type DirectRelation[O, U Object] interface {
	Relation[O]

	NewTuple(U, O) Tuple
}

type Userset interface {
	Object

	Relation() string
	IsUserset()
}

type Wildcard interface {
	Object

	IsWildcard()
}

type Condition interface {
	Name() string
	SdkRelationshipCondition() sdk.RelationshipCondition
}

type Tuple interface {
	SdkTupleKey() sdk.TupleKey
}

type TupleWithoutCondition interface {
	Tuple

	WithCondition(Condition) Tuple
	SdkTupleKeyWithoutCondition() sdk.TupleKeyWithoutCondition
}
