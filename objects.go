package fluentfga

import (
	"errors"
	"fmt"

	sdk "github.com/openfga/go-sdk"
)

// NewTuple creates a new Tuple with the given user, relation, and object.
//
// The user is not type-constrained, so it is up to the caller to ensure that the user is of a valid type
// for the relation.
func NewTuple[U Object, R Relation[O], O Object](user U, relation R, object O) TupleWithoutCondition {
	return tuple{
		u: user.String(),
		r: relation.Relation(),
		o: object.String(),
	}
}

type tuple struct {
	u string
	r string
	o string
	c *sdk.RelationshipCondition
}

func (t tuple) WithCondition(c Condition) Tuple {
	cond := c.SdkRelationshipCondition()
	t.c = &cond

	return t
}

func (t tuple) SdkTupleKey() sdk.TupleKey {
	return sdk.TupleKey{
		User:      t.u,
		Relation:  t.r,
		Object:    t.o,
		Condition: t.c,
	}
}

func (t tuple) SdkTupleKeyWithoutCondition() sdk.TupleKeyWithoutCondition {
	return sdk.TupleKeyWithoutCondition{
		User:     t.u,
		Relation: t.r,
		Object:   t.o,
	}
}

func ParseObjects[O Object](objects []string, p ObjectProvider) ([]O, error) {
	result := make([]O, 0, len(objects))

	var errs []error
	for _, objStr := range objects {
		var typ, id string
		n, err := fmt.Sscanf(objStr, "%s:%s", &typ, &id)
		if err != nil {
			err := fmt.Errorf("invalid object %q: %w", objStr, err)
			errs = append(errs, err)
			continue
		}

		if n < 2 {
			err := fmt.Errorf("invalid object %q: missing type or id", objStr)
			errs = append(errs, err)
			continue
		}

		obj, err := p.NewObject(typ, id, "")
		if err != nil {
			errs = append(errs, err)
			continue
		}

		o, ok := obj.(O)
		if !ok {
			err := fmt.Errorf("unexpected object type %T", obj)
			errs = append(errs, err)
			continue
		}

		result = append(result, o)
	}

	if err := errors.Join(errs...); err != nil {
		return result, fmt.Errorf("parse openfga objects: %w", err)
	}

	return result, nil
}

func NewUsers[U Object](users []sdk.User, reg ObjectProvider) ([]U, error) {
	result := make([]U, 0, len(users))

	var errs []error
	for _, u := range users {
		var obj Object
		var err error

		switch {
		case u.Object != nil:
			obj, err = reg.NewObject(u.Object.Type, u.Object.Id, "")

		case u.Wildcard != nil:
			obj, err = reg.NewObject(u.Wildcard.Type, "*", "")

		case u.Userset != nil:
			obj, err = reg.NewObject(u.Userset.Type, u.Userset.Id, u.Userset.Relation)

		default:
			err = fmt.Errorf("unknown user type %v", u)
		}

		if err != nil {
			errs = append(errs, err)
			continue
		}

		user, ok := obj.(U)
		if !ok {
			err := fmt.Errorf("unexpected object type %T", obj)
			errs = append(errs, err)
			continue
		}

		result = append(result, user)
	}

	return result, errors.Join(errs...)
}
