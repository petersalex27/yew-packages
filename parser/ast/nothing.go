package ast

type Nothing struct{}

func (Nothing) NodeType() Type {
	return None
}
