package compiler

import (
	"testing"
)

const (
	sample1 = `package main

func main() {
	println("hello world")
}
`
)

func TestCompilePackage(t *testing.T) {
	result, compilationError := Compile([]string{"../../../test/samples/sample-1"}, "", "", PrepareDefaultOptions())
	if compilationError != nil {
		t.Fatal(compilationError)
	}
	if len(result) == 0 {
		t.Fatal("CompilationError")
	}
}

func TestCompileString(t *testing.T) {
	result, compilationError := CompileString(sample1, "", ".", PrepareDefaultOptions())
	if compilationError != nil {
		t.Fatal(compilationError)
	}
	if len(result) == 0 {
		t.Fatal("CompilationError")
	}
}
