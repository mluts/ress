package main

type ajaxRequester interface {
	AjaxJsonRequest(method, url string) (ajaxResponse, error)
}

type ajaxResponse interface {
	GetCode() int
	GetBody() interface{}
}
