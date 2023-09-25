package expr

import "errors"

func redefineNameInTable(name Const) error {
	return errors.New("tried to redefine " + string(name))
}

func nameNotDefined(name Const) error {
	return errors.New(string(name) + " is not defined")
}

func redefineInv(name Const) error {
	return errors.New("cannot redefine the inverse of " + string(name))
}
