package parser

import (
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/util/stack"
)

// context for parsers can use to change how they parse
type ParserContext struct {
	// flag controlling table overwrite ability (default: false)
	allowTableOverwrites bool

	// table started on and fallback table
	defaultTable, 
	// actual table used
	currentTable ReduceTable  
	// all useable tables 
	tables map[ast.Type]ReduceTable

	// used to save nodes for later use
	stack stack.Stack[ast.Ast]

	// reduction for parser to use when signaled to do so
	reduction Reduction
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
func (cxt *ParserContext) MapTable(ty ast.Type, table ReduceTable) {
	if !cxt.allowTableOverwrites {
		_, found := cxt.tables[ty]
		if found {
			panic("tried to overwrite table")
		}
	}
	cxt.tables[ty] = table
}

func (cxt *ParserContext) Push(node ast.Ast) {
	cxt.stack.Push(node)
}

func (cxt *ParserContext) Pop() (node ast.Ast, stat stack.StackStatus) {
	node, stat = cxt.stack.Pop()
	return
}

func (cxt *ParserContext) Peek() (node ast.Ast, stat stack.StackStatus) {
	node, stat = cxt.stack.Peek()
	return
}

// sets reduction for parser to use when signaled. `reduction` can be nil: this
// sets a reduction that does nothing
func (cxt *ParserContext) SetReduction(reduction Reduction) {
	cxt.reduction = reduction
}

// creates a new parserContext with a default map set to `table` and adds the 
// map `ty -> table`. `mapHoldsAtLeast` sets the minimum inital table capacity
// for the parserContext's table map. `mapHoldsAtLeast` is ignored when its 
// value is less than 1
func makeParserContext(ty ast.Type, table ReduceTable, mapHoldsAtLeast int) ParserContext {
	var tableMap map[ast.Type]ReduceTable
	// make table map based on arg
	if mapHoldsAtLeast < 1 {
		tableMap = make(map[ast.Type]ReduceTable)
	} else {
		tableMap = make(map[ast.Type]ReduceTable, mapHoldsAtLeast)
	}

	// init
	cxt := ParserContext{
		allowTableOverwrites: false,
		defaultTable: table,
		currentTable: table,
		tables: tableMap,
	}

	cxt.MapTable(ty, table)

	return cxt
}
