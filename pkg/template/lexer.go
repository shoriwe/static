package template

import (
	"errors"
	"fmt"
	"text/scanner"
)

type Lexer struct {
	scanner     *scanner.Scanner
	hasNext     bool
	currentRune rune
}

func (lexer *Lexer) HasNext() bool {
	return lexer.currentRune != scanner.EOF
}

func (lexer *Lexer) next() rune {
	lexer.currentRune = lexer.scanner.Next()
	if lexer.currentRune == Return {
		return lexer.next()
	}
	return lexer.currentRune
}

func (lexer *Lexer) tokenizeHolder(startPosition scanner.Position, resultId uint, openRune, closeRune rune) (*Token, error) {
	if !lexer.HasNext() {
		return &Token{
			DirectValue: Raw,
			String:      string([]rune{lexer.currentRune}),
			Position:    lexer.scanner.Pos(),
			Length:      1,
		}, nil
	}
	if lexer.next() != openRune {
		return lexer.tokenizeRaw([]rune{openRune, lexer.currentRune}, startPosition)
	}
	if !lexer.HasNext() {
		return nil, errors.New(fmt.Sprintf(HolderNeverClosed, startPosition.Line, startPosition.Offset))
	}
	switch lexer.next() {
	case openRune:
		return &Token{
			DirectValue: Raw,
			String:      string([]rune{openRune, openRune}),
			Position:    startPosition,
			Length:      2,
		}, nil
	case WhiteSpace, NewLine, Tab:
		break
	default:
		return nil, errors.New(fmt.Sprintf(InvalidHolderDefinition, startPosition.Line, startPosition.Offset))
	}
	startPosition = lexer.scanner.Pos()
	var body []rune
	// Parse the body of the holder
	for {
		if !lexer.HasNext() {
			return nil, errors.New(fmt.Sprintf(HolderNeverClosed, startPosition.Line, startPosition.Offset))
		}
		switch lexer.next() {
		case WhiteSpace, NewLine, Tab:
			characterFound := lexer.currentRune
			if !lexer.HasNext() {
				return nil, errors.New(fmt.Sprintf(InvalidHolderDefinition, startPosition.Line, startPosition.Offset))
			}
			if lexer.scanner.Peek() != closeRune {
				body = append(body, characterFound)
				break
			}
			lexer.next()
			if lexer.next() != closeRune {
				body = append(body, characterFound, closeRune, lexer.currentRune)
				break
			}
			return &Token{
				DirectValue: resultId,
				String:      string(body),
				Position:    startPosition,
				Length:      len(body),
			}, nil
		default:
			body = append(body, lexer.currentRune)
		}
	}
}

func (lexer *Lexer) tokenizeRaw(previous []rune, startPosition scanner.Position) (*Token, error) {
	body := previous
rawLoop:
	for lexer.HasNext() {
		switch lexer.scanner.Peek() {
		case scanner.EOF, OpenSquare, OpenBrace, OpenParentheses, CloseSquare, CloseBrace, CloseParentheses:
			break rawLoop
		}
		body = append(body, lexer.next())
	}
	return &Token{
		DirectValue: Raw,
		String:      string(body),
		Position:    startPosition,
		Length:      len(body),
	}, nil
}

func (lexer *Lexer) Next() (*Token, error) {
	if !lexer.HasNext() {
		return &Token{
			DirectValue: EOF,
			String:      "EOF",
			Position:    lexer.scanner.Pos(),
			Length:      1,
		}, nil
	}
	position := lexer.scanner.Pos()
	switch lexer.next() {
	case scanner.EOF:
		return &Token{
			DirectValue: EOF,
			String:      "EOF",
			Position:    lexer.scanner.Pos(),
			Length:      1,
		}, nil
	case OpenBrace:
		return lexer.tokenizeHolder(position, GoCodeHolder, OpenBrace, CloseBrace)
	case OpenSquare:
		return lexer.tokenizeHolder(position, PlasmaCodeHolder, OpenSquare, CloseSquare)
	case OpenParentheses:
		return lexer.tokenizeHolder(position, Holder, OpenParentheses, CloseParentheses)
	case CloseBrace, CloseSquare, CloseParentheses:
		closing := lexer.currentRune
		if !lexer.HasNext() {
			return &Token{
				DirectValue: Raw,
				String:      string([]rune{lexer.currentRune}),
				Position:    position,
				Length:      1,
			}, nil
		}
		if lexer.next() == closing {
			if !lexer.HasNext() {
				return nil, errors.New("1")
			} else if lexer.next() != closing {
				return nil, errors.New(string(lexer.currentRune))
			}
			return &Token{
				DirectValue: Raw,
				String:      string([]rune{closing, closing}),
				Position:    position,
				Length:      2,
			}, nil
		} else if !lexer.HasNext() {
			return &Token{
				DirectValue: Raw,
				String:      string([]rune{closing}),
				Position:    position,
				Length:      1,
			}, nil
		}
		return lexer.tokenizeRaw([]rune{closing, lexer.currentRune}, position)
	default:
		return lexer.tokenizeRaw([]rune{lexer.currentRune}, position)
	}
}

func NewLexer(scanner *scanner.Scanner) *Lexer {
	return &Lexer{
		scanner: scanner,
		hasNext: true,
	}
}
