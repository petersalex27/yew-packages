package lexer

import (
	"github.com/petersalex27/yew-packages/token"
)

type TokenSequence struct {
	index uint32
	tokens []token.Token
}