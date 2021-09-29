package compiler

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path"
	"runtime"
)

func CompileString(sourceCode string) ([]byte, error) {
	tempFile, tempFileCreationError := os.CreateTemp("", "*.go")
	if tempFileCreationError != nil {
		return nil, tempFileCreationError
	}
	_, writeError := tempFile.WriteString(sourceCode)
	if writeError != nil {
		return nil, writeError
	}
	_ = tempFile.Close()
	defer os.Remove(tempFile.Name())
	return Compile([]string{tempFile.Name()})
}

func Compile(files []string) ([]byte, error) {
	tempFile, tempFileCreationError := os.CreateTemp("", "*.wasm")
	if tempFileCreationError != nil {
		return nil, tempFileCreationError
	}
	_ = tempFile.Close()
	defer os.Remove(tempFile.Name())
	extension := ""
	if os.PathSeparator == '\\' {
		extension = ".exe"
	}
	cmd := exec.Command(path.Join(runtime.GOROOT(), "bin", "go"+extension), append([]string{"build", "-o", tempFile.Name()}, files...)...)
	output, executionError := cmd.CombinedOutput()
	if executionError != nil {
		return nil, errors.New(executionError.Error() + "\n" + string(output))
	}
	file, openError := os.Open(tempFile.Name())
	if openError != nil {
		return nil, openError
	}
	defer file.Close()
	return io.ReadAll(file)
}
