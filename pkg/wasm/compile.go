package js

import (
	"io"
	"os"
)

func CompilePackageToScript(goCode io.Reader, symbolsToExport []string) ([]byte, error) {
	tempFile, tempFileCreationError := os.CreateTemp("", "")
	if tempFileCreationError != nil {
		return nil, tempFileCreationError
	}
	defer os.Remove(tempFile.Name())
	_, copyError := io.Copy(tempFile, goCode)
	if copyError != nil {
		return nil, copyError
	}
	_, seekError := tempFile.Seek(0, io.SeekStart)
	if seekError != nil {
		return nil, seekError
	}
	return nil, nil
}
