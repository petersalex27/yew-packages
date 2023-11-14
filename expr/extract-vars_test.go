package expr

import (
	"testing"

	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/util/testutil"
)

func TestExtractVars(t *testing.T) {
	x := Var[nameable.Testable]("x")
	y := Var[nameable.Testable]("y")
	z := Var[nameable.Testable]("z")
	vars_x := BindersOnly[nameable.Testable]{x}
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
		input Expression[nameable.Testable]
		expect BindersOnly[nameable.Testable]
	}{
		{
			`λx[1].x[1] => []`,
			idFunction,
			BindersOnly[nameable.Testable]{},
		},
		{
			`λx[1].y[2] => [ y[1] ]`,
			constFunction,
			BindersOnly[nameable.Testable]{y.UpdateVars(-1, 1).(Variable[nameable.Testable])},
		},
		{
			"λx[1].(x[1] (λx[1].x[1])) => []",
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
			BindersOnly[nameable.Testable]{},
		},
		{
			"λx[1].(y[2] (λx[1].x[1])) => [ y[1] ]",
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
			BindersOnly[nameable.Testable]{y.UpdateVars(-1, 1).(Variable[nameable.Testable])},
		},
		{
			"λx[1].(x[1] (λx[1].y[3])) => [ y[1] ]",
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
			BindersOnly[nameable.Testable]{y.UpdateVars(-1, 1).(Variable[nameable.Testable])},
		},
		{
			"λx[1].(x[1] (λx[1] . x[1] y[3] z[4])) => [ y[1], z[2] ]",
			Function[nameable.Testable]{
				vars_x.Update(1), // λx[1]
				Apply[nameable.Testable](
					x.UpdateVars(-1, 1), // x[1]
					Function[nameable.Testable]{
						vars_x.Update(1), // λx[1]
						Apply(
							x.UpdateVars(-1, 1), // x[1]
							y.UpdateVars(-1, 3), // y[3] 
							z.UpdateVars(-1, 4), // z[4]
						),
					},
				),
			},
			BindersOnly[nameable.Testable]{
				y.UpdateVars(-1, 1).(Variable[nameable.Testable]),
				z.UpdateVars(-1, 2).(Variable[nameable.Testable]),
			},
		},
	}

	for i, test := range tests {
		actual := BindersOnly[nameable.Testable](test.input.ExtractVariables(0))

		if !test.expect.StrictEquals(actual) {
			t.Fatal(testutil.
				Testing("equality", test.description).
				FailMessage(test.expect.StrictString(), actual.StrictString(), i),
			)
		}
	}
}