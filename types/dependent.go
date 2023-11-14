package types

import (
	expr "github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/nameable"
	str "github.com/petersalex27/yew-packages/stringable"
)

type DependentTyped[T nameable.Nameable] interface {
	Type[T]
	FreeInstantiation(cxt *Context[T]) Monotyped[T]
	ReplaceDependent(vs []Variable[T], with []Monotyped[T]) Monotyped[T]
}

//type DependentTypeFunction[T nameable.Nameable] Application[T]

// Dependent Type: `(mapcal (a: A) (b: B) ..) . (F a b ..)`
type DependentType[T nameable.Nameable] struct {
	mapval []TypeJudgement[T, expr.Variable[T]]
	DependentTypeInstance[T]
}

// creates a dependent type that depends on variables in `mapall` and is 
// instantiated by indexing `typeFunc`. This is sometimes called a "dependent 
// function type" or "dependent product type".
//
// example (removing type params from call for clarity):
//
//		n, Uint := expr.Var("n"), MakeConst("Uint")
//		mapval := []TypeJudgement{Judgement(n, Uint)}
//		ArraySelector, a := MakeConst("ArraySelector"), Var("a")
//		ArraySelector_a := Apply(ArraySelector, a)
//		MakeDependentType(mapval, ArraySelector_a).
//			String() == "mapval (n: Uint) . (ArraySelector a)"
func MakeDependentType[T nameable.Nameable](mapval []TypeJudgement[T, expr.Variable[T]], typeFunc Application[T]) DependentType[T] {
	return DependentType[T]{
		mapval: mapval,
		DependentTypeInstance: DependentTypeInstance[T]{
			Application: typeFunc,
			index: nil,
		},
	}
}

func (d DependentType[T]) String() string {
	return "mapval " + str.Join(d.mapval, str.String(" ")) + " . " + d.DependentTypeInstance.String()
}

// index dependent type, making it a dependent type index
func (d DependentType[T]) FreeIndex(cxt *expr.Context[T]) DependentTypeInstance[T] {
	return d.DependentTypeInstance.FreeInstantiateKinds(cxt, d.mapval...)
}

func kindInstantiation[T nameable.Nameable](d DependentType[T], defaultElem expr.Referable[T]) DependentTypeInstance[T] {
	index := make([]ExpressionJudgement[T, expr.Referable[T]], len(d.mapval))
	for i := range index {
		var elem expr.Referable[T]
		if defaultElem == nil {
			elem = d.mapval[i].expression
		} else {
			elem = defaultElem
		}
		index[i] = Judgement(elem, d.mapval[i].ty)
	}
	return DependentTypeInstance[T]{
		Application: d.DependentTypeInstance.Application,
		index:       index,
	}
}

func (d DependentType[T]) KindInstantiation() DependentTypeInstance[T] {
	return kindInstantiation(d, nil)
}

// just assumes e: A
// 	((mapval (a: A) (b: B) ..) . C) -> ((mapval (b: B) ..) . (C e))
func (cxt *Context[T]) InstantiateKind(d DependentType[T], e expr.Referable[T]) DependentTyped[T] {
	inst := d.DependentTypeInstance
	ty := d.mapval[0].ty // type of expression should be type of variable being replaced
	index := make([]ExpressionJudgement[T, expr.Referable[T]], len(inst.index)+1)
	copy(index, inst.index)
	index[len(inst.index)] = (FreeJudgement[T, expr.Referable[T]]{}).MakeJudgement(e, ty)
	
	out := DependentType[T]{
		mapval: d.mapval[1:],
		DependentTypeInstance: DependentTypeInstance[T]{
			Application: d.Application,
			index: index,
		},
	}

	if len(out.mapval) == 0 {
		return out.DependentTypeInstance
	}
	return out
}

func (d DependentType[T]) FreeInstantiation(cxt *Context[T]) Monotyped[T] {
	v := expr.Var(cxt.makeName("_"))
	return kindInstantiation(d, expr.Referable[T](v))
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

	return d.DependentTypeInstance.Equals(d2.DependentTypeInstance)
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
	res = append(res, d.DependentTypeInstance.Collect()...)
	return res
}

// replaces all occ. of each v in `vs` with corr. m in `ms`
func (d DependentType[T]) ReplaceDependent(vs []Variable[T], ms []Monotyped[T]) Monotyped[T] {
	// redo kind-variable binders's types
	mapall := fun.FMap(
		d.mapval, 
		func(tj TypeJudgement[T, expr.Variable[T]]) TypeJudgement[T, expr.Variable[T]] {
			mono, _ := tj.ty.(Monotyped[T]) // should always pass since dependent types only have monotype binders
			var ty Type[T] = mono.ReplaceDependent(vs, ms)
			return Judgement(tj.expression, ty) // updated judgement
		},
	)

	// replace vars in type function
	inst, _ := d.DependentTypeInstance.ReplaceDependent(vs, ms).(DependentTypeInstance[T])

	// create free kind variables (removing binder makes them free, which is what's done here)
	freeExprVars := fun.FMap(
		mapall, 
		func (tj TypeJudgement[T, expr.Variable[T]]) expr.Variable[T] {
			return tj.expression
		},
	)

	// index (call) type function
	return inst.AsFreeInstance(freeExprVars, mapall)
}
