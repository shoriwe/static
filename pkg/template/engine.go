package template

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/shoriwe/gplasma"
	"io"
	"text/scanner"
)

type Template struct {
	lexer         *Lexer
	RawReplace    map[string]string
	PlasmaReplace map[string]string
	vm            *gplasma.VirtualMachine
}

func (template *Template) Compile() ([]byte, error) {
	bufferHandler := bytes.NewBuffer([]byte{})
	for template.lexer.HasNext() {
		token, lexingError := template.lexer.Next()
		if lexingError != nil {
			return nil, errors.New("During compilation -> " + lexingError.Error())
		}
		var writeError error = nil
		switch token.DirectValue {
		case Raw:
			_, writeError = bufferHandler.WriteString(token.String)
		case Holder:
			replace, found := template.RawReplace[token.String]
			if !found {
				return nil, errors.New(fmt.Sprintf(HolderNotFound, token.String))
			}
			_, writeError = bufferHandler.WriteString(replace)
		case CodeHolder:
			// Compile plasma and execute it
			code, found := template.PlasmaReplace[token.String]
			if !found {
				code = token.String
			}
			output := bytes.NewBuffer([]byte{})
			template.vm.Stdout = output
			template.vm.Stderr = template.vm.Stdout
			executionError, success := template.vm.ExecuteMain(code)
			if !success {
				return nil, errors.New(executionError.GetClass(template.vm.Plasma).Name + ": " + executionError.String)
			}
			_, writeError = bufferHandler.Write(output.Bytes())
		case EOF:
			break
		}
		if writeError != nil {
			return nil, writeError
		}
	}
	return bufferHandler.Bytes(), nil
}

func NewTemplate(reader io.Reader, rawReplace map[string]string, plasmaReplace map[string]string) *Template {
	templateScanner := new(scanner.Scanner)
	templateScanner.Init(reader)
	return &Template{
		lexer:         NewLexer(templateScanner),
		RawReplace:    rawReplace,
		PlasmaReplace: plasmaReplace,
		vm:            gplasma.NewVirtualMachine(),
	}
}

func NewTemplateCustomVM(reader io.Reader, rawReplace map[string]string, plasmaReplace map[string]string, vm *gplasma.VirtualMachine) *Template {
	templateScanner := new(scanner.Scanner)
	templateScanner.Init(reader)
	return &Template{
		lexer:         NewLexer(templateScanner),
		RawReplace:    rawReplace,
		PlasmaReplace: plasmaReplace,
		vm:            vm,
	}
}
