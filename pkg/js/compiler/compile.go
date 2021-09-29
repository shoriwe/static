package compiler

import (
	"fmt"
	gbuild "github.com/gopherjs/gopherjs/build"
	"go/build"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func PrepareDefaultOptions() *gbuild.Options {
	return &gbuild.Options{
		GOROOT:         "",
		GOPATH:         "",
		Verbose:        false,
		Quiet:          false,
		Watch:          false,
		CreateMapFile:  true,
		MapToLocalDisk: false,
		Minify:         false,
		Color:          true,
		BuildTags:      nil,
	}
}

func CompileString(fileContent string, tags, packagePath string, options *gbuild.Options) ([]byte, error) {
	tempCodeFile, tempCodeFileCreationError := os.CreateTemp("", "code-*.go")
	if tempCodeFileCreationError != nil {
		return nil, tempCodeFileCreationError
	}
	defer os.Remove(tempCodeFile.Name())
	_, writeError := tempCodeFile.WriteString(fileContent)
	if writeError != nil {
		return nil, writeError
	}
	_ = tempCodeFile.Close()
	return Compile([]string{tempCodeFile.Name()}, tags, packagePath, options)
}

func Compile(args []string, tags, currentDirectory string, options *gbuild.Options) ([]byte, error) {
	tempOutputFile, tempOutputFileCreationError := os.CreateTemp("", "output-*.js")
	if tempOutputFileCreationError != nil {
		return nil, tempOutputFileCreationError
	}
	defer os.Remove(tempOutputFile.Name())
	_ = tempOutputFile.Close()

	pkgObj := tempOutputFile.Name()

	options.BuildTags = strings.Fields(tags)
	for {
		s, err := gbuild.NewSession(options)
		if err != nil {
			return nil, err
		}

		err = func() error {
			// Handle "gopherjs build [files]" ad-hoc package mode.
			if len(args) > 0 && (strings.HasSuffix(args[0], ".go") || strings.HasSuffix(args[0], ".inc.js")) {
				for _, arg := range args {
					if !strings.HasSuffix(arg, ".go") && !strings.HasSuffix(arg, ".inc.js") {
						return fmt.Errorf("named files must be .go or .inc.js files")
					}
				}
				if pkgObj == "" {
					basename := filepath.Base(args[0])
					pkgObj = basename[:len(basename)-3] + ".js"
				}
				names := make([]string, len(args))
				for i, name := range args {
					name = filepath.ToSlash(name)
					names[i] = name
					if s.Watcher != nil {
						s.Watcher.Add(name)
					}
				}
				err := s.BuildFiles(args, pkgObj, currentDirectory)
				return err
			}

			xctx := gbuild.NewBuildContext(s.InstallSuffix(), options.BuildTags)
			// Expand import path patterns.
			pkgs, err := xctx.Match(args)
			if err != nil {
				return fmt.Errorf("failed to expand patterns %v: %w", args, err)
			}

			for _, pkgPath := range pkgs {
				if s.Watcher != nil {
					pkg, err := xctx.Import(pkgPath, "", build.FindOnly)
					if err != nil {
						return err
					}
					s.Watcher.Add(pkg.Dir)
				}
				pkg, err := xctx.Import(pkgPath, ".", 0)
				if err != nil {
					return err
				}
				archive, err := s.BuildPackage(pkg)
				if err != nil {
					return err
				}
				if len(pkgs) == 1 { // Only consider writing output if single package specified.
					if pkgObj == "" {
						pkgObj = filepath.Base(pkg.Dir) + ".js"
					}
					if pkg.IsCommand() && !pkg.UpToDate {
						if err := s.WriteCommandPackage(archive, pkgObj); err != nil {
							return err
						}
					}
				}
			}
			return nil
		}()
		if err != nil {
			return nil, err
		}

		output, openError := os.Open(tempOutputFile.Name())
		if openError != nil {
			return nil, openError
		}
		defer output.Close()
		return io.ReadAll(output)
	}
}
