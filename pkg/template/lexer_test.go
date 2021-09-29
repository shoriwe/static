package template

import (
	"fmt"
	"strings"
	"testing"
	"text/scanner"
)

const (
	test1 = "Hello world"

	test2 = "{{ 1 + 2 }}"
	test3 = "{{ 1 + 2 }} Message"
	test4 = "Message {{ 1 + 2 }}"
	test5 = "{{{ hello world }}}"
	test6 = "{ Hello World! }"

	test7  = "[[ 1 + 2 ]]"
	test8  = "[[ 1 + 2 ]] Message"
	test9  = "Message [[ 1 + 2 ]]"
	test10 = "[[[ hello world ]]]"
	test11 = "[ Hello World! ]"

	test12 = "(( 1 + 2 ))"
	test13 = "(( 1 + 2 )) Message"
	test14 = "Message (( 1 + 2 ))"
	test15 = "((( hello world )))"
	test16 = "( Hello World! )"
)

func TestRawOnly(t *testing.T) {
	var testScanner scanner.Scanner
	testScanner.Init(strings.NewReader(test1))
	lexer := NewLexer(&testScanner)
	for lexer.HasNext() {
		token, err := lexer.Next()
		if err != nil {
			panic(err)
		}
		fmt.Println(token)
	}
}

func TestBraceOnly(t *testing.T) {
	var testScanner scanner.Scanner
	testScanner.Init(strings.NewReader(test2))
	lexer := NewLexer(&testScanner)
	for lexer.HasNext() {
		token, err := lexer.Next()
		if err != nil {
			panic(err)
		}
		fmt.Println(token)
	}
}

func TestBraceThenRaw(t *testing.T) {
	var testScanner scanner.Scanner
	testScanner.Init(strings.NewReader(test3))
	lexer := NewLexer(&testScanner)
	for lexer.HasNext() {
		token, err := lexer.Next()
		if err != nil {
			panic(err)
		}
		fmt.Println(token)
	}
}

func TestRawThenBrace(t *testing.T) {
	var testScanner scanner.Scanner
	testScanner.Init(strings.NewReader(test4))
	lexer := NewLexer(&testScanner)
	for lexer.HasNext() {
		token, err := lexer.Next()
		if err != nil {
			panic(err)
		}
		fmt.Println(token)
	}
}

func TestEscapedBrace(t *testing.T) {
	var testScanner scanner.Scanner
	testScanner.Init(strings.NewReader(test5))
	lexer := NewLexer(&testScanner)
	for lexer.HasNext() {
		token, err := lexer.Next()
		if err != nil {
			panic(err)
		}
		fmt.Println(token)
	}
}

func TestRawOpenBrace(t *testing.T) {
	var testScanner scanner.Scanner
	testScanner.Init(strings.NewReader(test6))
	lexer := NewLexer(&testScanner)
	for lexer.HasNext() {
		token, err := lexer.Next()
		if err != nil {
			panic(err)
		}
		fmt.Println(token)
	}
}

func TestSquareOnly(t *testing.T) {
	var testScanner scanner.Scanner
	testScanner.Init(strings.NewReader(test7))
	lexer := NewLexer(&testScanner)
	for lexer.HasNext() {
		token, err := lexer.Next()
		if err != nil {
			panic(err)
		}
		fmt.Println(token)
	}
}

func TestSquareThenRaw(t *testing.T) {
	var testScanner scanner.Scanner
	testScanner.Init(strings.NewReader(test8))
	lexer := NewLexer(&testScanner)
	for lexer.HasNext() {
		token, err := lexer.Next()
		if err != nil {
			panic(err)
		}
		fmt.Println(token)
	}
}

func TestRawThenSquare(t *testing.T) {
	var testScanner scanner.Scanner
	testScanner.Init(strings.NewReader(test9))
	lexer := NewLexer(&testScanner)
	for lexer.HasNext() {
		token, err := lexer.Next()
		if err != nil {
			panic(err)
		}
		fmt.Println(token)
	}
}

func TestEscapedSquare(t *testing.T) {
	var testScanner scanner.Scanner
	testScanner.Init(strings.NewReader(test10))
	lexer := NewLexer(&testScanner)
	for lexer.HasNext() {
		token, err := lexer.Next()
		if err != nil {
			panic(err)
		}
		fmt.Println(token)
	}
}

func TestRawOpenSquare(t *testing.T) {
	var testScanner scanner.Scanner
	testScanner.Init(strings.NewReader(test11))
	lexer := NewLexer(&testScanner)
	for lexer.HasNext() {
		token, err := lexer.Next()
		if err != nil {
			panic(err)
		}
		fmt.Println(token)
	}
}

func TestParenthesesOnly(t *testing.T) {
	var testScanner scanner.Scanner
	testScanner.Init(strings.NewReader(test12))
	lexer := NewLexer(&testScanner)
	for lexer.HasNext() {
		token, err := lexer.Next()
		if err != nil {
			panic(err)
		}
		fmt.Println(token)
	}
}

func TestParenthesesThenRaw(t *testing.T) {
	var testScanner scanner.Scanner
	testScanner.Init(strings.NewReader(test13))
	lexer := NewLexer(&testScanner)
	for lexer.HasNext() {
		token, err := lexer.Next()
		if err != nil {
			panic(err)
		}
		fmt.Println(token)
	}
}

func TestRawThenParentheses(t *testing.T) {
	var testScanner scanner.Scanner
	testScanner.Init(strings.NewReader(test14))
	lexer := NewLexer(&testScanner)
	for lexer.HasNext() {
		token, err := lexer.Next()
		if err != nil {
			panic(err)
		}
		fmt.Println(token)
	}
}

func TestEscapedParentheses(t *testing.T) {
	var testScanner scanner.Scanner
	testScanner.Init(strings.NewReader(test15))
	lexer := NewLexer(&testScanner)
	for lexer.HasNext() {
		token, err := lexer.Next()
		if err != nil {
			panic(err)
		}
		fmt.Println(token)
	}
}

func TestRawOpenParentheses(t *testing.T) {
	var testScanner scanner.Scanner
	testScanner.Init(strings.NewReader(test16))
	lexer := NewLexer(&testScanner)
	for lexer.HasNext() {
		token, err := lexer.Next()
		if err != nil {
			panic(err)
		}
		fmt.Println(token)
	}
}
