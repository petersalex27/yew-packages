package types

import (
	"errors"
	"strconv"

	"github.com/petersalex27/yew-packages/nameable"
)

func typeMismatch[T nameable.Nameable](a, b Type[T]) error { 
	return errors.New(
		"the type " + a.String() + " is not in the same equivalence class as " + b.String())
}

func alreadyDefined(name string) error {
	return errors.New("the type " + name + " has already been defined")
}

func ruleAlreadyDefined(name string) error {
	return errors.New("the rule " + name + " has already been defined")
}

func ruleIdDNE(id RuleID) error {
	return errors.New("no rule with id #" + strconv.FormatInt(int64(id), 10))
}