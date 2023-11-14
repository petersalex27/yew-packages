package expr

import (
	"testing"

	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/util/testutil"
)

func TestBind(t *testing.T) {
	x := Var[nameable.Testable]("x")
	y := Var[nameable.Testable]("y")
	vars_x := BindersOnly[nameable.Testable]{x}
	//vars_x_y := BindersOnly[nameable.Testable]{x, y}
	idFunction := Function[nameable.Testable]{
		vars_x.Update(1),
		x.UpdateVars(-1, 1),
	}
	constFunction := Function[nameable.Testable]{
		vars_x.Update(1),
		y.UpdateVars(-1, 2),
	}

	tests := []struct{
		description string
		binders BindersOnly[nameable.Testable]
		bound Expression[nameable.Testable]
		expect Expression[nameable.Testable]
	}{
		{
			`Bind(x[0]).In(x[0]) == λx[1].x[1]`,
			vars_x, x,
			idFunction,
		},
		{
			`Bind(x[0]).In(y[0]) == λx[1].y[2]`,
			vars_x, y,
			constFunction,
		},
		{
			"Bind(x[0]).In(x[0] (λx[1].x[1])) == λx[1].(x[1] (λx[1].x[1]))",
			vars_x, Apply[nameable.Testable](x, idFunction),
			Function[nameable.Testable]{
				vars_x.Update(1),
				Apply[nameable.Testable](
					x.UpdateVars(-1, 1),
					Function[nameable.Testable]{
						vars_x.Update(1),
						x.UpdateVars(-1, 1),
					},
				),
			},
		},
		{
			"Bind(x[0]).In(y[0] (λx[1].x[1])) == λx[1].(y[2] (λx[1].x[1]))",
			vars_x, Apply[nameable.Testable](y, idFunction),
			Function[nameable.Testable]{
				vars_x.Update(1),
				Apply[nameable.Testable](
					y.UpdateVars(-1, 2),
					Function[nameable.Testable]{
						vars_x.Update(1),
						x.UpdateVars(-1, 1),
					},
				),
			},
		},
		{
			"Bind(x[0]).In(x[0] (λx[1].y[2])) == λx[1].(x[1] (λx[1].y[3]))",
			vars_x, Apply[nameable.Testable](x, constFunction),
			Function[nameable.Testable]{
				vars_x.Update(1),
				Apply[nameable.Testable](
					x.UpdateVars(-1, 1),
					Function[nameable.Testable]{
						vars_x.Update(1),
						y.UpdateVars(-1, 3),
					},
				),
			},
		},
	}

	for i, test := range tests {
		actual := Bind(test.binders...).In(test.bound)
		if !test.expect.StrictEquals(actual) {
			t.Fatal(testutil.
				Testing("equality", test.description).
				FailMessage(test.expect.StrictString(), actual.StrictString(), i),
			)
		}
	}
}