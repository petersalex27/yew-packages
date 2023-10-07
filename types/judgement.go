package types

import (
	expr "github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/stringable"
)

type ExpressionJudgement[T nameable.Nameable, E expr.Expression[T]] interface {
	stringable.Stringable
	collectable[T]
	asTypeJudgement() TypeJudgement[T, E]
	MakeJudgement(E, Type[T]) ExpressionJudgement[T, E]
}

func GetType[T nameable.Nameable, E expr.Expression[T]](j ExpressionJudgement[T, E]) Type[T] {
	return j.asTypeJudgement().ty
}

func GetExpression[T nameable.Nameable, E expr.Expression[T]](j ExpressionJudgement[T, E]) E {
	return j.asTypeJudgement().expression
}

func Replace[T nameable.Nameable, E expr.Expression[T]](j ExpressionJudgement[T, E], v Variable[T], m Monotyped[T]) ExpressionJudgement[T,E] {
	ex, ty := GetExpression(j), GetType(j)
	return j.MakeJudgement(ex, MaybeReplace(ty, v, m))
}

func JudgesEquals[N nameable.Nameable, T, U expr.Expression[N]](j1 ExpressionJudgement[N,T], j2 ExpressionJudgement[N,U]) bool {
	return Equals(j1.asTypeJudgement(), j2.asTypeJudgement())
}

func Collect[T nameable.Nameable, E expr.Expression[T]](j ExpressionJudgement[T, E]) []T {
	return j.Collect()
}

func String[T nameable.Nameable, E expr.Expression[T]](j ExpressionJudgement[T, E]) string {
	return j.String()
}

type TypeJudgement[T nameable.Nameable, E expr.Expression[T]] struct{
	expression E
	ty Type[T]
}

func (j TypeJudgement[T, E]) asTypeJudgement() TypeJudgement[T, E] {
	return j
}

func (j TypeJudgement[T,_]) GetType() Type[T] {
	return j.ty
}

func (j TypeJudgement[_, E]) GetExpression() E {
	return j.expression
}

/*
decons: [a; n+1] -> (a, [a; n])
decons (x::xs) = (x, xs)

typeof(decons [1, 2, 3, 4])
decons: forall a . map (Arr a)(n: Uint) . (Arr a)(n+1) -> (a, (Arr a)(n))
[_, ..]: forall a . map (Arr a)(n: Uint) . (Arr a)(n)
1: Int
2: Int
3: Int
4: Int
let x = [] in 
	(Cons 1) . (Cons 2) . (Cons 3) . (Cons 4 x)
	(Cons 4 x): (Arr Int)(1)
	(Cons 3) . (Cons 4 x): (Arr Int)(2)
	...
	: (Arr Int)(4)
decons$(Arr Int)(4): (Arr Int)(4) -> (Int, (Arr Int)(3))
*/

func (j TypeJudgement[T,_]) String() string {
	return "(" + j.expression.String() + ": " + j.ty.String() + ")"
}

// Judgement makes the trivial type judgement `𝚪, e: ty ⊢ e: ty`
func Judgement[T nameable.Nameable, E expr.Expression[T]](e E, ty Type[T]) TypeJudgement[T,E] {
	return TypeJudgement[T,E]{
		expression: e,
		ty: ty,
	}
}

func (TypeJudgement[T, E]) MakeJudgement(e E, ty Type[T]) ExpressionJudgement[T, E] {
	return Judgement[T,E](e, ty)
}

func (j TypeJudgement[T,E]) Replace(v Variable[T], m Monotyped[T]) TypeJudgement[T,E] {
	return Judgement(j.expression, MaybeReplace(j.ty, v, m))
}

func Equals[N nameable.Nameable, T, U expr.Expression[N]](j1 TypeJudgement[N,T], j2 TypeJudgement[N,U]) bool {
	return j1.ty.Equals(j2.ty) && 
			j1.expression.StrictEquals(j2.expression)
			// TODO: ??? j1.expression.Equals(j2.expression) instead???
}

func (j TypeJudgement[T, E]) Collect() []T {
	return append(j.expression.Collect(), j.ty.Collect()...)
}