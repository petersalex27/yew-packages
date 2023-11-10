package types

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
)

type TypedJudgement[N nameable.Nameable, E expr.Expression[N], T Type[N]] struct {
	expression E
	typing T
}

// returns typing of judgement
//
// e.g., given a judgement `e: t`
//		return t
func (judgement TypedJudgement[_,_,T]) GetType() T { return judgement.typing }

// returns expression of judgement
//
// e.g., given a judgement `e: t`
//		return e
func (judgement TypedJudgement[_,E,_]) GetExpression() E { return judgement.expression }

// splits the judgement into its components
func (judgement TypedJudgement[_,E,T]) GetExpressionAndType() (E, T) {
	return judgement.expression, judgement.typing
}


func (judgement TypedJudgement[N, E, T]) String() string {
	return "(" + judgement.expression.String() + ": " + judgement.typing.String() + ")"
}

func (judgement TypedJudgement[N,E,_]) AsTypeJudgement() TypeJudgement[N, E] {
	return TypeJudgement[N, E]{judgement.expression, judgement.typing}
}

func (judgement TypedJudgement[N,E,_]) MakeJudgement(e E, ty Type[N]) ExpressionJudgement[N, E] {
	return TypedJudgement[N,E,Type[N]]{e,ty}
}

func JudgementEquals[N nameable.Nameable, E expr.Expression[N], T Type[N]](j0 ExpressionJudgement[N,E], j1 ExpressionJudgement[N,E]) bool {
	var strictJudgement0, strictJudgement1 TypedJudgement[N,E,T]
	var ok0, ok1 bool

	// assert go type
	strictJudgement0, ok0 = AsJudgement[N,E,T](j0)
	strictJudgement1, ok1 = AsJudgement[N,E,T](j1)
	if !(ok0 && ok1) {
		return false
	}

	// split judgement into comp.
	e0, t0 := strictJudgement0.GetExpressionAndType()
	e1, t1 := strictJudgement1.GetExpressionAndType()

	// test equality of components
	return e0.StrictEquals(e1) && t0.Equals(t1)
}

func (judgement TypedJudgement[N,E,T]) Collect() []N {
	return append(judgement.expression.Collect(), judgement.typing.Collect()...)
}