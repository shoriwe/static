package compiler

import (
	"testing"
)

const (
	script = `package main
func main() {
	println("hello world")
}`
)

func TestCompileString(t *testing.T) {
	result, compilationError := CompileString(script)
	if compilationError != nil {
		t.Fatal(compilationError)
	}
	if len(result) == 0 {
		t.Fatal("Nothing compiled")
	}
}

func TestCompilePackage(t *testing.T) {
	result, compilationError := Compile([]string{"../../../samples/sample-1"})
	if compilationError != nil {
		t.Fatal(compilationError)
	}
	if len(result) == 0 {
		t.Fatal("Nothing compiled")
	}
}
