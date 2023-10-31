package parser

type ActionFunction func(any)

type Action struct {
	Name string
	Does func(any)
}

type actionMap map[string]ActionFunction

type actionRequester struct {
	failGet func(name string) func(any)
	actionMap
}

func defaultFailGet(name string) func(any) { 
	return func (a any) {
		panic("tried to load action `" + name + "`, but it does not exist")
	}
}

// returns action named `name`, panics if name does not map to an action
func (ar *actionRequester) get(name string) func(any) {
	action, found := ar.actionMap[name]
	if !found {
		return ar.failGet(name)
	}
	return action
}

// panics on duplicate map attempt
func (ar *actionRequester) attach(action Action) {
	if _, found := ar.actionMap[action.Name]; found {
		panic("tried to attach action `" + action.Name + "` more than once")
	}
	ar.actionMap[action.Name] = action.Does
}

// expandActions adds arguments to existing actions in actionMap. Function
// panics if an attempt to map an action name to an action is made
func (p *parser) expandActions(actions ...Action) *parser {
	// store old actions to be re-added at end of function
	old := p.actions

	// make actionMap that can hold both new and old actions
	p.actions = actionRequester{
		old.failGet,
		make(actionMap, len(actions)+len(old.actionMap)),
	}

	// attach actions from arguments
	for _, act := range actions {
		p.actions.attach(act)
	}

	// re-attach old actions
	for name, does := range old.actionMap {
		p.actions.attach(Action{name, does})
	}

	return p
}

// Attach adds actions to parser to be used with action reductions. Function
// panics if an attempt to map an action name to an action is made
func (p *parser) Attach(actions ...Action) *parser {
	if len(p.actions.actionMap) != 0 {
		return p.expandActions(actions...)
	}

	// map w/ enough room for len(actions)
	p.actions = actionRequester{
		defaultFailGet,
		make(actionMap, len(actions)),
	}

	// map names to actions
	for _, act := range actions {
		p.actions.attach(act)
	}

	return p
}

// updates (overwrites) action if it exists in action map, else adds it to map
func (p *parser) UpdateAction(action Action) {
	if len(p.actions.actionMap) == 0 {
		p.Attach(action) // makes new actionRequest object
		return
	}

	p.actions.actionMap[action.Name] = action.Does
}

func (p *parser) SetActionRequestFail(onFailGet func(name string) func(any)) {
	p.actions.failGet = onFailGet
}