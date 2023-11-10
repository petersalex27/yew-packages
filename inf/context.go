package inf

import (
	"github.com/petersalex27/yew-packages/bridge"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/table"
	"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-packages/util/stack"
)

type Context[T nameable.Nameable] struct {
	es []error
	syms *table.Table[Symbol[T]]
	removeActions *stack.Stack[T]
	typeContext *types.Context[T]
	exprContext *expr.Context[T]
}

// creates new inf context
func NewContext[T nameable.Nameable]() *Context[T] {
	cxt := new(Context[T])
	cxt.removeActions = stack.NewStack[T](8) // 8 is arbitrary
	cxt.syms = table.NewTable[Symbol[T]]()
	cxt.exprContext = expr.NewContext[T]()
	cxt.typeContext = types.NewContext[T]()
	cxt.es = []error{}
	return cxt
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
	if found{
		judgedName = sym.Get()
	}
	return
}

func (cxt *Context[T]) RemoveLater(name expr.Const[T]) {
	cxt.removeActions.Push(name.Name)
}
