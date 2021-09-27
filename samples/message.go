package main

import (
	"github.com/gopherjs/gopherjs/js"
)

func Hello() {
	println("Hello World")
}

func main(){
	js.Global.Set("Hello", Hello)
}
