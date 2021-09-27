package template

import (
	"fmt"
	"strings"
	"testing"
	"text/scanner"
)

const (
	test1  = "Hello world"
	test2  = "{{ 1 + 2 }}"
	test3  = "{{ 1 + 2 }} Message"
	test4  = "Message {{ 1 + 2 }}"
	test5  = "{{{ hello world }}}"
	test6  = "{ Hello World! }"
	test7  = "[[ 1 + 2 ]]"
	test8  = "[[ 1 + 2 ]] Message"
	test9  = "Message [[ 1 + 2 ]]"
	test10 = "[[[ hello world ]]]"
	test11 = "[ Hello World! ]"
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

func TestCodeHolderOnly(t *testing.T) {
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

func TestCodeHolderThenRaw(t *testing.T) {
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

func TestRawThenCodeHolder(t *testing.T) {
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

func TestEscapedCodeHolder(t *testing.T) {
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

func TestHolderOnly(t *testing.T) {
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

func TestHolderThenRaw(t *testing.T) {
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

func TestRawThenHolder(t *testing.T) {
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

func TestEscapedHolder(t *testing.T) {
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
