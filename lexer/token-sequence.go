package lexer

import (
	"alex.peters/yew/token"
)

type TokenSequence struct {
	index uint32
	tokens []token.Token
}