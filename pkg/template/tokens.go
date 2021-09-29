package template

import "text/scanner"

const (
	EOF = iota
	Raw
	Holder
	PlasmaCodeHolder
	GoCodeHolder
)

const (
	OpenParentheses  rune = '('
	CloseParentheses rune = ')'
	OpenBrace        rune = '{'
	CloseBrace       rune = '}'
	OpenSquare       rune = '['
	CloseSquare      rune = ']'
	WhiteSpace       rune = ' '
	NewLine          rune = '\n'
)

type Token struct {
	DirectValue uint
	String      string
	Position    scanner.Position
	Length      int
}
