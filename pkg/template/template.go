package template

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/shoriwe/gplasma"
	"github.com/shoriwe/static/pkg/js/compiler"
	"io"
	"strings"
	"text/scanner"
)

type Template struct {
	lexer      *Lexer
	RawReplace map[string]string
	vm         *gplasma.VirtualMachine
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
		case GoCodeHolder:
			file, compilationError := compiler.CompileReader(strings.NewReader(token.String), "", compiler.PrepareDefaultOptions())
			if compilationError != nil {
				return nil, compilationError
			}
			compiledCode, readError := io.ReadAll(file)
			if readError != nil {
				return nil, readError
			}
			_, writeError = bufferHandler.Write(compiledCode)
			_ = file.Close()
		case PlasmaCodeHolder:
			// Compile plasma and execute it
			output := bytes.NewBuffer([]byte{})
			template.vm.Stdout = output
			template.vm.Stderr = template.vm.Stdout
			executionError, success := template.vm.ExecuteMain(token.String)
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

func NewTemplate(reader io.Reader, rawReplace map[string]string) *Template {
	templateScanner := new(scanner.Scanner)
	templateScanner.Init(reader)
	return &Template{
		lexer:      NewLexer(templateScanner),
		RawReplace: rawReplace,
		vm:         gplasma.NewVirtualMachine(),
	}
}

func NewTemplateCustomVM(reader io.Reader, rawReplace map[string]string, plasmaReplace map[string]string, vm *gplasma.VirtualMachine) *Template {
	templateScanner := new(scanner.Scanner)
	templateScanner.Init(reader)
	return &Template{
		lexer:      NewLexer(templateScanner),
		RawReplace: rawReplace,
		vm:         vm,
	}
}
