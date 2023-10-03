package expr

type test_named string

func (t test_named) GetName() string {
	return string(t)
}