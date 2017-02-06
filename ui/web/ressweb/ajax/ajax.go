package ajax

import (
	"errors"
	"github.com/gopherjs/gopherjs/js"
)

// XHR ajax implementation
type XHR struct{}

// Response returned by XHR
type Response struct {
	Code  int
	Body  interface{}
	Error error
}

// New XHR request
func New() *XHR {
	return &XHR{}
}

// JSONRequest makes  XHR request with json responseType "json"
func (x *XHR) JSONRequest(method, url string) chan *Response {
	var (
		ch  chan *Response
		err error
		xhr *js.Object
	)
	ch = make(chan *Response, 1)

	xhr = js.Global.Get("XMLHttpRequest").New()

	xhr.Call("open", method, url)
	xhr.Set("responseType", "json")
	xhr.Set("onloadend", func() {
		go func() {
			ch <- &Response{
				xhr.Get("status").Int(),
				xhr.Get("response").Interface(),
				err,
			}
		}()
	})
	xhr.Set("onerror", func() {
		err = errors.New("An error occured during the request")
	})
	xhr.Call("send")

	return ch
}
