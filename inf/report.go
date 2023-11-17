package inf

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
)

type errorReport[N nameable.Nameable] struct {
	DuringRule    string
	Status        Status
	TermsInvolved []TypeJudgement[N]
	Names         []expr.Const[N]
}

// creates an errorReport for a failed rule
func makeReport[N nameable.Nameable](duringRule string, status Status, withTerms ...TypeJudgement[N]) errorReport[N] {
	return errorReport[N]{duringRule, status, withTerms, nil}
}

// creates an errorReport for a failed context lookup
func makeNameReport[N nameable.Nameable](duringRule string, status Status, withNames ...expr.Const[N]) errorReport[N] {
	return errorReport[N]{duringRule, status, nil, withNames}
}
