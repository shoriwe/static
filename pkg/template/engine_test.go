package template

import (
	"fmt"
	"strings"
	"testing"
)

const (
	rawHolderReplace  = "Hello [[ Name ]]!"
	rawCodeHolderEval = "Hello {{ print(\"Antonio\") }}!"
)

func TestHolderReplace(t *testing.T) {
	rawReplace := map[string]string{
		"Name": "Antonio",
	}
	codeReplace := map[string]string{}
	template := NewTemplate(strings.NewReader(rawHolderReplace), rawReplace, codeReplace)
	content, compileError := template.Compile()
	if compileError != nil {
		panic(compileError)
	}
	fmt.Println(string(content))
}

func TestCodeHolderReplace(t *testing.T) {
	rawReplace := map[string]string{}
	codeReplace := map[string]string{}
	template := NewTemplate(strings.NewReader(rawCodeHolderEval), rawReplace, codeReplace)
	content, compileError := template.Compile()
	if compileError != nil {
		panic(compileError)
	}
	fmt.Println(string(content))
}
