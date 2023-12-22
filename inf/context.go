package inf

import (
	"github.com/petersalex27/yew-packages/bridge"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/table"
	"github.com/petersalex27/yew-packages/types"
)

// Type a = Data a Int
//
// Data = (\x y -> Data x y): a -> Int -> Type a
type consJudge[N nameable.Nameable] struct {
	forType      types.Polytype[N]
	constructors constructorMapType[N]
}

type constructorMapType[N nameable.Nameable] map[string]types.TypedJudgment[N, expr.Function[N], types.Polytype[N]]

// tries to find constructor named `constructorName` w/in construtor map receiver
func (constructors consJudge[N]) Find(constructorName N) (constructor types.TypedJudgment[N, expr.Function[N], types.Polytype[N]], found bool) {
	if constructors.constructors == nil {
		found = false
	} else {
		constructor, found = constructors.constructors[constructorName.GetName()]
	}
	return
}

func (constructors consJudge[N]) GetType() types.Polytype[N] {
	if constructors.constructors == nil {
		panic("bug: constructor map is uninitialized")
	}

	return constructors.forType
}

type Context[N nameable.Nameable] struct {
	reports     []errorReport[N]
	typeSubs    *table.Table[types.Monotyped[N]]
	exprSubs    *table.Table[expr.Referable[N]]
	consTable   *table.Table[consJudge[N]]
	syms        *table.Table[Symbol[N]]
	TypeContext *types.Context[N]
	ExprContext *expr.Context[N]
}

// convenience method for a type judgment with a new, free type variable; i.e., for an expression e,
//
//	e: newvar
func (cxt *Context[N]) Judge(e expr.Expression[N]) TypeJudgment[N] {
	var newvar types.Type[N] = cxt.TypeContext.NewVar()
	return bridge.Judgment(e, newvar)
}

type ExportableContext[N nameable.Nameable] struct {
	name      N
	consTable *table.Table[consJudge[N]]
	syms      *table.Table[Symbol[N]]
}

func (ecxt *ExportableContext[N]) export(name N, sym Symbol[N]) Status {
	// re-exported?
	_, ok := ecxt.syms.Get(name)
	if ok {
		return IllegalShadow
	}

	// add symbol to table
	ecxt.syms.Add(name, sym)
	return Ok
}

func newConsAndSymsTables[N nameable.Nameable]() (*table.Table[consJudge[N]], *table.Table[Symbol[N]]) {
	return table.NewTable[consJudge[N]](), table.NewTable[Symbol[N]]()
}

func NewExportableContext[N nameable.Nameable]() *ExportableContext[N] {
	cxt := new(ExportableContext[N])
	cxt.consTable, cxt.syms = newConsAndSymsTables[N]()
	return cxt
}

// creates new inf context
func NewContext[N nameable.Nameable]() *Context[N] {
	cxt := new(Context[N])
	cxt.typeSubs = table.NewTable[types.Monotyped[N]]()
	cxt.exprSubs = table.NewTable[expr.Referable[N]]()
	cxt.consTable, cxt.syms = newConsAndSymsTables[N]()
	cxt.ExprContext = expr.NewContext[N]()
	cxt.TypeContext = types.NewContext[N]()
	cxt.reports = []errorReport[N]{}
	return cxt
}

func (cxt *Context[N]) Inst(sigma types.Polytype[N]) types.Monotyped[N] {
	var t types.DependentTyped[N] = sigma.GetBound()
	typeVars := sigma.GetBinders()

	// create new type variables
	vs := fun.FMap(
		typeVars,
		func(v types.Variable[N]) types.Monotyped[N] {
			return cxt.TypeContext.NewVar()
		},
	)

	if d, ok := t.(types.DependentType[N]); ok {
		// replace all bound expression variables w/ new expression variables
		t = d.FreeIndex(cxt.ExprContext)
	}

	// replace all bound variables w/ newly created type variables
	return t.ReplaceDependent(typeVars, vs)
}

func NewTestableContext() *Context[nameable.Testable] {
	cxt := NewContext[nameable.Testable]()
	cxt.TypeContext = cxt.TypeContext.SetNameMaker(nameable.MakeTestable)
	cxt.ExprContext = cxt.ExprContext.SetNameMaker(nameable.MakeTestable)
	return cxt
}

// applies kind and type substitutions to expression and type of judgment respectively
func (cxt *Context[N]) judgmentSubstitution(judge bridge.JudgmentAsExpression[N, expr.Expression[N]]) bridge.JudgmentAsExpression[N, expr.Expression[N]] {
	referable, monotype := GetExpressionAndType[N, expr.Referable[N], types.Monotyped[N]](judge)

	var kindSubResult expr.Expression[N] = cxt.GetKindSub(referable)
	var typeSubResult types.Type[N] = cxt.GetSub(monotype)

	return bridge.Judgment(kindSubResult, typeSubResult)
}

// applies kind substitutions to `postFindKind`
//
// ASSUMPTION: `postFindKind` is
//
//	cxt.findKindSub(someKind) = postFindKind
func (cxt *Context[N]) applyKindSubstitutions(postFindKind expr.Referable[N]) expr.Referable[N] {
	data, isData := postFindKind.(bridge.Data[N])
	if !isData {
		return postFindKind
	}

	memsSubResult := fun.FMap(data.Members, cxt.judgmentSubstitution)
	return bridge.MakeData(data.GetTag(), memsSubResult...)
}

// returns the result of applying all applicable substitutions to `rawKind`.
//
// For example, given substitutions
//
//	> Sub = { n ⟼ 0, k ⟼ Succ n },
//
// and given an input of
//
//	> Succ k
//
// return
//
//	> Succ (Succ 0)
func (cxt *Context[N]) GetKindSub(rawKind expr.Referable[N]) (kind expr.Referable[N]) {
	kind, _ = cxt.findKindSub(rawKind) // returns rawKind if no sub exists
	return cxt.applyKindSubstitutions(kind)
}

func (cxt *Context[N]) GetSub(m types.Monotyped[N]) (out types.Monotyped[N]) {
	var found bool

	out, found = cxt.findSub(m)

	if !found {
		out = m
	}

	if function, ok := out.(types.TypeFunction[N]); ok {
		out = function.Rebuild(cxt.GetSub, cxt.GetKindSub)
	}

	return out
}

// first return value is base substitution for `m` (or `m` itself when second return value is false)
//
// second return value is true iff `m` is a variable and `m` has a registered substitution
func (cxt *Context[N]) findSub(m types.Monotyped[N]) (out types.Monotyped[N], found bool) {
	found = false
	if nm, ok := m.(types.Variable[N]); ok {
		out, found = cxt.typeSubs.Get(nm)
	}

	if !found {
		out = m
	}

	return
}

// first return value is base substitution for `e` (or `e` itself when second return value is false)
//
// second return value is true iff `e` is a variable and `e` has a registered substitution
func (cxt *Context[N]) findKindSub(e expr.Referable[N]) (out expr.Referable[N], found bool) {
	found = false
	if v, ok := e.(expr.Variable[N]); ok {
		out, found = cxt.exprSubs.Get(v.GetReferred())
	}

	if !found {
		out = e
	}

	return
}

// returns representative for type equiv. class
func (cxt *Context[N]) Find(m types.Monotyped[N]) (representative types.Monotyped[N]) {
	representative, _ = cxt.findSub(m)
	return
}

// returns representative for kind equiv. class
func (cxt *Context[N]) FindKind(e expr.Referable[N]) (representative expr.Referable[N]) {
	representative, _ = cxt.findKindSub(e)
	return
}
