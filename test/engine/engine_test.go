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
	handlingError := myEngine.HandlePath("/index.html",
		func(e *engine.Engine) ([]byte, error) {
			scriptUrl, getError := e.ScriptURL("hello-world")
			if getError != nil {
				return nil, getError
			}
			renderedTemplate, renderError := e.RenderTemplate("index.html",
				map[string]string{
					"url": scriptUrl,
				},
			)
			if renderError != nil {
				return nil, renderError
			}
			return e.MinifyHTML(renderedTemplate)
		},
	)
	if handlingError != nil {
		t.Fatal(handlingError)
	}
	generationError := myEngine.Generate("output/sample-1")
	if generationError != nil {
		t.Fatal(generationError)
	}
}
