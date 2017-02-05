package console

import (
	"github.com/gopherjs/gopherjs/js"
)

// Log calls console.log
func Log(msg ...interface{}) {
	js.Global.Get("console").Call("log", msg...)
}

// Err calls console.error
func Err(msg ...interface{}) {
	js.Global.Get("console").Call("error", msg...)
}
