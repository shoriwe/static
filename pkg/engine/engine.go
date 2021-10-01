package engine

import (
	"errors"
	"fmt"
	"github.com/shoriwe/static/pkg/js/compiler"
	"github.com/shoriwe/static/pkg/template"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ContentGenerator func(engine *Engine) ([]byte, error)

type Script struct {
	WebPagePath string
	PackagePath string
}

type Scripts map[string]Script

type Asset struct {
	WebPagePath    string
	FileSystemPath string
}

type Assets map[string]Asset

func LoadScripts(location, output string) (Scripts, error) {
	listedResults, readError := os.ReadDir(location)
	if readError != nil {
		return nil, readError
	}
	files := map[string]Script{}
	for _, file := range listedResults {
		if !file.IsDir() {
			return nil, errors.New(fmt.Sprintf(ScriptIsNotAGoPackage, file.Name()))
		}
		packagePath := path.Join(location, file.Name())
		if packagePath[0] != '/' && packagePath[0] != '.' {
			packagePath = "./" + packagePath
		}
		files[file.Name()] = Script{
			WebPagePath: "/" + path.Join(output, file.Name()+".js"),
			PackagePath: packagePath,
		}
	}
	return files, nil
}

func LoadAssets(location, outputFolder string) (Assets, error) {
	files := map[string]Asset{}
	err := filepath.Walk(location,
		func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			relativePath := strings.ReplaceAll(path, location, "")
			for relativePath[0] == '/' {
				relativePath = relativePath[1:]
			}
			files[relativePath] = Asset{
				WebPagePath:    "/" + strings.ReplaceAll(path, location, outputFolder),
				FileSystemPath: path,
			}
			return nil
		},
	)
	return files, err
}

type Templates map[string]string

func LoadTemplates(location string) (Templates, error) {
	files := map[string]string{}
	err := filepath.Walk(location,
		func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			relativePath := strings.ReplaceAll(path, location, "")
			for relativePath[0] == '/' {
				relativePath = relativePath[1:]
			}
			files[relativePath] = path
			return nil
		},
	)
	return files, err
}

type Engine struct {
	Templates Templates
	Scripts   Scripts
	Assets    Assets
	paths     map[string]ContentGenerator
}

func NewEngine(templates Templates, scripts Scripts, assets Assets) *Engine {
	return &Engine{
		Templates: templates,
		Scripts:   scripts,
		Assets:    assets,
		paths:     map[string]ContentGenerator{},
	}
}

func (engine *Engine) RenderTemplate(template_ string, symbolMap map[string]string) ([]byte, error) {
	templatePath, ok := engine.Templates[template_]
	if !ok {
		return nil, errors.New(fmt.Sprintf(TemplateNotFound, template_))
	}
	file, openError := os.Open(templatePath)
	if openError != nil {
		return nil, openError
	}
	defer file.Close()
	return template.NewTemplate(file, symbolMap).Compile()
}

func (engine *Engine) AssetURL(asset string) (string, error) {
	a, ok := engine.Assets[asset]
	if !ok {
		return "", errors.New(fmt.Sprintf(AssetNotFound, asset))
	}
	return a.WebPagePath, nil
}

func (engine *Engine) ScriptURL(script string) (string, error) {
	a, ok := engine.Scripts[script]
	if !ok {
		return "", errors.New(fmt.Sprintf(ScriptNotFound, script))
	}
	return a.WebPagePath, nil
}

func (engine *Engine) HandlePath(newPath string, loader ContentGenerator) error {
	if _, ok := engine.paths[newPath]; ok {
		return errors.New(fmt.Sprintf(PathAlreadyInUse, newPath))
	}
	engine.paths[newPath] = loader
	return nil
}

func (engine *Engine) prepareStatic() error {
	for _, script := range engine.Scripts {
		engine.paths[script.WebPagePath] = func(engine *Engine) ([]byte, error) {
			file, compilationError := compiler.Compile([]string{script.PackagePath}, "", compiler.PrepareDefaultOptions())
			if compilationError != nil {
				return nil, compilationError
			}
			defer file.Close()
			return io.ReadAll(file)
		}
	}
	for _, asset := range engine.Assets {
		engine.paths[asset.WebPagePath] = func(engine *Engine) ([]byte, error) {
			file, openError := os.Open(asset.FileSystemPath)
			if openError != nil {
				return nil, openError
			}
			defer file.Close()
			return io.ReadAll(file)
		}
	}
	return nil
}

func createFile(filePath string) (*os.File, error) {
	directory, _ := filepath.Split(filePath)
	directoryCreationError := os.MkdirAll(directory, os.ModePerm)
	if directoryCreationError != nil {
		return nil, directoryCreationError
	}
	return os.Create(filePath)
}

func (engine *Engine) Generate(output string) error {
	preparationError := engine.prepareStatic()
	if preparationError != nil {
		return preparationError
	}
	for webPath, contentGenerator := range engine.paths {
		file, creationError := createFile(path.Join(output, webPath))
		if creationError != nil {
			return creationError
		}
		content, contentGenerationError := contentGenerator(engine)
		if contentGenerationError != nil {
			return contentGenerationError
		}
		_, writeError := file.Write(content)
		if writeError != nil {
			return writeError
		}
		closeError := file.Close()
		if closeError != nil {
			return closeError
		}
	}
	return nil
}
