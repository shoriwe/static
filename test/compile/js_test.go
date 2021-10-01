package compile

import (
	"github.com/shoriwe/static/pkg/js/compiler"
	"io"
	"strings"
	"testing"
)

const (
	sample1 = `package main

func main() {
	println("hello world")
}
`
)

func TestCompilePackageToJS(t *testing.T) {
	file, compilationError := compiler.Compile([]string{"./sample-1"}, ".", compiler.PrepareDefaultOptions())
	if compilationError != nil {
		t.Fatal(compilationError)
	}
	defer file.Close()
	result, _ := io.ReadAll(file)
	if len(result) == 0 {
		t.Fatal("CompilationError")
	}
}

func TestCompileReaderToJS(t *testing.T) {
	file, compilationError := compiler.CompileReader(strings.NewReader(sample1), ".", compiler.PrepareDefaultOptions())
	if compilationError != nil {
		t.Fatal(compilationError)
	}
	defer file.Close()
	result, _ := io.ReadAll(file)
	if len(result) == 0 {
		t.Fatal("CompilationError")
	}
}
