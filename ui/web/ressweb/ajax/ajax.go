package main

import (
	"errors"
	"github.com/gopherjs/gopherjs/js"
)

// XHR ajax implementation
type XHR struct {
	o *js.Object
}

// Response returned by XHR
type Response struct {
	Code int
	Body interface{}
}

// New XHR request
func New() *XHR {
	return &XHR{js.Global.Get("XMLHttpRequest").New()}
}

// AjaxJSONRequest makes  XHR request with json responseType "json"
func (xhr *XHR) AjaxJSONRequest(method, url string) (response *Response, err error) {
	var ch chan int

	xhr.o.Call("open", method, url)
	xhr.o.Set("responseType", "json")
	xhr.o.Set("onloadend", func() { ch <- 1 })
	xhr.o.Set("onerror", func() {
		err = errors.New("An error occured during the request")
	})
	<-ch

	return &Response{
		xhr.o.Get("status").Int(),
		xhr.o.Get("response").Interface(),
	}, err
}

// GetCode returns HTTP status code
func (r *Response) GetCode() int {
	return r.Code
}

// GetBody returns response body
func (r *Response) GetBody() interface{} {
	return r.Body
}
