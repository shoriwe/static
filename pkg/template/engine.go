package template

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/shoriwe/gplasma"
	"io"
	"text/scanner"
)

type Asset struct {
	lexer         *Lexer
	RawReplace    map[string]string
	PlasmaReplace map[string]string
}

func (asset *Asset) Compile() ([]byte, error) {
	bufferHandler := bytes.NewBuffer([]byte{})
	vm := gplasma.NewVirtualMachine()
	for asset.lexer.HasNext() {
		token, lexingError := asset.lexer.Next()
		if lexingError != nil {
			return nil, errors.New("During compilation -> " + lexingError.Error())
		}
		var writeError error = nil
		switch token.DirectValue {
		case Raw:
			_, writeError = bufferHandler.WriteString(token.String)
		case Holder:
			replace, found := asset.RawReplace[token.String]
			if !found {
				return nil, errors.New(fmt.Sprintf(HolderNotFound, token.String))
			}
			_, writeError = bufferHandler.WriteString(replace)
		case CodeHolder:
			// Compile plasma and execute it
			code, found := asset.PlasmaReplace[token.String]
			if !found {
				code = token.String
			}
			output := bytes.NewBuffer([]byte{})
			vm.Stdout = output
			vm.Stderr = vm.Stdout
			executionError, success := vm.ExecuteMain(code)
			if !success {
				return nil, errors.New(executionError.GetClass(vm.Plasma).Name + ": " + executionError.String)
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

func NewAsset(reader io.Reader, rawReplace map[string]string, plasmaReplace map[string]string) *Asset {
	assetScanner := new(scanner.Scanner)
	assetScanner.Init(reader)
	return &Asset{
		lexer:         NewLexer(assetScanner),
		RawReplace:    rawReplace,
		PlasmaReplace: plasmaReplace,
	}
}