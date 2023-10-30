package parser

import (
	"github.com/petersalex27/yew-packages/parser/ast"
)

type needPattern Production

// creates production
//
//	replacement ::= handle
//
// where
//
//	replacement = productionFunction.do(handle)
func (productionFunction NamedProduction) From(handle ...ast.Type) productionInterface {
	return needPattern(Production{productionFunction}).From(handle...)
}

// uses `productionFunction` to create (and return) a value ready
// that can be turned into a reduction rule
func Get(productionFunction func(...ast.Ast) ast.Ast) needPattern {
	return needPattern(Production{ProductionFunction(productionFunction)})
}
