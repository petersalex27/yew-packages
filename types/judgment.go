package types

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
)

type TypedJudgment[N nameable.Nameable, E expr.Expression[N], T Type[N]] struct {
	expression E
	typing     T
}

func TypedJudge[N nameable.Nameable, E expr.Expression[N], T Type[N]](e E, t T) TypedJudgment[N, E, T] {
	return TypedJudgment[N, E, T]{e, t}
}

// returns typing of judgment
//
// e.g., given a judgment `e: t`
//
//	return t
func (judgment TypedJudgment[_, _, T]) GetType() T { return judgment.typing }

// returns expression of judgment
//
// e.g., given a judgment `e: t`
//
//	return e
func (judgment TypedJudgment[_, E, _]) GetExpression() E { return judgment.expression }

// splits the judgment into its components
func (judgment TypedJudgment[N, _, _]) GetExpressionAndType() (expr.Expression[N], Type[N]) {
	return judgment.expression, judgment.typing
}

func (judgment TypedJudgment[N, E, T]) String() string {
	return "(" + judgment.expression.String() + ": " + judgment.typing.String() + ")"
}

func (judgment TypedJudgment[N, E, _]) AsTypeJudgment() TypeJudgment[N, E] {
	return TypeJudgment[N, E]{judgment.expression, judgment.typing}
}

func (judgment TypedJudgment[N, E, _]) MakeJudgment(e E, ty Type[N]) ExpressionJudgment[N, E] {
	return TypedJudgment[N, E, Type[N]]{e, ty}
}

func JudgmentEquals[N nameable.Nameable, E expr.Expression[N], T Type[N]](j0 ExpressionJudgment[N, E], j1 ExpressionJudgment[N, E]) bool {
	var strictJudgment0, strictJudgment1 TypedJudgment[N, E, T]
	var ok0, ok1 bool

	// assert go type
	strictJudgment0, ok0 = AsJudgment[N, E, T](j0)
	strictJudgment1, ok1 = AsJudgment[N, E, T](j1)
	if !(ok0 && ok1) {
		return false
	}

	// split judgment into comp.
	e0, t0 := strictJudgment0.GetExpressionAndType()
	e1, t1 := strictJudgment1.GetExpressionAndType()

	// test equality of components
	return e0.StrictEquals(e1) && t0.Equals(t1)
}

func (judgment TypedJudgment[N, E, T]) Collect() []N {
	return append(judgment.expression.Collect(), judgment.typing.Collect()...)
}
