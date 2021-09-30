package engine

import (
	"github.com/shoriwe/static/pkg/engine"
	"testing"
)

func TestEngine(t *testing.T) {
	var (
		templates engine.Templates
		scripts   engine.Scripts
		assets    engine.Assets
		loadError error
	)
	templates, loadError = engine.LoadTemplates("sample-1/templates")
	if loadError != nil {
		t.Fatal(loadError)
	}
	scripts, loadError = engine.LoadScripts("sample-1/scripts", "js")
	if loadError != nil {
		t.Fatal(loadError)
	}
	assets, loadError = engine.LoadAssets("sample-1/assets", "static")
	if loadError != nil {
		t.Fatal(loadError)
	}
	myEngine := engine.NewEngine(
		templates,
		scripts,
		assets,
	)
	myEngine.HandlePath("/index.html",
		func(engine *engine.Engine) ([]byte, error) {
			scriptUrl, getError := engine.ScriptURL("hello-world")
			if getError != nil {
				return nil, getError
			}
			return engine.RenderTemplate("index.html",
				map[string]string{
					"url": scriptUrl,
				},
			)
		},
	)
	generationError := myEngine.Generate("output/sample-1")
	if generationError != nil {
		t.Fatal(generationError)
	}
}
