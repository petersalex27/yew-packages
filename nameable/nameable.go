package nameable

// something with a name
type Nameable interface {
	GetName() string
}

// basic implementation for ease of testing w/ things that require Nameable
// implementations
type Testable string

// wraps `name`
func MakeTestable(name string) Testable {
	return Testable(name)
}

// returns wrapped string
func (t Testable) GetName() string {
	return string(t)
}