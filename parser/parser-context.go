package parser

import (
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/util/stack"
)

// context for parsers that can use to change how they parse
type ParserContext struct {
	// flag controlling table overwrite ability (default: false)
	allowTableOverwrites bool
	// table started on and fallback table
	defaultTable ReductionTable
	// actual table used
	currentTable ReductionTable
	// all useable tables
	tables map[ast.Type]ReductionTable
	// used to save nodes for later use
	stack stack.Stack[ast.Ast]
	// reduction for parser to use when signaled to do so
	reduction Productions
}

// SwitchTable tries to switch the table used by the parser based on `index`.
// if a table is not mapped at `index`, then the default table is used and
// SwitchTable returns false. If a table is mapped at `index`, SwitchTable
// returns true
func (cxt *ParserContext) SwitchTable(index ast.Type) (found bool) {
	cxt.currentTable, found = cxt.tables[index]
	if !found {
		cxt.currentTable = cxt.defaultTable
	}
	return
}

// MapTable adds a new ReduceTable to the map of ReduceTables. If table
// overwrites are allowed, then this function just adds the table mapping to it
// with `ty`. If table overwrites are not allowed (this is the default), then
// this function panics when trying to map from `ty` when `ty` already maps to
// a table
func (cxt *ParserContext) MapTable(ty ast.Type, table ReductionTable) {
	if !cxt.allowTableOverwrites {
		_, found := cxt.tables[ty]
		if found {
			panic("tried to overwrite table")
		}
	}
	cxt.tables[ty] = table
}

// adds a node to the top of the context stack
func (cxt *ParserContext) Push(node ast.Ast) {
	cxt.stack.Push(node)
}

// removes and returns the top node of the context stack
func (cxt *ParserContext) Pop() (node ast.Ast, stat stack.StackStatus) {
	return cxt.stack.Pop()
}

// returns the top node of the context stack (but does not remove it)
func (cxt *ParserContext) Peek() (node ast.Ast, stat stack.StackStatus) {
	return cxt.stack.Peek()
}

// sets reduction for parser to use when signaled. `reduction` can be nil: this
// sets a reduction that does nothing
func (cxt *ParserContext) SetReduction(reduction Productions) {
	cxt.reduction = reduction
}

// creates a new parserContext with a default map set to `table` and adds the
// map `ty -> table`. `mapHoldsAtLeast` sets the minimum inital table capacity
// for the parserContext's table map. `mapHoldsAtLeast` is ignored when its
// value is less than 1
func makeParserContext(ty ast.Type, table ReductionTable, mapHoldsAtLeast int) ParserContext {
	var tableMap map[ast.Type]ReductionTable
	// make table map based on arg
	if mapHoldsAtLeast < 1 {
		tableMap = make(map[ast.Type]ReductionTable)
	} else {
		tableMap = make(map[ast.Type]ReductionTable, mapHoldsAtLeast)
	}

	// init
	cxt := ParserContext{
		allowTableOverwrites: false,
		defaultTable:         table,
		currentTable:         table,
		tables:               tableMap,
	}

	cxt.MapTable(ty, table)

	return cxt
}
