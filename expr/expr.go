package expr

import (
	"github.com/petersalex27/yew-packages/equality"
	"github.com/petersalex27/yew-packages/nameable"
	str "github.com/petersalex27/yew-packages/stringable"
)

type Expression[T nameable.Nameable] interface {
	str.Stringable
	//equality.Eq[Expression[T]]
	Equals(*Context[T], Expression[T]) bool
	StrictString() string
	StrictEquals(Expression[T]) bool
	Replace(Variable[T], Expression[T]) (Expression[T], bool)
	UpdateVars(gt int, by int) Expression[T]
	Again() (Expression[T], bool)
	Bind(BindersOnly[T]) Expression[T]
	Find(Variable[T]) bool
	PrepareAsRHS() Expression[T]
	Rebind() Expression[T]
	Copy() Expression[T]
	ForceRequest() Expression[T]
	Collect() []T
}

type ApplicableExpression[T nameable.Nameable] interface {
	DoApplication(e Expression[T]) Expression[T]
	AgainApply(e Expression[T]) (res Expression[T], again bool)
}

type SplitableExpression[T nameable.Nameable] interface {
	Expression[T]
	Split() (left, right Expression[T])
}

type InvariableExpression[T nameable.Nameable] interface {
	str.Stringable
	equality.Eq[Expression[T]]
	DeepCopy() InvariableExpression[T]
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
