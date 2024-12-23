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

type relation interface {
	Relation() string
}

type Relation[O Object] interface {
	FgaType() string
	Relation() string
	String() string

	ObjectType(O)
}

type DirectRelation[U, O Object] interface {
	Relation[O]

	UserType(U)
	Tuple(U, O) Tuple
}

type Userset interface {
	Object

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
