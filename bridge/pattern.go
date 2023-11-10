package bridge

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/types"
)

// representation of a value, but not a value--think of as a layout for a value
type Pattern[N nameable.Nameable] struct {
	ty types.Monotyped[N]
	pattern expr.AlmostPattern[N]
}

// returns `p, true`
func ToPattern[N nameable.Nameable](p expr.AlmostPattern[N], ty types.Monotyped[N]) (Pattern[N], bool) { 
	return Pattern[N]{ty, p}, true 
}


// Match tries to match two patterns. Match returns true on success, else false
func (this Pattern[N]) Match(that Pattern[N]) bool {
	return this.ty.Equals(that.ty) && this.pattern.Match(that.pattern)
}