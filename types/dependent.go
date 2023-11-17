package types

import (
	expr "github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
	str "github.com/petersalex27/yew-packages/stringable"
)

type DependentTyped[T nameable.Nameable] interface {
	Type[T]
	ReplaceDependent(vs []Variable[T], with []Monotyped[T]) Monotyped[T]
}

// Dependent Type: `(mapval (a: A) (b: B) ..) . (F a b ..)`
type DependentType[T nameable.Nameable] struct {
	mapval   []TypeJudgement[T, expr.Variable[T]]
	Function TypeFunction[T]
}

// creates a dependent type that depends on variables in `mapall` and is
// instantiated by indexing `typeFunc`. This is sometimes called a "dependent
// function type" or "dependent product type".
//
// example (removing type params from call for clarity):
//
//	n, Uint := expr.Var("n"), MakeConst("Uint")
//	mapval := []TypeJudgement{Judgement(n, Uint)}
//	ArraySelector, a := MakeConst("ArraySelector"), Var("a")
//	ArraySelector_a := Apply(ArraySelector, a)
//	MakeDependentType(mapval, ArraySelector_a).
//		String() == "mapval (n: Uint) . (ArraySelector a)"
func MakeDependentType[T nameable.Nameable](mapval []TypeJudgement[T, expr.Variable[T]], typeFunc TypeFunction[T]) DependentType[T] {
	return DependentType[T]{
		mapval:   mapval,
		Function: typeFunc,
	}
}

type mapNeedsBound[T nameable.Nameable] []TypeJudgement[T, expr.Variable[T]]

func Map[T nameable.Nameable](ts ...TypeJudgement[T, expr.Variable[T]]) mapNeedsBound[T] {
	return ts
}

func (mapval mapNeedsBound[T]) To(typeFunc TypeFunction[T]) DependentType[T] {
	return MakeDependentType(mapval, typeFunc)
}

func (d DependentType[T]) String() string {
	if len(d.mapval) == 0 {
		return d.Function.String()
	}

	return "mapval " + str.Join(d.mapval, str.String(" ")) + " . " + d.Function.String()
}

// index dependent type, making it a dependent type index
func (d DependentType[T]) FreeIndex(cxt *expr.Context[T]) TypeFunction[T] {
	return d.Function.SubVars(
		d.mapval, 
		cxt.NumNewReferable(len(d.mapval)),
	)
}

// just assumes e: A
//
//	((mapval (a: A) (b: B) ..) . C) -> ((mapval (b: B) ..) . (C e))
func (cxt *Context[T]) InstantiateKind(d DependentType[T], e expr.Referable[T]) DependentTyped[T] {
	if len(d.mapval) == 0 {
		return d.Function
	}

	preSub := d.mapval[0:1]
	postSub := []expr.Referable[T]{e}
	d.Function = d.Function.SubVars(preSub, postSub)
	if len(d.mapval) == 1 {
		return d.Function
	}
	d.mapval = d.mapval[1:]
	return d
}

// This test exact equality, not judgemental equality. For example,
//
//	(mapval (a: A) (b: B) . (C b)) != (mapval (b: B) . (C b))
//	despite the two being equiv. in some (probably useful) sense.
//
// Additionally, the following is not equiv. either:
//
//	(mapval (a: A) (b: B) . (C b)) != (mapval (b: B) (a: A) . (C b))
func (d DependentType[T]) Equals(t Type[T]) bool {
	d2, ok := t.(DependentType[T])
	if !ok {
		return false
	}

	if len(d.mapval) != len(d2.mapval) {
		return false
	}

	for i, judge := range d.mapval {
		if !judge.ty.Equals(d2.mapval[i].ty) {
			return false
		}
		// nil is okay here because variables don't require context object for equality
		if !judge.expression.Equals(nil, d2.mapval[i].expression) {
			return false
		}
	}

	return d.Function.Equals(d2.Function)
}

func (d DependentType[T]) Generalize(cxt *Context[T]) Polytype[T] {
	return Polytype[T]{
		typeBinders: cxt.MakeDummyVars(1),
		bound:       d,
	}
}

func (d DependentType[T]) Collect() []T {
	var res []T = []T{}
	if len(d.mapval) != 0 {
		res = d.mapval[0].Collect()
		for _, m := range d.mapval[1:] {
			res = append(res, m.Collect()...)
		}
	}
	res = append(res, d.Function.Collect()...)
	return res
}

// replaces all occ. of each v in `vs` with corr. m in `ms`
func (d DependentType[T]) ReplaceDependent(vs []Variable[T], ms []Monotyped[T]) Monotyped[T] {
	// redo kind-variable binders's types
	// mapval := fun.FMap(
	// 	d.mapval,
	// 	func(tj TypeJudgement[T, expr.Variable[T]]) TypeJudgement[T, expr.Variable[T]] {
	// 		mono, _ := tj.ty.(Monotyped[T]) // should always pass since dependent types only have monotype binders
	// 		var ty Type[T] = mono.ReplaceDependent(vs, ms)
	// 		return Judgement(tj.expression, ty) // updated judgement
	// 	},
	// )

	// replace vars in type function
	inst, _ := d.Function.ReplaceDependent(vs, ms).(TypeFunction[T])

	// index (call) type function
	return inst //.AsFreeInstance(mapval)
}
