package types

import (
	expr "github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
)

type TypeJudgement[T nameable.Nameable, E expr.Expression] struct{
	expression E
	ty Type[T]
}

func (j TypeJudgement[T,_]) GetType() Type[T] {
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

func (j TypeJudgement[T,_]) String() string {
	return "(" + j.expression.String() + ": " + j.ty.String() + ")"
}

// Judgement makes the trivial type judgement `𝚪, e: ty ⊢ e: ty`
func Judgement[T nameable.Nameable, E expr.Expression](e E, ty Type[T]) TypeJudgement[T,E] {
	return TypeJudgement[T,E]{
		expression: e,
		ty: ty,
	}
}

func (j TypeJudgement[T,E]) Replace(v Variable[T], m Monotyped[T]) TypeJudgement[T,E] {
	return Judgement(j.expression, MaybeReplace(j.ty, v, m))
}

func Equals[N nameable.Nameable, T, U expr.Expression](j1 TypeJudgement[N,T], j2 TypeJudgement[N,U]) bool {
	return j1.ty.Equals(j2.ty) && j1.expression.Equals(j2.expression)
}

func (j TypeJudgement[T, E]) Collect() []T {
	return j.ty.Collect()
}