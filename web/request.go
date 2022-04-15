//+build tinygo wasm,js

package web

import (
	"syscall/js"
)

const (
	xhttp__constructor = "xhttp"

	XHTTP__response = "response"
	XHTTP__error    = "Error"

	function__xhttp_open = "open"
	function__xhttp_send = "send"

	PROMISE        = "promise"
	PROMISE__then  = "then"
	PROMISE__catch = "catch"
)

var (
	xhttp js.Value
)

func init() {
	xhttp = js.Global().Get(xhttp__constructor)
}

func NewXHTTPRequest(_type string, _url string, _async bool) (js.Value, error) {
	err := ValidJSValue(xhttp__constructor, xhttp)
	if err != nil {
		return js.ValueOf(nil), err
	}

	xhttp.Call(function__xhttp_open, _type, _url, _async)
	resp := xhttp.Call(function__xhttp_send)
	if err = ValidJSValue(XHTTP__response, resp); err != nil {
		return js.ValueOf(nil), err
	}

	return resp, nil
}
