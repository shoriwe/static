package compile

import (
	"github.com/shoriwe/static/pkg/wasm/compiler"
	"io"
	"strings"
	"testing"
)

const (
	script = `package main
func main() {
	println("hello world")
}`
)

func TestCompileReaderToWASM(t *testing.T) {
	file, compilationError := compiler.CompileReader(strings.NewReader(script))
	if compilationError != nil {
		t.Fatal(compilationError)
	}
	defer file.Close()
	result, _ := io.ReadAll(file)
	if len(result) == 0 {
		t.Fatal("Nothing compiled")
	}
}

func TestCompilePackageToWASM(t *testing.T) {
	file, compilationError := compiler.Compile([]string{"./sample-1"})
	if compilationError != nil {
		t.Fatal(compilationError)
	}
	defer file.Close()
	result, _ := io.ReadAll(file)
	if len(result) == 0 {
		t.Fatal("Nothing compiled")
	}
}
