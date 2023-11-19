package expr

import (
	"github.com/petersalex27/yew-packages/equality"
	"github.com/petersalex27/yew-packages/nameable"
	str "github.com/petersalex27/yew-packages/stringable"
)

type Bindable[T nameable.Nameable] interface {
	Rebind() Expression[T]
	Bind(BindersOnly[T]) Expression[T]
}

type Replaceable[T nameable.Nameable] interface {
	Replace(Variable[T], Expression[T]) (Expression[T], bool)
	UpdateVars(gt int, by int) Expression[T]
	BodyAbstract(v Variable[T], name Const[T]) Expression[T]
}

type Referable[T nameable.Nameable] interface {
	Expression[T]
	GetReferred() T
}

type Expression[T nameable.Nameable] interface {
	str.Stringable
	Bindable[T]
	nameable.Collectable[T]
	Replaceable[T]
	Equals(*Context[T], Expression[T]) bool
	StrictString() string
	StrictEquals(Expression[T]) bool
	Again() (Expression[T], bool)
	Find(Variable[T]) bool
	PrepareAsRHS() Expression[T]
	Copy() Expression[T]
	ForceRequest() Expression[T]
	ExtractVariables(gt int) []Variable[T]
	Flatten() []Expression[T]
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
var hidden_binder_string string = "Λ"
var to_string string = " . "
var apply_string string = " "
var list_open string = "["
var list_close string = "]"
var list_sep string = ", "
var case_sep string = " | "
var match_head string = "match "
var case_onMatch string = " -> "
var rec_head string = "rec "
var rec_sep string = " and "
var let_head string = "let "
var contextualized_sep string = " in "
var where_infix string = " where "
var grouping_open string = "("
var grouping_close string = ")"

func GetBinderString() string { return binder_string }

func GetToString() string { return to_string }

func GetApplyString() string { return apply_string }

func encloseListString(list string) string {
	return list_open + list + list_close
}

func listSepString() string { return list_sep }

func applyString(left, right string) string {
	return left + apply_string + right
}

func recHeadString() string { return rec_head }

func recString(defs, contextualized string) string {
	return contextualizeString(recHeadString() + defs, contextualized)
}

func recSepString() string { return rec_sep }

func bindingString(binders, bound string) string {
	return binder_string + binders + to_string + bound
} 

func hiddenBindingString(binders, bound string) string {
	return hidden_binder_string + binders + to_string + bound
}

func contextualizeString(left, right string) string {
	return left + contextualized_sep + right
}

func caseSepString() string {
	return case_sep
}

func onMatchString() string {
	return case_onMatch
}

func matchHeadString(matching, cases string) string {
	return contextualizeString(match_head + matching, cases)
}

func letString(def, contextualized string) string {
	return contextualizeString(let_head + def, contextualized)
}

func fixWhere(left, right string) string {
	return left + where_infix + right
}

func groupStringed(stringed string) string {
	return grouping_open + stringed + grouping_close
}

func GenHiddenBinder(hiddenBinder string) func() {
	return func() {
		hidden_binder_string = hiddenBinder
	}
}

func GenSetWhere(where string) func() {
	return func() {
		where_infix = where
	}
}

func GenSetFunc(binder, to string) func() {
	return func() {
		binder_string, to_string = binder, to
	}
}

func GenSetApply(apply string) func() {
	return func() {
		apply_string = apply
	}
}

func GenSetList(open, close, sep string) func() {
	return func() {
		list_open, list_close, list_sep = open, close, sep
	}
}

func GenSetGrouping(open, close string) func() {
	return func() {
		grouping_open, grouping_close = open, close
	}
}

func GenSetMatch(head, sepCases, sepCaseHeadCaseTail string) func() {
	return func() {
		match_head, case_sep, case_onMatch = head, sepCases, sepCaseHeadCaseTail
	}
}

// creates think for init. context sep. related strings
func GenSetContextSep(sep string) func() {
	return func() {
		contextualized_sep = sep
	}
}

// creates thunk for init. rec expression related strings
func GenSetRec(head, sep string) func() {
	return func() {
		rec_head, rec_sep = head, sep
	}
}

func Init(initializerThunks ...func()) {
	// force init thunks
	for _, thunk := range initializerThunks {
		thunk() 
	}
}
