package types

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/str"
)

type ProductTyped interface {
	str.Stringable
	IsKindParameterized() bool
}

type DependentProductType struct {
	iterators []TypeJudgement[expr.Variable]
	Monotyped
}

func (d DependentProductType) IsKindParameterized() bool { return true  }

func (d DependentProductType) String() string {
	return "mapall " + str.Join(d.iterators, str.String(" ")) + " . " + d.Monotyped.String()
}