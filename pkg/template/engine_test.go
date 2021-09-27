package template

import (
	"fmt"
	"strings"
	"testing"
)

const (
	rawHolderReplace  = "Hello [[ Name ]]"
	rawCodeHolderEval = "Hello {{ println(\"Antonio\") }}"
)

func TestHolderReplace(t *testing.T) {
	rawReplace := map[string]string{
		"Name": "Antonio",
	}
	codeReplace := map[string]string{}
	asset := NewAsset(strings.NewReader(rawHolderReplace), rawReplace, codeReplace)
	content, compileError := asset.Compile()
	if compileError != nil {
		panic(compileError)
	}
	fmt.Println(string(content))
}

func TestCodeHolderReplace(t *testing.T) {
	rawReplace := map[string]string{}
	codeReplace := map[string]string{}
	asset := NewAsset(strings.NewReader(rawCodeHolderEval), rawReplace, codeReplace)
	content, compileError := asset.Compile()
	if compileError != nil {
		panic(compileError)
	}
	fmt.Println(string(content))
}
