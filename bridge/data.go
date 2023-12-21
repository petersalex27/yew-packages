package bridge

import (
	"strings"

	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/nameable"
)

type Data[T nameable.Nameable] struct {
	tag     expr.Const[T]
	Members []JudgmentAsExpression[T, expr.Expression[T]]
}

type Constructor[N nameable.Nameable] func(args []expr.Expression[N]) (data Data[N], success bool)

func applyArg[N nameable.Nameable](constructor, arg expr.Expression[N]) expr.Expression[N] {
	return constructor.(expr.Function[N]).Apply(arg)
}

// assumes length of `args` is the number arguments constructor takes
func applyArgumentsToConstructor[N nameable.Nameable](constructor expr.Function[N], args []expr.Expression[N]) (data Data[N], success bool) {
	// result of applying 0 to len(args) arguments
	var constructed expr.Expression[N] = constructor

	// check if constructor has arguments
	if len(constructor.GetBinders()) != 0 {
		// apply all arguments to function
		constructed = fun.FoldLeft(constructed, args, applyArg[N])
	} else {
		// no arg constructor
		constructed = constructor.GetBound()
	}

	data, success = constructed.(Data[N])
	return data, success
}

// creates data type constructor from function
func AsConstructor[N nameable.Nameable](f expr.Function[N]) Constructor[N] {
	binderLength := len(f.GetBinders())

	return func(args []expr.Expression[N]) (data Data[N], success bool) {
		if success = len(args) == binderLength; success {
			data, success = applyArgumentsToConstructor(f, args)
		}
		return
	}
}

func (data Data[T]) GetTag() expr.Const[T] { return data.tag }

func (data Data[T]) Flatten() []expr.Expression[T] {
	f := (JudgmentAsExpression[T, expr.Expression[T]]).Flatten
	fold := func(l, r []expr.Expression[T]) []expr.Expression[T] {
		return append(l, r...)
	}

	left := data.tag.Flatten()
	right := fun.FoldLeft([]expr.Expression[T]{}, fun.FMap(data.Members, f), fold)

	return append(left, right...)
}

func (data Data[T]) GetReferred() T {
	return data.tag.GetReferred()
}

func MakeData[T nameable.Nameable](tag expr.Const[T], members ...JudgmentAsExpression[T, expr.Expression[T]]) Data[T] {
	return Data[T]{tag, members}
}

func makeData[T nameable.Nameable](tag expr.Const[T], members []JudgmentAsExpression[T, expr.Expression[T]]) Data[T] {
	return Data[T]{tag, members}
}

// string rep. of data type
//
//	 MakeData(MyData, App((+), n, 1), App(whatever, thing, this, 1)).String()
//		=> "(MyData (n + 1) (whatever thing this 1))"
func (data Data[T]) String() string {
	strs := make([]string, 1, 1+len(data.Members))
	strs[0] = data.tag.String()
	for _, mems := range data.Members {
		strs = append(strs, mems.String())
	}

	return "(" + strings.Join(strs, " ") + ")"
}

func (data Data[T]) Rebind() expr.Expression[T] {
	f := func(j JudgmentAsExpression[T, expr.Expression[T]]) JudgmentAsExpression[T, expr.Expression[T]] {
		return j.Rebind().(JudgmentAsExpression[T, expr.Expression[T]])
	}
	return makeData(data.tag, fun.FMap(data.Members, f))
}

func (data Data[T]) Bind(bs expr.BindersOnly[T]) expr.Expression[T] {
	f := func(j JudgmentAsExpression[T, expr.Expression[T]]) JudgmentAsExpression[T, expr.Expression[T]] {
		return j.Bind(bs).(JudgmentAsExpression[T, expr.Expression[T]])
	}
	return makeData(data.tag, fun.FMap(data.Members, f))
}

func (data Data[T]) Replace(v expr.Variable[T], e expr.Expression[T]) (expr.Expression[T], bool) {
	again := false
	f := func(j JudgmentAsExpression[T, expr.Expression[T]]) JudgmentAsExpression[T, expr.Expression[T]] {
		res, tmp := j.Replace(v, e)
		again = again || tmp
		return res.(JudgmentAsExpression[T, expr.Expression[T]])
	}

	return makeData(data.tag, fun.FMap(data.Members, f)), again
}

func (data Data[T]) UpdateVars(gt int, by int) expr.Expression[T] {
	f := func(j JudgmentAsExpression[T, expr.Expression[T]]) JudgmentAsExpression[T, expr.Expression[T]] {
		return j.UpdateVars(gt, by).(JudgmentAsExpression[T, expr.Expression[T]])
	}
	return makeData(data.tag, fun.FMap(data.Members, f))
}

func (data Data[T]) BodyAbstract(v expr.Variable[T], name expr.Const[T]) expr.Expression[T] {
	if name.StrictEquals(data.tag) { // don't abstract data type name to variable
		return data
	}

	f := func(j JudgmentAsExpression[T, expr.Expression[T]]) JudgmentAsExpression[T, expr.Expression[T]] {
		return j.BodyAbstract(v, name).(JudgmentAsExpression[T, expr.Expression[T]])
	}

	return makeData(data.tag, fun.FMap(data.Members, f))
}

func (data Data[T]) Equals(cxt *expr.Context[T], e expr.Expression[T]) bool {
	data2, ok := e.(Data[T])
	if !ok {
		return false
	}

	if len(data.Members) != len(data2.Members) {
		return false
	}

	if !data.tag.Equals(cxt, data2.tag) {
		return false
	}

	for i, mem := range data.Members {
		if !mem.Equals(cxt, data2.Members[i]) {
			return false
		}
	}

	return true
}

func (data Data[T]) StrictEquals(e expr.Expression[T]) bool {
	data2, ok := e.(Data[T])
	if !ok {
		return false
	}

	if len(data.Members) != len(data2.Members) {
		return false
	}

	if !data.tag.StrictEquals(data2.tag) {
		return false
	}

	for i, mem := range data.Members {
		if !mem.StrictEquals(data2.Members[i]) {
			return false
		}
	}

	return true
}

func (data Data[T]) StrictString() string {
	strs := make([]string, 1, 1+len(data.Members))
	strs[0] = data.tag.StrictString()
	for _, mems := range data.Members {
		strs = append(strs, mems.StrictString())
	}

	return "(" + strings.Join(strs, " ") + ")"
}

func (data Data[T]) Again() (expr.Expression[T], bool) {
	return data, false
}

func (data Data[T]) Find(v expr.Variable[T]) bool {
	for _, mem := range data.Members {
		if mem.Find(v) {
			return true
		}
	}
	return false
}

func (data Data[T]) PrepareAsRHS() expr.Expression[T] {
	f := func(j JudgmentAsExpression[T, expr.Expression[T]]) JudgmentAsExpression[T, expr.Expression[T]] {
		return j.PrepareAsRHS().(JudgmentAsExpression[T, expr.Expression[T]])
	}
	return makeData(data.tag, fun.FMap(data.Members, f))
}

func (data Data[T]) Copy() expr.Expression[T] {
	f := func(j JudgmentAsExpression[T, expr.Expression[T]]) JudgmentAsExpression[T, expr.Expression[T]] {
		return j.Copy().(JudgmentAsExpression[T, expr.Expression[T]])
	}
	return makeData(data.tag, fun.FMap(data.Members, f))
}

func (data Data[T]) ForceRequest() expr.Expression[T] {
	f := func(j JudgmentAsExpression[T, expr.Expression[T]]) JudgmentAsExpression[T, expr.Expression[T]] {
		return j.ForceRequest().(JudgmentAsExpression[T, expr.Expression[T]])
	}
	return makeData(data.tag, fun.FMap(data.Members, f))
}

func (data Data[T]) ExtractVariables(gt int) []expr.Variable[T] {
	f := func(j JudgmentAsExpression[T, expr.Expression[T]]) []expr.Variable[T] {
		return j.ExtractVariables(gt)
	}
	app := func(l, r []expr.Variable[T]) []expr.Variable[T] {
		return append(l, r...)
	}
	return fun.FoldLeft([]expr.Variable[T]{}, fun.FMap(data.Members, f), app)
}

func (data Data[T]) Collect() []T {
	f := func(j JudgmentAsExpression[T, expr.Expression[T]]) []T {
		return j.Collect()
	}
	app := func(l, r []T) []T {
		return append(l, r...)
	}
	return fun.FoldLeft([]T{}, fun.FMap(data.Members, f), app)
}
