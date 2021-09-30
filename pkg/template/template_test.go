package template

import (
	"fmt"
	"strings"
	"testing"
)

const (
	rawHolderReplace    = "Hello (( Name ))!"
	rawGoCodeHolderEval = `<script>
{{

package main

func main() {
	println("Hello world")
}

}}
</script>`
	rawPlasmaCodeHolderEval = "Hello [[ print(\"Antonio\") ]]!"
)

func TestHolderReplace(t *testing.T) {
	rawReplace := map[string]string{
		"Name": "Antonio",
	}
	template := NewTemplate(strings.NewReader(rawHolderReplace), rawReplace)
	content, compileError := template.Compile()
	if compileError != nil {
		panic(compileError)
	}
	fmt.Println(string(content))
}

func TestGoCodeHolderReplace(t *testing.T) {
	rawReplace := map[string]string{}
	template := NewTemplate(strings.NewReader(rawGoCodeHolderEval), rawReplace)
	content, compileError := template.Compile()
	if compileError != nil {
		panic(compileError)
	}
	fmt.Println(string(content))
}

func TestPlasmaCodeHolderReplace(t *testing.T) {
	rawReplace := map[string]string{}
	template := NewTemplate(strings.NewReader(rawPlasmaCodeHolderEval), rawReplace)
	content, compileError := template.Compile()
	if compileError != nil {
		panic(compileError)
	}
	fmt.Println(string(content))
}
