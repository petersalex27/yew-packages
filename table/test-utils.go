package table

type test_nameable string

func (n test_nameable) GetName() string { return string(n) }