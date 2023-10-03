package parser

import (
	"github.com/petersalex27/yew-packages/parser/ast"
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

type Mapper []ast.Type

func LA(tys ...ast.Type) Mapper { return Mapper(tys) }

type ReductionRules interface {
	GetRuleSet() ruleSet
}

type ReductionRule struct {
	lookaheads Mapper
	ruleSet
}

func (r ReductionRule) GetRuleSet() ruleSet {
	return r.ruleSet
}

func (r ReductionRule) ElseShift() ReductionRule {
	rule_set := r.ruleSet
	rule_set.shiftAtEnd = true
	return ReductionRule{lookaheads: r.lookaheads, ruleSet: rule_set}
}

var shiftRuleSet ruleSet = ruleSet{rules: nil, shiftAtEnd: true}

func (tys Mapper) Shift() ReductionRule {
	return ReductionRule{tys, shiftRuleSet}
}

func (tys Mapper) Then(rs ruleSet) ReductionRule {
	return ReductionRule{tys, rs}
}

type needEndReduction ReduceTable

type ForTypesThrough ast.Type

func (lastType ForTypesThrough) UseReductions(reductionRules ...ReductionRule) needEndReduction {
	m := ReduceTable{table: make(map[ast.Type]ruleSet)}
	m.root = initRoot(ast.Type(lastType))

	for _, rrule := range reductionRules {
		ty := m.root.set(rrule.lookaheads...)
		rs, found := m.table[ty]
		var in ruleSet
		if found {
			in = rs.Union(rrule.ruleSet)
		} else {
			in = rrule.ruleSet
		}
		m.table[ty] = in
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
