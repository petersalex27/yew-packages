package inf

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/types"
)

type Conclusion[N nameable.Nameable, E expr.Expression[N], T types.Type[N]] struct {
	judgement types.TypedJudgement[N, E, T]
	Status
}

func (c Conclusion[N, E, T]) Judgement() types.TypedJudgement[N, E, T] { return c.judgement }

func (c Conclusion[N, E, T]) String() string {
	if c.NotOk() {
		return "_: ⊥"
	}
	return c.judgement.String()
}

func Conclude[N nameable.Nameable, E expr.Expression[N], T types.Type[N]](e E, t T) Conclusion[N, E, T] {
	return Conclusion[N, E, T]{types.TypedJudge[N](e, t), Ok}
}

func CannotConclude[N nameable.Nameable, E expr.Expression[N], T types.Type[N]](stat Status) Conclusion[N, E, T] {
	return Conclusion[N, E, T]{Status: stat}
}
