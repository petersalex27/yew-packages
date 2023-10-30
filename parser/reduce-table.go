package parser

import (
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/util/iterator"
)

// table of reductions
type ReductionTable struct {
	root *combinerTrieRoot
	// maps lookahead to an ordered list of productions
	table map[ast.Type]productionOrder
}

type ReductionRules interface {
	GetProductionOrder() productionOrder
}

// reduction rule class. reduction rules classified based on lookahead
type ReductionRuleClass struct {
	lookaheads mappable
	productionOrder
}

// return productions in given class
func (r ReductionRuleClass) GetProductionOrder() productionOrder {
	return r.productionOrder
}

// create a shift reduction rule class
func (r ReductionRuleClass) ElseShift() ReductionRuleClass {
	rule_set := r.productionOrder
	rule_set.elseShift = true
	return ReductionRuleClass{lookaheads: r.lookaheads, productionOrder: rule_set}
}

var shiftRuleSet productionOrder = productionOrder{rules: nil, classes: nil, elseShift: true}

type needEndReduction ReductionTable

type ForTypesThrough ast.Type

// Note: mapping from an element of mems that already exists in the table will overwrite
// the previous map!
func (m *ReductionTable) setInTable(ruleset productionOrder, rep ast.Type, mems []ast.Type) {
	rs, found := m.table[rep]
	var in productionOrder
	if found {
		in = Union(rs, ruleset)
	} else {
		in = ruleset
	}

	// iterator is used here and not `range` because the first iteration uses `rep` and
	// not the first value from the iterator
	ty, it := rep, iterator.Iterator(mems)
	for tyExists := true; tyExists; ty, tyExists = it.Next() {
		m.table[ty] = in
	}
}

func (lastType ForTypesThrough) UseReductions(reductionRules ...ReductionRuleClass) needEndReduction {
	m := ReductionTable{table: make(map[ast.Type]productionOrder)}
	m.root = initRoot(ast.Type(lastType))

	for _, rrule := range reductionRules {
		var mp Mapper
		var tys []ast.Type = nil
		if cls, isClass := class(rrule.lookaheads); isClass {
			mp = cls.rep
			tys = m.root.setMems(cls.mems...)
		} else {
			mp, _ = justMapper(rrule.lookaheads)
		}

		ty := m.root.set(mp...)
		m.setInTable(rrule.productionOrder, ty, tys)
	}

	return needEndReduction(m)
}

// production rules to apply once all tokens have been shifted
func (m needEndReduction) Finally(rs productionOrder) ReductionTable {
	ty := ast.None
	if _, found := m.table[ty]; found {
		panic("terminal reduction rule(s) already exist(s)")
	}

	m.table[ty] = rs
	return ReductionTable(m)
}
