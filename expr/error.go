package expr

import (
	"errors"

	"github.com/petersalex27/yew-packages/nameable"
)

func redefineNameInTable[T nameable.Nameable](name Const[T]) error {
	return errors.New("tried to redefine " + name.String())
}

func nameNotDefined[T nameable.Nameable](name Const[T]) error {
	return errors.New(name.String() + " is not defined")
}

func redefineInv[T nameable.Nameable](name Const[T]) error {
	return errors.New("cannot redefine the inverse of " + name.String())
}
