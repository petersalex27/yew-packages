package parser

import (
	"sort"

	"github.com/petersalex27/yew-packages/parser/ast"
)

type Reduction func(nodes ...ast.Ast) ast.Ast

type rule struct {
	pattern
	Reduction
}

type pattern []ast.Type

/*
RuleSet(
	Rule(Id, TypeJudgement, Id, Assign, Value).Reduce(assignment),
	Rule(Id, Assign, Value).Reduce(assignment),
)
*/

func Rule(types ...ast.Type) pattern { return append(pattern{}, types...) }

type ruleSet []rule

func (rs ruleSet) Less(i, j int) bool { return len(rs[i].pattern) < len(rs[j].pattern) }

func (rs ruleSet) Len() int { return len(rs) }

func (rs ruleSet) Swap(i, j int) { rs[i], rs[j] = rs[j], rs[i] }

func (rs ruleSet) Union(rs2 ruleSet) ruleSet {
	out := make(ruleSet, 0, len(rs)+len(rs2))
	out = append(out, rs...)
	out = append(out, rs2...)
	sort.Sort(out)
	return out
}

func RuleSet(rules ...rule) ruleSet {
	out := append(ruleSet{}, rules...)
	sort.Sort(out)
	return out
}

func (p pattern) Reduce(f Reduction) rule { return rule{p, f} }

func (p pattern) equals_len_known(q []ast.Ast) bool {
	for i := range p {
		if p[i] != q[i].NodeType() {
			return false
		}
	}
	return true
}

type ReduceTable struct {
	table map[ast.Type]ruleSet
}

func MakeTable(table map[ast.Type]ruleSet) ReduceTable {
	return ReduceTable{
		table: table,
	}
}

func (rt *ReduceTable) ReductionExists(t ast.Type) bool {
	r, found := rt.table[t]
	return found && len(r) != 0
}