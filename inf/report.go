package inf

import "github.com/petersalex27/yew-packages/nameable"

type errorReport[N nameable.Nameable] struct {
	DuringRule string
	Status
	TermsInvolved []TypeJudgement[N]
}

// creates an errorReport for a failed rule
func makeReport[N nameable.Nameable](duringRule string, status Status, withTerms ...TypeJudgement[N]) errorReport[N] {
	return errorReport[N]{duringRule, status, withTerms}
}
