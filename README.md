# static

A simple static frontend generator library.

## Examples

- [https://shoriwe.github.io](https://shoriwe.github.io)

## Features

- Go to JavaScript (**Powered by [GopherJS](https://github.com/gopherjs/gopherjs)**)
- Go to WASM (**Powered by the standard library**)
- Templating system
- Templating scripting (**Powered by [gplasma](https://github.com/shoriwe/gplasma)**)

## Important

If you are on **Windows** use the WSL to compile your project. This happens cause **GopherJS** still does not support
this platform.

## Quicklook

For this example we will virtually store our scripts in the root of our project, in a folder called `scripts`. Something
similar with the assets and templates.

`my-project/main.go`

```go
package main

import (
	"github.com/shoriwe/static/pkg/engine"
)

func main() {
	var (
		templates engine.Templates
		scripts   engine.Scripts
		assets    engine.Assets
		loadError error
	)
	templates, loadError = engine.LoadTemplates("./templates")
	if loadError != nil {
		panic(loadError)
	}
	scripts, loadError = engine.LoadScripts(
		"./scripts", // This is the directory where the `Go` to then `JavaScript` scripts will be located
		"js",        // When the project start to be generated, the compiled JavaScript will be stored in this folder relative to the output one
	)
	if loadError != nil {
		panic(loadError)
	}
	assets, loadError = engine.LoadAssets(
		"./assets", // This is the directory where all static assets will be searched for
		"static",   // This is the directory where all assets will be stored once the project is generated.
	)
	if loadError != nil {
		panic(loadError)
	}
	myProjectEngine := engine.NewEngine(
		templates,
		scripts,
		assets,
	)
	handlingError := myProjectEngine.HandlePath("/index.html", // Define a new path to generate
		func(e *engine.Engine) ([]byte, error) {
			scriptUrl, getError := e.ScriptURL("hello-world")
			if getError != nil {
				return nil, getError
			}
			renderedTemplate, renderError := e.RenderTemplate(
				"index.html", // Search for the index.html templates in the specified templates directory
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
	if handlingError != nil { // These throws and error when the path was already in use
		panic(handlingError)
	}
	generationError := myProjectEngine.Generate("output/sample-1") // Create the path "output/sample-1" and write to it all assets, scripts and defined paths
	if generationError != nil {
		panic(generationError)
	}
}
```

Finally, to generate the page simply execute in the root of your project

```shell
go run
```

## Documentation

You can find the documentation on the wiki section.
