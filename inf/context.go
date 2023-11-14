package inf

import (
	"github.com/petersalex27/yew-packages/bridge"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/table"
	"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-packages/util/stack"
)

// Type a = Data a Int
//
// Data = (\x y -> Data x y): a -> Int -> Type a
type consJudge[N nameable.Nameable] struct {
	types.TypedJudgement[N, expr.Function[N], types.Monotyped[N]]
}

type Context[T nameable.Nameable] struct {
	reports       []errorReport[T]
	typeSubs      *table.Table[types.Monotyped[T]]
	consTable     *table.Table[consJudge[T]]
	syms          *table.Table[Symbol[T]]
	removeActions *stack.Stack[T]
	typeContext   *types.Context[T]
	exprContext   *expr.Context[T]
}

// creates new inf context
func NewContext[T nameable.Nameable]() *Context[T] {
	cxt := new(Context[T])
	cxt.typeSubs = table.NewTable[types.Monotyped[T]]()
	cxt.consTable = table.NewTable[consJudge[T]]()
	cxt.removeActions = stack.NewStack[T](8) // 8 is arbitrary
	cxt.syms = table.NewTable[Symbol[T]]()
	cxt.exprContext = expr.NewContext[T]()
	cxt.typeContext = types.NewContext[T]()
	cxt.reports = []errorReport[T]{}
	return cxt
}

func (cxt *Context[T]) appendReport(report errorReport[T]) {
	cxt.reports = append(cxt.reports, report)
}

func (cxt *Context[T]) GetReports() []errorReport[T] {
	return cxt.reports
}

func (cxt *Context[T]) HasErrors() bool {
	return len(cxt.reports) != 0
}

func NewTestableContext() *Context[nameable.Testable] {
	cxt := NewContext[nameable.Testable]()
	cxt.typeContext = cxt.typeContext.SetNameMaker(nameable.MakeTestable)
	cxt.exprContext = cxt.exprContext.SetNameMaker(nameable.MakeTestable)
	return cxt
}

// removes name binding from context
func (cxt *Context[T]) Remove(name expr.Const[T]) {
	key := name.Name
	sym, ok := cxt.syms.Get(key)
	if !ok {
		// TODO: do nothing, ig?
		return
	}

	// unshadow/remove sym
	remove := sym.Unshadow()
	if remove {
		// symbol is not shadowed, remove it
		cxt.syms.Remove(key)
	}
}

// adds name judging it to have the given type
func (cxt *Context[T]) AddWithType(name expr.Const[T], ty types.Type[T]) {
	key := name.Name
	sym, ok := cxt.syms.Get(key)
	if !ok {
		sym = MakeSymbol[T]()
	}
	// create/shadow symbol
	sym.Shadow(name, ty)
	// add symbol to table
	cxt.syms.Add(name.Name, sym)
}

func (cxt *Context[T]) Add(name expr.Const[T]) {
	ty := cxt.typeContext.NewVar()
	cxt.AddWithType(name, ty)
}

// tries to find symbol w/ name
func (cxt *Context[T]) Get(name expr.Const[T]) (judgedName bridge.JudgementAsExpression[T, expr.Const[T]], found bool) {
	key := name.Name
	var sym Symbol[T]
	sym, found = cxt.syms.Get(key)
	if found {
		judgedName = sym.Get()
	}
	return
}

func (cxt *Context[T]) RemoveLater(name expr.Const[T]) {
	cxt.removeActions.Push(name.Name)
}

// returns representative for equiv. class
func (cxt *Context[T]) Find(m types.Monotyped[T]) (out types.Monotyped[T]) {
	found := false // true iff m
	if nm, ok := m.(nameableMonotype[T]); ok {
		out, found = cxt.typeSubs.Get(nm)
	}

	if !found {
		out = m // finds itself
	}

	return out
}
