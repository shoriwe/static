package engine

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/shoriwe/static/pkg/js/compiler"
	"github.com/shoriwe/static/pkg/template"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	JSExtension   = "js"
	JsonExtension = "json"
	HTMLExtension = "html"
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
	m         *minify.M
}

func NewEngine(templates Templates, scripts Scripts, assets Assets) *Engine {
	m := minify.New()
	m.Add("application/javascript", js.DefaultMinifier)
	m.Add("text/json", json.DefaultMinifier)
	m.Add("text/html", html.DefaultMinifier)
	return &Engine{
		Templates: templates,
		Scripts:   scripts,
		Assets:    assets,
		paths:     map[string]ContentGenerator{},
		m:         m,
	}
}

func (engine *Engine) MinifyJS(content []byte) ([]byte, error) {
	output := bytes.NewBuffer([]byte{})
	minifyError := engine.m.Minify("application/javascript", output, bytes.NewReader(content))
	return output.Bytes(), minifyError
}

func (engine *Engine) MinifyJson(content []byte) ([]byte, error) {
	output := bytes.NewBuffer([]byte{})
	minifyError := engine.m.Minify("text/json", output, bytes.NewReader(content))
	return output.Bytes(), minifyError
}

func (engine *Engine) MinifyHTML(content []byte) ([]byte, error) {
	output := bytes.NewBuffer([]byte{})
	minifyError := engine.m.Minify("text/html", output, bytes.NewReader(content))
	return output.Bytes(), minifyError
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

func newAssetLoader(asset Asset) ContentGenerator {
	return func(e *Engine) ([]byte, error) {
		file, openError := os.Open(asset.FileSystemPath)
		if openError != nil {
			return nil, openError
		}
		defer file.Close()
		content, readError := io.ReadAll(file)
		if readError != nil {
			return nil, readError
		}
		split := strings.Split(file.Name(), ".")
		if len(split) == 1 {
			return content, nil
		}
		fileFormat := split[len(split)-1]
		switch fileFormat {
		case JSExtension:
			return e.MinifyJS(content)
		case JsonExtension:
			return e.MinifyJson(content)
		case HTMLExtension:
			return e.MinifyHTML(content)
		default:
			return content, nil
		}
	}
}

func newScriptLoader(script Script) ContentGenerator {
	return func(e *Engine) ([]byte, error) {
		file, compilationError := compiler.Compile([]string{script.PackagePath}, "", compiler.PrepareDefaultOptions())
		if compilationError != nil {
			return nil, compilationError
		}
		defer file.Close()
		content, readError := io.ReadAll(file)
		if readError != nil {
			return nil, readError
		}
		return e.MinifyJS(content)
	}
}

func (engine *Engine) prepareStatic() error {
	for _, script := range engine.Scripts {
		if _, ok := engine.paths[script.WebPagePath]; ok {
			return errors.New(fmt.Sprintf(PathAlreadyInUse, script.WebPagePath))
		}
		engine.paths[script.WebPagePath] = newScriptLoader(script)
	}
	for _, asset := range engine.Assets {
		if _, ok := engine.paths[asset.WebPagePath]; ok {
			return errors.New(fmt.Sprintf(PathAlreadyInUse, asset.WebPagePath))
		}
		engine.paths[asset.WebPagePath] = newAssetLoader(asset)
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

func generate(e *Engine, file *os.File, generator ContentGenerator) error {
	content, contentGenerationError := generator(e)
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
	return nil
}
func (engine *Engine) Generate(output string) error {
	preparationError := engine.prepareStatic()
	if preparationError != nil {
		return preparationError
	}
	generationErrors := make(chan error, len(engine.paths))
	for webPath, contentGenerator := range engine.paths {
		file, creationError := createFile(path.Join(output, webPath))
		if creationError != nil {
			return creationError
		}
		go func(generator ContentGenerator) {
			err := generate(engine, file, generator)
			if err != nil {
				generationErrors <- err
			} else {
				generationErrors <- nil
			}
		}(contentGenerator)
	}
	for range engine.paths {
		err := <-generationErrors
		if err != nil {
			return err
		}
	}
	return nil
}
