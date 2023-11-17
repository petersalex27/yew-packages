package types

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
)

// possibly indexed dependent type function
type TypeFunction[N nameable.Nameable] interface {
	Monotyped[N]
	//AsFreeInstance(vs []TypeJudgement[N, expr.Variable[N]]) TypeFunction[N]
	SubVars(preSub []TypeJudgement[N, expr.Variable[N]], postSub []expr.Referable[N]) TypeFunction[N]
	FunctionAndIndexes() (function Application[N], indexes Indexes[N])
	Rebuild(findMono func(Monotyped[N]) Monotyped[N], findKind func(expr.Referable[N]) expr.Referable[N]) TypeFunction[N]
}
