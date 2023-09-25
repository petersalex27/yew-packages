package types

import (
	expr "github.com/petersalex27/yew-packages/expression"
)

type TypeJudgement[E expr.Expression] struct{
	expression E
	ty Type
}

func (j TypeJudgement[_]) GetType() Type {
	return j.ty
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

func (j TypeJudgement[_]) String() string {
	return "(" + j.expression.String() + ": " + j.ty.String() + ")"
}

// Judgement makes the trivial type judgement `𝚪, e: ty ⊢ e: ty`
func Judgement[E expr.Expression](e E, ty Type) TypeJudgement[E] {
	return TypeJudgement[E]{
		expression: e,
		ty: ty,
	}
}

func (j TypeJudgement[E]) Replace(v Variable, m Monotyped) TypeJudgement[E] {
	return Judgement(j.expression, MaybeReplace(j.ty, v, m))
}

func Equals[T, U expr.Expression](j1 TypeJudgement[T], j2 TypeJudgement[U]) bool {
	return j1.ty.Equals(j2.ty) && j1.expression.Equals(j2.expression)
}