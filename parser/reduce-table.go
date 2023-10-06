package parser

import (
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/util/iterator"
)

type ReduceTable struct {
	root *combinerTrieRoot
	//classes map[ast.Type]reps
	table map[ast.Type]ruleSet
}

type reps []ast.Type

/*
func (rt *ReduceTable) SetClass(representative ast.Type, members ...ast.Type) {
	for _, member := range members {
		rs, found := rt.classes[member]
		if found {
			rs = append(rs, representative)
		} else {
			rs = reps{representative}
		}
		rt.classes[member] = rs
	}
}

func (rt *ReduceTable) getMembers(ty ast.Type) reps {
	if ty2, found := rt.classes[ty]; found {
		return ty2
	}
	return reps{ty}
}
*/

func (rt *ReduceTable) Match(pat pattern, nodes ...ast.Ast) bool {
	if len(pat) != len(nodes) {
		return false
	}

	for i, node := range nodes {
		if pat[i] != node.NodeType() {
			return false
		}
	}
	return true

	/*
		for i, node := range nodes {
			rs := rt.getMembers(node.NodeType())
			ok := false
			for _, r := range rs {
				if pat[i] == r {
					ok = true
					break
				}
			}
			if !ok {
				return false
			}
		}
		return true*/
}

type ReductionRules interface {
	GetRuleSet() ruleSet
}

type ReductionRule struct {
	lookaheads mappable
	ruleSet
}

func (r ReductionRule) GetRuleSet() ruleSet {
	return r.ruleSet
}

func (r ReductionRule) ElseShift() ReductionRule {
	rule_set := r.ruleSet
	rule_set.elseShift = true
	return ReductionRule{lookaheads: r.lookaheads, ruleSet: rule_set}
}

var shiftRuleSet ruleSet = ruleSet{rules: nil, elseShift: true}

type needEndReduction ReduceTable

type ForTypesThrough ast.Type

// Note: mapping from an element of mems that already exists in the table will overwrite
// the previous map!
func (m *ReduceTable) setInTable(ruleset ruleSet, rep ast.Type, mems []ast.Type) {
	rs, found := m.table[rep]
	var in ruleSet
	if found {
		in = rs.Union(ruleset)
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

func (lastType ForTypesThrough) UseReductions(reductionRules ...ReductionRule) needEndReduction {
	m := ReduceTable{table: make(map[ast.Type]ruleSet)}
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
		m.setInTable(rrule.ruleSet, ty, tys)
	}

	return needEndReduction(m)
}

func (m needEndReduction) Finally(rs ruleSet) ReduceTable {
	ty := ast.None
	if _, found := m.table[ty]; found {
		panic("terminal reduction rule(s) already exist(s)")
	}

	m.table[ty] = rs
	return ReduceTable(m)
}
