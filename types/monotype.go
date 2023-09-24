package types

type Monotyped interface {
	Type
	DependentTyped
	ReplaceKindVar(replacing Variable, with Monotyped) Monotyped
	Replace(Variable, Monotyped) Monotyped
}

type Splitable interface {
	Monotyped
	Split() (name string, params []Monotyped)
}

func Name(m Monotyped) (name string, hasName bool) {
	var st Splitable
	st, hasName = m.(Splitable)
	if !hasName {
		name = m.String()
	} else {
		name, _ = st.Split()
	}
	return
}

func IsVariable(m Monotyped) bool {
	_, isVar := m.(Variable)
	return isVar
}

func Split(m Monotyped) (name string, params []Monotyped) {
	st, hasName := m.(Splitable)
	if !hasName {
		name, params = m.String(), nil
	} else {
		name, params = st.Split()
	}
	return
}