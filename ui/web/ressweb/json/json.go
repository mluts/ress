package json

import (
	"github.com/gopherjs/gopherjs/js"
)

// Stringify generates JSON string using Web API
func Stringify(data interface{}) string {
	return js.Global.Get("JSON").Call("stringify", data).String()
}
