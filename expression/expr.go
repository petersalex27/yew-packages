package expr

import (
	"github.com/petersalex27/yew-packages/equality"
	str "github.com/petersalex27/yew-packages/stringable"
)

type Expression interface {
	str.Stringable
	equality.Eq[Expression]
	StrictString() string
	StrictEquals(Expression) bool
	Replace(Variable, Expression) (Expression, bool)
	UpdateVars(gt int, by int) Expression
	Again() (Expression, bool)
	Bind(BindersOnly) Expression
	Find(Variable) bool
	PrepareAsRHS() Expression
	Rebind() Expression
	Copy() Expression
	ForceRequest() Expression
}

type ApplicableExpression interface {
	DoApplication(e Expression) Expression
	AgainApply(e Expression) (res Expression, again bool)
}

type SplitableExpression interface {
	Expression
	Split() (left, right Expression)
}

type InvariableExpression interface {
	str.Stringable
	equality.Eq[Expression]
	DeepCopy() InvariableExpression
}

var binder_string string = "λ"
var to_string string = " . "
var apply_string string = " "
var list_l_enclose string = "["
var list_r_enclose string = "]"
var list_split string = ", "
var strPtrs = []*string{
	&binder_string,
	&to_string,
	&apply_string,
	&list_l_enclose,
	&list_r_enclose,
	&list_split,
}

func GetBinderString() string { return binder_string }

func GetToString() string { return to_string }

func GetApplyString() string { return apply_string }

func GetListEnclose() (string, string) {
	return list_l_enclose, list_r_enclose
}

func doInit(strs ...string) {
	for i, s := range strs {
		if len(s) != 0 {
			(*strPtrs[i]) = s
		}
	}
}

func Init(binder string, to string, apply string, listEnclose [2]string, listSplit string) {
	doInit(binder, to, apply, listEnclose[0], listEnclose[1], listSplit)
}
