package token

type Token interface {
	SetLineChar(line, char int) Token
	GetLineChar() (line, char int)
	SetType(ty uint) Token
	GetType() uint
	SetValue(value string) (Token, error)
	GetValue() string
}